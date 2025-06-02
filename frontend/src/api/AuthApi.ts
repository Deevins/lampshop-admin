import { axiosInstance, handleResponse, BASE_URL } from "./api";

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    token: string;
}

export const login = async (
    username: string,
    password: string
): Promise<string> => {
    const data = await handleResponse<LoginResponse>(
        axiosInstance.post(`${BASE_URL}/login`, { username, password })
    );
    return data.token;
};

export const logout = () => {
    localStorage.removeItem("jwtToken");

};
