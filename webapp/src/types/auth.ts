export interface User {
    id: number;
    username: string;
    role: "user" | "admin";
    created_at: string;
    updated_at: string;
}