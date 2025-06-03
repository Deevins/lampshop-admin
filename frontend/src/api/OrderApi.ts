import {axiosInstance, BASE_URL, handleResponse} from "./api";
import type {AttributeOption, Category} from "../types/Product.ts";
import type {Order, OrderStatus} from "../types/Order.ts";

export const getCategories = async (): Promise<Category[]> => {
    return handleResponse<Category[]>(
        axiosInstance.get(`${BASE_URL}/categories`)
    );
};

export const getAttributeOptions = async (
    category_id: string
): Promise<AttributeOption[]> => {
    console.log(category_id)
    return handleResponse<AttributeOption[]>(
        axiosInstance.get(
            `${BASE_URL}/categories/${category_id}/attributes`
        )
    );
};

export const getOrders = async (): Promise<Order[]> => {
    return handleResponse<Order[]>(
        axiosInstance.get(`${BASE_URL}/orders`)
    );
};

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