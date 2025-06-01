import axios, {type AxiosResponse, type InternalAxiosRequestConfig} from "axios";


export const BASE_URL = "http://localhost:8080";


// Создаём axios-инстанс, чтобы можно было глобально настроить interceptor
export const axiosInstance = axios.create();

// Интерцептор: если в localStorage есть jwtToken, кладём его в заголовок Authorization
axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const token = localStorage.getItem("jwtToken");
        if (token) {
            // Гарантируем, что headers не undefined
            config.headers = config.headers ?? {};
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Обёртка для обработки ответа и ошибок
export async function handleResponse<T>(
    promise: Promise<AxiosResponse<T>>
): Promise<T> {
    try {
        const response = await promise;
        return response.data;
    } catch (err: any) {
        // Если получен 401 — очищаем токен; ProtectedRoute редиректит на /login
        if (err.response?.status === 401) {
            localStorage.removeItem("jwtToken");
        }
        const msg =
            err.response?.data?.error || `HTTP error! status: ${err.response?.status}`;
        throw new Error(msg);
    }
}
