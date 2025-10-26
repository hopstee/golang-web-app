import { apiFetch } from "@/lib/api";
import type { Page, PageData, PageSchema } from "@/types/pages";

export async function fetchAllSchemas(): Promise<PageSchema[]> {
    return await apiFetch("/api/v1/admin/pages");
}

export async function fetchPageSchema(slug: string): Promise<PageSchema> {
    return await apiFetch(`/api/v1/admin/pages/${slug}/schema`);
}

export async function fetchPageData(slug: string): Promise<PageData> {
    return await apiFetch(`/api/v1/admin/pages/${slug}/data`);
}

export async function updatePageData(slug: string, data: Page): Promise<void> {
    return apiFetch(`/api/v1/admin/pages/${slug}/data`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });
}

export async function deletePageData(slug: string): Promise<void> {
    return apiFetch(`/api/v1/admin/pages/${slug}`, { method: "DELETE", credentials: "include" });
}
