import {axiosInstance, BASE_URL, handleResponse} from "./api";
import type {AttributeOption, Category} from "../types/Product.ts";
import type {Order, OrderStatus} from "../types/Order.ts";

// Получить категории
export const getCategories = async (): Promise<Category[]> => {
    return handleResponse<Category[]>(
        axiosInstance.get(`${BASE_URL}/categories`)
    );
};

// Получить опции атрибутов для категории
export const getAttributeOptions = async (
    categoryId: string
): Promise<AttributeOption[]> => {
    return handleResponse<AttributeOption[]>(
        axiosInstance.get(
            `${BASE_URL}/categories/${encodeURIComponent(categoryId)}/attributes`
        )
    );
};

// Получить все заказы
export const getOrders = async (): Promise<Order[]> => {
    return handleResponse<Order[]>(
        axiosInstance.get(`${BASE_URL}/orders`)
    );
};

// Обновить статус заказа
export const updateOrderStatus = async (
    id: number,
    status: OrderStatus
): Promise<Order> => {
    return handleResponse<Order>(
        axiosInstance.put(
            `${BASE_URL}/orders/${id}/status`,
            {status},
            {headers: {"Content-Type": "application/json"}}
        )
    );
};