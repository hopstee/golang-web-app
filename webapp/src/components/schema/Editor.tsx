import { useEntityStore } from "@/store/EntitiesStore";
import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { Spinner } from "@/components/ui/spinner";
import { centerContent } from "@/lib/render";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import FieldsScaffold from "./FieldsScaffold";
import { EntityDataTypes, type EntityData } from "@/types/entities";
import { Button } from "../ui/button";
import { Save } from "lucide-react";
import { toast } from "sonner";
import { deleteFile, uploadFile } from "@/api/files";
import { getFileType } from "@/lib/utils";

interface EditorProps {
    type: string;
}

export default function Editor(props: EditorProps) {
    const { slug } = useParams<{ slug: string }>();
    const { entityData, entitySchema, fetchEntity, saveEntityData, loading, updating } = useEntityStore();
    const [searchParams, setSearchParams] = useSearchParams();

    const currentTabFromUrl = searchParams.get("tab") || "content";
    const [selectedTab, setSelectedTab] = useState<string>(currentTabFromUrl);

    const [pageValues, setPageValues] = useState<EntityData>({ content: {} });
    const [changed, setChanged] = useState(false);

    useEffect(() => {
        if (slug) fetchEntity(props.type, slug);
    }, [fetchEntity, slug]);

    const page = slug ? entityData[slug] : null;
    const schema = slug ? entitySchema[slug] : null;

    const pageHeader = props.type === "page"
        ? "Редактор страницы"
        : "Редактор компонента"

    useEffect(() => {
        if (page && slug) {
            setPageValues(page);
            setChanged(false);
        }
    }, [page, slug]);

    useEffect(() => {
        if (schema && slug) {
            const availableTabs = ["content", ...(schema.children?.map((c) => c.id) || [])];
            console.log(availableTabs)
            console.log(selectedTab)
            console.log(availableTabs.includes(selectedTab))
            if (!availableTabs.includes(selectedTab)) {
                const firstTab = availableTabs[0] || "content";
                console.log(firstTab)
                handleTabChange(firstTab);
            }
        }
    }, [schema, slug]);

    const handleTabChange = (tab: string) => {
        setSelectedTab(tab);
        setSearchParams({ tab })
    }

    if (loading) return centerContent(<Spinner />);
    if (!slug || !page || !schema) return centerContent("Нет данных");

    const handlePartChange = (key: string, value: Record<string, unknown>) => {
        setPageValues((prev) => ({
            ...prev,
            [key]: value,
        }));
        setChanged(true);
    }

    const handleSave = async () => {
        if (!slug || !changed) return;

        const uploadImg = async (blobURL: string): Promise<string> => {
            let res: string;
            try {
                const response = await fetch(blobURL);
                const blob = await response.blob();
                const formData = new FormData();
                const fileType = getFileType(blob.type);
                formData.append("file", blob, `uploaded_file.${fileType}`);

                res = await uploadFile(formData);
            } catch (error) {
                throw new Error("Ошибка при получении данных изображения");
            }
            return res;
        }

        const prepareFiles = async (data: unknown): Promise<unknown> => {
            if (typeof data === "string") {
                if ((data as string).startsWith("blob:")) {
                    return await uploadImg(data);
                }
                if ((data as string).startsWith("deleted:")) {
                    if ((data as string).indexOf("blob:") === -1) {
                        const withoutDeletedPrefix = (data as string).replace("deleted:", "");
                        try {
                            await deleteFile(withoutDeletedPrefix);
                        } catch (error) {
                            console.error("Ошибка при удалении файла:", error);
                        }
                    }
                    return "";
                }
                return data;
            }

            if (Array.isArray(data)) {
                return await Promise.all(data.map(prepareFiles));
            }


            if (typeof data === "object" && data !== null) {
                const newData: Record<string, unknown> = {};
                for (const [key, value] of Object.entries(data)) {
                    newData[key] = await prepareFiles(value);
                }
                return newData;
            }

            return data;
        }

        toast.promise(
            async () => {
                const updatedValues = await prepareFiles(pageValues);
                await saveEntityData(slug, (updatedValues as Record<string, unknown>));
                setChanged(false);
            },
            {
                loading: "Обновляем данные...",
                success: "Данные страницы обновленны",
                error: "Не удалось обновить данные страницы",
            }
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex justify-between">
                <h1 className="text-muted-foreground">{pageHeader} <span className="font-bold text-foreground">{schema.title}</span></h1>
                <Button
                    onClick={handleSave}
                    disabled={!changed || updating}
                    size="sm"
                    className="disabled:cursor-not-allowed"
                >
                    {updating ? <Spinner /> : <Save />}
                    Сохранить
                </Button>
            </div>

            <Tabs value={selectedTab} onValueChange={handleTabChange} className="space-y-4 max-w-3xl w-full mx-auto">
                <TabsList>
                    <TabsTrigger value="content">Содержимое страницы</TabsTrigger>

                    {schema.children?.length > 0 && schema.children.map((child) => (
                        <TabsTrigger key={child.id} value={child.id}>{child.title}</TabsTrigger>
                    ))}
                </TabsList>

                <TabsContent value="content">
                    <FieldsScaffold
                        fields={schema.content || []}
                        data={pageValues.content || {}}
                        onChange={(data) => handlePartChange(EntityDataTypes.CONTENT, data)}
                    />
                </TabsContent>

                {schema.children?.length > 0 && schema.children.map((child) => (
                    <TabsContent key={child.id} value={child.id}>
                        <FieldsScaffold
                            fields={child.content || []}
                            data={pageValues[child.id] as Record<string, unknown> || {}}
                            onChange={(data) => handlePartChange(child.id, data)}
                        />
                    </TabsContent>
                ))}
            </Tabs>
        </div>
    )
}
