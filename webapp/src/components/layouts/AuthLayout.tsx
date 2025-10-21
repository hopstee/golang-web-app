import { useAuth } from "@/hooks/useAuth";
import { useEffect, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";

interface AuthLayoutProps {
    children: ReactNode
}

export default function AuthLayout({ children }: AuthLayoutProps) {
    const navigate = useNavigate();
    const { user, isLoading } = useAuth();

    useEffect(() => {
        if (!isLoading && user) {
            navigate("/admin/dashboard");
        }
    }, [isLoading, user, navigate]);

    return (
        <div className="w-screen h-screen flex items-center justify-center">
            {children}
        </div>
    )
}