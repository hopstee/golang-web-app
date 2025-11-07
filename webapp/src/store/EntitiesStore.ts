import { fetchEntitiesNames, fetchEntityData, fetchEntitySchema, updateEntityData } from "@/api/entities"
import type { EntityData, EntityDataData, EntitySchema } from "@/types/pages"
import { create } from "zustand"
import { persist } from "zustand/middleware"

interface EntitiesState {
    pages: string[]
    shared: string[]
    pageData: Record<string, EntityData>
    sharedData: Record<string, EntityData>
    loading: boolean
    updating: boolean
    error: boolean
    errorMessage: string
    fetchAll: (force?: boolean) => Promise<void>
    fetchEntity: (type: string, slug: string) => Promise<void>
    saveEntityData: (slug: string, updated: EntityDataData) => Promise<void>
}

export const useEntityStore = create<EntitiesState>()(
    persist(
        (set, get) => ({
            pages: [],
            shared: [],
            pageData: {},
            sharedData: {},
            loading: false,
            updating: false,
            error: false,
            errorMessage: "",

            async fetchAll(force = false) {
                if (!force && get().pages.length > 0) return;

                let pages: string[] = [];
                let shared: string[] = [];
                set({ loading: true });
                try {
                    pages = await fetchEntitiesNames("page");
                    shared = await fetchEntitiesNames("shared");
                } catch (err) {
                    set({ error: true, errorMessage: (err as Error).message });
                } finally {
                    set({ loading: false, pages, shared });
                }
            },

            async fetchEntity(type: string, slug: string) {
                if (get().pageData[slug]) return;

                let data: EntityData | undefined;
                let schema: EntitySchema;
                set({ loading: true });
                try {
                    data = await fetchEntityData(slug);
                    schema = await fetchEntitySchema(type, slug);
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

            async saveEntityData(slug, data) {
                set({ updating: true });
                try {
                    await updateEntityData(slug, data);
                    set(state => ({
                        pageData: {
                            ...state.pageData,
                            [slug]: { ...state.pageData[slug], data }
                        },
                    }));
                } catch (err) {
                    set({ error: true, errorMessage: (err as Error).message });
                } finally {
                    set({ updating: false });
                }
            },
        }),
        { name: "entities-store" }
    )
)