import { Link } from "react-router-dom";

export default function NotFound() {
    return (
        <div className="p-6 text-center">
            <h1 className="text-4xl font-bold">404</h1>
            <p className="mt-4">Страница не найдена</p>
            <Link to="/admin/dashboard" className="mt-2 inline-block text-blue-600 hover:underline">
                В админку
            </Link>
        </div>
    );
}
