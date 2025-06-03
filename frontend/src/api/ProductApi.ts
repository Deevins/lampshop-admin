import {axiosInstance, BASE_URL, handleResponse} from "./api";
import type {AttributeOption, Product} from "../types/Product.ts";

export const getProducts = async (): Promise<Product[]> => {
    return handleResponse<Product[]>(
        axiosInstance.get(`${BASE_URL}/products`)
    );
};

export const getProductById = async (id: string): Promise<Product> => {
    return handleResponse<Product>(
        axiosInstance.get(`${BASE_URL}/products/${id}`)
    );
};

type UpsertProductDTO = Omit<Product, "id" | "createdAt" | "updatedAt">;

export const addProduct = async (
    data: UpsertProductDTO
): Promise<Product> => {
    return handleResponse<Product>(
        axiosInstance.post(`${BASE_URL}/products`, data, {
            headers: {"Content-Type": "application/json"},
        })
    );
};

export const updateProduct = async (
    id: string,
    data: UpsertProductDTO
): Promise<Product> => {
    return handleResponse<Product>(
        axiosInstance.put(`${BASE_URL}/products/${id}`, data, {
            headers: {"Content-Type": "application/json"},
        })
    );
};

export const deleteProduct = async (id: string): Promise<void> => {
    await handleResponse<void>(
        axiosInstance.delete(`${BASE_URL}/products/${id}`)
    );
};

export const getAttributeOptions = async (
    category_id: string
): Promise<AttributeOption[]> => {
    return handleResponse<AttributeOption[]>(
        axiosInstance.get(
            `${BASE_URL}/categories/${category_id}/attributes`
        )
    );
};