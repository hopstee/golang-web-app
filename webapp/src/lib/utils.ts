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
