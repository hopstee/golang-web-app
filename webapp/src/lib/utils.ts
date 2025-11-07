import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function normalizeErrors(errors?: unknown): { message: string }[] | undefined {
    if (!errors) return undefined;
    if (Array.isArray(errors)) {
        return errors.map((e) =>
            typeof e === "string"
                ? { message: e }
                : "message" in (e as Error)
                    ? { message: (e as Error).message }
                    : { message: String(e) }
        );
    }
    if (typeof errors === "string") return [{ message: errors }];
    if (typeof errors === "object" && "message" in (errors as Error))
        return [{ message: (errors as Error).message }];
    return [{ message: String(errors) }];
};

export function getDefaultValue(type: string): unknown {
    if (type.startsWith("list[")) return [];
    if (type === "bool" || type === "boolean") return false;
    return "";
}

export function getFileType(type: string): string {
    const mimeToExt: Record<string, string> = {
        "image/png": "png",
        "image/jpeg": "jpg",
        "image/webp": "webp",
        "image/svg+xml": "svg",
        "image/gif": "gif",
        "application/pdf": "pdf",
        "text/plain": "txt",
        "application/json": "json",
        "application/zip": "zip",
        "audio/mpeg": "mp3",
        "audio/wav": "wav",
        "video/mp4": "mp4",
    };

    return mimeToExt[type] ?? "bin";
}
