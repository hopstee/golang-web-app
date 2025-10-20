import { AuthContext, type AuthContextType } from "@/context/AuthContext";
import { useContext } from "react";

export function useAuth(): AuthContextType {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error("useAuth must be used within AuthProvider");
    return ctx;
}