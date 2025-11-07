import { apiFetch } from "@/lib/api";
import type { EntityData, EntityDataData, EntitySchema } from "@/types/pages";

export async function fetchEntitiesNames(type: string): Promise<string[]> {
    return await apiFetch(`/api/v1/admin/entity/${type}/names`);
}

export async function fetchEntitySchema(type: string, slug: string): Promise<EntitySchema> {
    return await apiFetch(`/api/v1/admin/entity/${type}/${slug}/schema`);
}

export async function fetchEntityData(slug: string): Promise<EntityData> {
    return await apiFetch(`/api/v1/admin/entity/${slug}/data`);
}

export async function updateEntityData(slug: string, data: EntityDataData): Promise<void> {
    return apiFetch(`/api/v1/admin/entity/${slug}/data`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });
}

export async function deleteEntityData(slug: string): Promise<void> {
    return apiFetch(`/api/v1/admin/entity/${slug}`, { method: "DELETE", credentials: "include" });
}
