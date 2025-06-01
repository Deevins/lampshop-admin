import React from "react";
import { Navigate } from "react-router-dom";

// Простая проверка: если в localStorage есть jwtToken, считаем, что авторизация пройдена.
// (В реальном приложении можно распарсить токен и проверить срок годности.)
function isAuthenticated(): boolean {
    return Boolean(localStorage.getItem("jwtToken"));
}

interface ProtectedRouteProps {
    children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    if (!isAuthenticated()) {
        return <Navigate to="/login" replace />;
    }
    return <>{children}</>;
};

export default ProtectedRoute;
