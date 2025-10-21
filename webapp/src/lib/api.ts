const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "";

async function apiFetch<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const res = await fetch(url, {
        ...options,
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            ...(options.headers || {}),
        },
    });

    if (!res.ok) {
        const text = await res.text();
        throw new Error(`API error (${res.status}): ${text}`);
    }

    const text = await res.text();

    try {
        return JSON.parse(text) as T;
    } catch {
        return text as unknown as T;
    }
}

export { apiFetch };
