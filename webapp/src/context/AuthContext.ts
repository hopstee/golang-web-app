import { createContext } from "react";
import type { User } from "../types/auth";

export interface AuthContextType {
    user: User | undefined;
    isLoading: boolean;
    isError: boolean;
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | null>(null);
