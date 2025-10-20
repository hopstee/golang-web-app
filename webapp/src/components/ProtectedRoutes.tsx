import { useAuth } from "@/hooks/useAuth";
import { Navigate, Outlet } from "react-router-dom";

export const ProtectedRoute: React.FC = () => {
    const { user, isLoading, isError } = useAuth();

    if (isLoading) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="h-10 w-10 animate-spin rounded-full border-4 border-gray-300 border-t-transparent" />
            </div>
        );
    }

    if (isError || !user) {
        return <Navigate to="/admin/login" replace />;
    }

    return <Outlet />;
};
