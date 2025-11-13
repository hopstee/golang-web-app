import { fetchEntitiesNames, fetchEntityData, fetchEntitySchema, updateEntityData } from "@/api/entities"
import type { EntityData, EntitySchema, ShortEntityData } from "@/types/entities"
import { create } from "zustand"
// import { persist } from "zustand/middleware"

interface EntitiesState {
    pages: ShortEntityData[]
    shared: ShortEntityData[]
    entityData: Record<string, EntityData>
    entitySchema: Record<string, EntitySchema>
    loading: boolean
    updating: boolean
    error: boolean
    errorMessage: string
    fetchAll: (force?: boolean) => Promise<void>
    fetchEntity: (type: string, slug: string) => Promise<void>
    saveEntityData: (slug: string, updated: Record<string, unknown>) => Promise<void>
}

export const useEntityStore = create<EntitiesState>()(
    // persist(
        (set, get) => ({
            pages: [],
            shared: [],
            entityData: {},
            entitySchema: {},
            loading: false,
            updating: false,
            error: false,
            errorMessage: "",

            async fetchAll(force = false) {
                if (!force && get().pages.length > 0 && get().shared.length > 0) return;

                let pages: ShortEntityData[] = [];
                let shared: ShortEntityData[] = [];
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
                if (get().entityData[slug] && get().entitySchema[slug]) return;

                let data: EntityData | undefined;
                let schema: EntitySchema | undefined;
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
                    console.log(data)
                    set(state => ({ entityData: { ...state.entityData, [slug]: { ...data } } }));
                }
                if (schema) {
                    set(state => ({ entitySchema: { ...state.entitySchema, [slug]: { ...schema } } }));
                }

                set({ loading: false });
            },

            async saveEntityData(slug, data) {
                set({ updating: true });
                try {
                    await updateEntityData(slug, data);

                    set(state => ({
                        entityData: {
                            ...state.entityData,
                            [slug]: { ...state.entityData[slug], ...data }
                        },
                    }));
                } catch (err) {
                    set({ error: true, errorMessage: (err as Error).message });
                } finally {
                    set({ updating: false });
                }
            },
        }),
    //     { name: "entities-store" }
    // )
)