import { apiFetch } from "@/lib/api";
import type { User } from "@/types/auth";

export async function fetchMe(): Promise<User> {
    return await apiFetch("/api/v1/admin/auth/me", { credentials: "include" });
}

export async function loginRequest(username: string, password: string): Promise<User> {
    return apiFetch("/api/v1/admin/auth/login", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    });
}

export async function logoutRequest(): Promise<void> {
    return apiFetch("/api/v1/admin/auth/logout", { method: "POST", credentials: "include" });
}
