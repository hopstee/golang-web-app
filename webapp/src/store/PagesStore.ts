import { fetchAllSchemas, fetchPageData, updatePageData } from "@/api/pages"
import type { PageData, PageDataData, PageSchema } from "@/types/pages"
import { create } from "zustand"
import { persist } from "zustand/middleware"

interface PagesState {
    pages: PageSchema[]
    pageData: Record<string, PageData>
    loading: boolean
    error: boolean
    errorMessage: string
    fetchAll: (force?: boolean) => Promise<void>
    fetchPage: (slug: string) => Promise<void>
    savePageData: (slug: string, updated: PageDataData) => Promise<void>
}

export const usePageStore = create<PagesState>()(
    persist(
        (set, get) => ({
            pages: [],
            pageData: {},
            loading: false,
            error: false,
            errorMessage: "",

            async fetchAll(force = false) {
                if (!force && get().pages.length > 0) return;

                let schemas: PageSchema[] = [];
                set({ loading: true });
                try {
                    schemas = await fetchAllSchemas();
                } catch (err) {
                    set({ error: true, errorMessage: (err as Error).message });
                } finally {
                    set({ loading: false, pages: schemas });
                }
            },

            async fetchPage(slug: string) {
                if (get().pageData[slug]) return;

                let data: PageData | undefined;
                set({ loading: true });
                try {
                    data = await fetchPageData(slug);
                } catch (error) {
                    set({ error: true, errorMessage: (error as Error).message });
                } finally {
                    set({ loading: false });
                }

                if (data) {
                    set(state => ({
                        pageData: { ...state.pageData, [slug]: { ...data } },
                        loading: false,
                    }));
                }
            },

            async savePageData(slug, data) {
                set({ loading: true });
                try {
                    await updatePageData(slug, data);
                    set(state => ({
                        pageData: {
                            ...state.pageData,
                            [slug]: { ...state.pageData[slug], data }
                        },
                    }));
                } catch (err) {
                    set({ error: true, errorMessage: (err as Error).message });
                } finally {
                    set({ loading: false });
                }
            },
        }),
        { name: "pages-store" }
    )
)