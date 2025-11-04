import { apiFetch } from "@/lib/api";

export async function uploadFile(formData: FormData): Promise<string> {
    return await apiFetch("/api/v1/admin/files", {
        method: "POST",
        credentials: "include",
        body: formData,
    });
}

export async function deleteFile(path: string): Promise<void> {
    return await apiFetch(`/api/v1/admin/files?path=${encodeURIComponent(path)}`, {
        method: "DELETE",
        credentials: "include",
    });
}