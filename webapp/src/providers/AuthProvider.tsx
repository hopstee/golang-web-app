import { fetchMe, loginRequest, logoutRequest } from "@/api/auth";
import { AuthContext } from "@/context/AuthContext";
import type { User } from "@/types/auth";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";

interface AuthProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
    const queryClient = useQueryClient();

    const {
        data: user,
        isLoading,
        isFetching,
        isError,
    } = useQuery<User, Error>({
        queryKey: ["me"],
        queryFn: fetchMe,
        retry: false,
        refetchOnWindowFocus: false,
    });

    const loginMutation = useMutation<User, Error, { email: string; password: string }>({
        mutationFn: async ({ email, password }) => {
            await loginRequest(email, password);
            const me = await fetchMe();
            queryClient.setQueryData(["me"], me);
            return me;
        },
    });

    const logoutMutation = useMutation<void, Error>({
        mutationFn: async () => {
            await logoutRequest();
            queryClient.removeQueries({ queryKey: ["me"] });
        },
    });

    const login = async (email: string, password: string): Promise<void> => {
        await loginMutation.mutateAsync({ email, password });
    };

    const logout = async (): Promise<void> => {
        await logoutMutation.mutateAsync();
    };

    const isAuthLoading =
        isLoading || isFetching || loginMutation.isPending || logoutMutation.isPending;

    return (
        <AuthContext.Provider value={{ user, isLoading: isAuthLoading, isError, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};