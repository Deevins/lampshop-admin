import type {Category} from "../types/Product.ts";
import type {AxiosResponse} from "axios";

export const categories: Category[] = [
    {id: "bulb", name: "Лампочки"},
    {id: "cable", name: "Кабели"},
    {id: "equipment", name: "Оборудование"}
];

export interface AttributeOption {
    key: string;
    label: string;
    type: "text" | "number";
}

export const attributeOptions: Record<string, AttributeOption[]> = {
    bulb: [
        {key: "power", label: "Мощность (Вт)", type: "number"},
        {key: "color", label: "Цвет", type: "text"},
        {key: "temperature", label: "Температура (K)", type: "number"},
        {key: "socketType", label: "Тип цоколя", type: "text"}
    ],
    cable: [
        {key: "length", label: "Длина (м)", type: "number"},
        {key: "material", label: "Материал", type: "text"},
        {key: "color", label: "Цвет", type: "text"}
    ],
    equipment: [
        {key: "manufacturer", label: "Производитель", type: "text"},
        {key: "model", label: "Модель", type: "text"},
        {key: "warranty", label: "Гарантия (мес.)", type: "number"}
    ]
};


export const BASE_URL = "http://localhost:8080";

// Вспомогательная функция для обработки ответа
// Если нужно, можно расширить логику обработки ошибок
export async function handleResponse<T>(promise: Promise<AxiosResponse<T>>): Promise<T> {
    try {
        const response = await promise;
        return response.data;
    } catch (err: unknown) {
        // @ts-ignore
        throw new Error(err.response?.data?.error || `HTTP error! status: ${err.response?.status}`);
    }
}