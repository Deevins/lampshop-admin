// src/services/authApi.ts
import { axiosInstance, handleResponse, BASE_URL } from "./api";

// Тело запроса для /login
export interface LoginRequest {
    username: string;
    password: string;
}

// Тело ответа от /login
export interface LoginResponse {
    token: string;
}

// Выполняем POST /login
export const login = async (
    username: string,
    password: string
): Promise<string> => {
    const data = await handleResponse<LoginResponse>(
        axiosInstance.post(`${BASE_URL}/login`, { username, password })
    );
    return data.token;
};

// Логаут: просто очищаем локально хранимый JWT
export const logout = () => {
    localStorage.removeItem("jwtToken");
};
