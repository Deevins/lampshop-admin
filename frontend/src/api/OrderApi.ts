import type {Category} from "../types/Product.ts";
import {type AttributeOption, BASE_URL, handleResponse} from "./api.ts";
import axios from "axios";
import type {Order, OrderStatus} from "../types/Order.ts";

export const getCategories = async (): Promise<Category[]> => {
    return handleResponse<Category[]>(
        axios.get(`${BASE_URL}/categories`)
    );
};

export const getAttributeOptions = async (
    categoryId: string
): Promise<AttributeOption[]> => {
    return handleResponse<AttributeOption[]>(
        axios.get(`${BASE_URL}/categories/${encodeURIComponent(categoryId)}/attributes`)
    );
};

export const getOrders = async (): Promise<Order[]> => {
    return handleResponse<Order[]>(
        axios.get(`${BASE_URL}/orders`)
    );
};

export const updateOrderStatus = async (
    id: number,
    status: OrderStatus
): Promise<Order> => {
    return handleResponse<Order>(
        axios.put(
            `${BASE_URL}/orders/${id}/status`,
            {status},
            {
                headers: {
                    "Content-Type": "application/json",
                },
            }
        )
    );
};