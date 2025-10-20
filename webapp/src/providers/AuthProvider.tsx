import { fetchMe, loginRequest, logoutRequest } from "@/api/auth";
import { AuthContext } from "@/context/AuthContext";
import type { User } from "@/types/auth";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";

interface AuthProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
    const queryClient = useQueryClient();

    const { data: user, isLoading, isError } = useQuery<User, Error>({
        queryKey: ["me"],
        queryFn: fetchMe,
        retry: false,
    });

    const login = async (email: string, password: string) => {
        const user = await loginRequest(email, password);
        queryClient.setQueryData(["me"], user);
    };

    const logout = async () => {
        await logoutRequest();
        queryClient.removeQueries({ queryKey: ["me"] });
    };

    return (
        <AuthContext.Provider value= {{ user, isLoading, isError, login, logout }}>
            { children }
        </AuthContext.Provider>
    );
};