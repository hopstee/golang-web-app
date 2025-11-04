import { usePageStore } from "@/store/PagesStore";
import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { Spinner } from "@/components/ui/spinner";
import { centerContent } from "@/lib/render";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import FieldsScaffold from "./FieldsScaffold";
import { PageDataTypes, type PageDataData } from "@/types/pages";
import { Button } from "../ui/button";
import { Save } from "lucide-react";
import { toast } from "sonner";
import { deleteFile, uploadFile } from "@/api/files";
import { getFileType } from "@/lib/utils";

export default function Editor() {
    const { slug } = useParams<{ slug: string }>();
    const { pageData, fetchPage, savePageData, loading, updating } = usePageStore();
    const [searchParams, setSearchParams] = useSearchParams();

    const currentTabFromUrl = searchParams.get("tab") || "content";
    const [selectedTab, setSelectedTab] = useState<string>(currentTabFromUrl);

    const [pageValues, setPageValues] = useState<PageDataData>({ layout_fields: {}, content: {} });
    const [changed, setChanged] = useState(false);

    useEffect(() => {
        if (slug) fetchPage(slug);
    }, [fetchPage, slug]);

    const page = slug ? pageData[slug] : null;

    useEffect(() => {
        if (page?.data && slug) {
            setPageValues(page.data);
            setChanged(false);
        }
    }, [page?.data, slug]);

    const handleTabChange = (tab: string) => {
        setSelectedTab(tab);
        setSearchParams({ tab })
    }

    if (!slug || !page) return centerContent("Нет данных");
    if (loading) return centerContent(<Spinner />);

    const schema = page.schema || {};

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
                await savePageData(slug, (updatedValues as PageDataData));
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
                <h1 className="text-muted-foreground">Редактор страницы <span className="font-bold text-foreground">{page.schema.title}</span></h1>
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

            <Tabs defaultValue={selectedTab} onValueChange={handleTabChange} className="space-y-4 max-w-3xl w-full mx-auto">
                <TabsList>
                    <TabsTrigger value="content">Содержимое страницы</TabsTrigger>
                    <TabsTrigger value="layout">Макет</TabsTrigger>
                </TabsList>

                <TabsContent value="content">
                    <FieldsScaffold
                        fields={schema.content || []}
                        data={pageValues.content || {}}
                        onChange={(data) => handlePartChange(PageDataTypes.CONTENT, data)}
                    />
                </TabsContent>

                <TabsContent value="layout">
                    <FieldsScaffold
                        fields={schema.layout_fields || []}
                        data={pageValues.layout_fields || {}}
                        onChange={(data) => handlePartChange(PageDataTypes.LAYOUT_FIELDS, data)}
                    />
                </TabsContent>
            </Tabs>
        </div>
    )
}
