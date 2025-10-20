import { useAuth } from "@/hooks/useAuth";
import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";

interface AuthLayoutProps {
    children: ReactNode
}

export default function AuthLayout({ children }: AuthLayoutProps) {
    const navigate = useNavigate();
    const { user } = useAuth();

    if (user) {
        navigate("/admin/dashboard");
    }

    return (
        <div className="w-screen h-screen flex items-center justify-center">
            { children }
        </div>
    )
}