import { usePageStore } from "@/store/PagesStore"
import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { Spinner } from "@/components/ui/spinner"
import { centerContent } from "@/lib/render"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import FieldsScaffold from "./FieldsScaffold"
import { PageDataTypes, type PageDataData } from "@/types/pages"
import { Button } from "../ui/button"
import { Save } from "lucide-react"
import { toast } from "sonner"

export default function Editor() {
    const { slug } = useParams<{ slug: string }>()
    const { pageData, fetchPage, savePageData, loading, updating } = usePageStore()

    const [pageValues, setPageValues] = useState<PageDataData>({ layout_fields: {}, content: {} })
    const [changed, setChanged] = useState(false)

    useEffect(() => {
        if (slug) fetchPage(slug)
    }, [fetchPage, slug])

    const page = slug ? pageData[slug] : null

    useEffect(() => {
        if (page?.data && slug) {
            setPageValues(page.data)
            setChanged(false)
        }
    }, [page?.data, slug])

    if (!slug || !page) return centerContent("Нет данных")
    if (loading) return centerContent(<Spinner />)

    const schema = page.schema || {}

    const handlePartChange = (key: string, value: Record<string, unknown>) => {
        setPageValues((prev) => ({
            ...prev,
            [key]: value,
        }))
        setChanged(true)
    }

    const handleSave = async () => {
        if (!slug || !changed) return

        toast.promise(
            async () => {
                await savePageData(slug, pageValues)
                setChanged(false)
            },
            {
                loading: "Обновляем данные...",
                success: "Данные страницы обновленны",
                error: "Не удалось обновить данные страницы",
            }
        )
    }

    return (
        <div className="space-y-4">
            <div className="flex justify-between">
                <h1>Редактор страницы <span className="font-bold">{page.schema.title}</span></h1>
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

            <Tabs defaultValue="layout">
                <TabsList>
                    <TabsTrigger value="layout">Макет</TabsTrigger>
                    <TabsTrigger value="content">Содержимое страницы</TabsTrigger>
                </TabsList>

                <TabsContent value="layout">
                    <FieldsScaffold
                        fields={schema.layout_fields || []}
                        data={pageValues.layout_fields || {}}
                        onChange={(data) => handlePartChange(PageDataTypes.LAYOUT_FIELDS, data)}
                    />
                </TabsContent>

                <TabsContent value="content">
                    <FieldsScaffold
                        fields={schema.content || []}
                        data={pageValues.content || {}}
                        onChange={(data) => handlePartChange(PageDataTypes.CONTENT, data)}
                    />
                </TabsContent>
            </Tabs>
        </div>
    )
}
