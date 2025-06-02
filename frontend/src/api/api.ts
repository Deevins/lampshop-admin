import axios, {type AxiosResponse, type InternalAxiosRequestConfig} from "axios";


export const BASE_URL = "http://localhost:8083";


// Создаём axios-инстанс, чтобы можно было глобально настроить interceptor
export const axiosInstance = axios.create();

axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const token = localStorage.getItem("jwtToken");
        if (token) {
            config.headers = config.headers ?? {};
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

export async function handleResponse<T>(
    promise: Promise<AxiosResponse<T>>
): Promise<T> {
    try {
        const response = await promise;
        return response.data;
    } catch (err: any) {
        if (err.response?.status === 401) {
            localStorage.removeItem("jwtToken");
        }
        const msg =
            err.response?.data?.error || `HTTP error! status: ${err.response?.status}`;
        throw new Error(msg);
    }
}
