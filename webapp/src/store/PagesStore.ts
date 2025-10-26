import { fetchAllSchemas, fetchPageData } from "@/api/pages"
import type { PageData, PageSchema } from "@/types/pages"
import { create } from "zustand"
import { persist } from "zustand/middleware"

interface PagesState {
    pages: PageSchema[]
    pageData: Record<string, PageData>
    loading: boolean
    fetchAll: (force?: boolean) => Promise<void>
    fetchPage: (slug: string) => Promise<void>
}

export const usePageStore = create<PagesState>()(
    persist(
        (set, get) => ({
            pages: [],
            pageData: {},
            loading: false,

            async fetchAll(force = false) {
                if (!force && get().pages.length > 0) return;

                set({ loading: true });
                const schemas = await fetchAllSchemas();
                set({ loading: false, pages: schemas });
            },

            async fetchPage(slug: string) {
                if (get().pageData[slug]) return;

                set({ loading: true });
                const data: PageData = await fetchPageData(slug);
                set(state => ({
                    pageData: { ...state.pageData, [slug]: data },
                    loading: false,
                }));
            }
        }),
        { name: "pages-store" }
    )
)