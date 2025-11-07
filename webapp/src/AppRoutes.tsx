import { Routes, Route } from "react-router-dom";
import Login from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import { ProtectedRoute } from "@/components/ProtectedRoutes";
import NotFound from "@/pages/NotFound";
import EditorPage from "./pages/Editor";

const AppRoutes: React.FC = () => (
    <Routes>
        <Route path="/admin/login" element={<Login />} />

        <Route path="/admin" element={<ProtectedRoute />}>
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="page/schemas/:slug" element={<EditorPage type="page" />} />
            <Route path="shared/schemas/:slug" element={<EditorPage type="shared" />} />
        </Route>

        <Route path="*" element={<NotFound />} />
    </Routes>
);

export default AppRoutes;
