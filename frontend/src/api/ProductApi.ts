import type {Product} from "../types/Product.ts";
import {BASE_URL, handleResponse} from "./api.ts";
import axios from "axios";

export const getProducts = async (): Promise<Product[]> => {
    return handleResponse<Product[]>(
        axios.get(`${BASE_URL}/products`)
    );
};

export const getProductById = async (id: number): Promise<Product> => {
    return handleResponse<Product>(
        axios.get(`${BASE_URL}/products/${id}`)
    );
};

type UpsertProductDTO = Omit<Product, "id" | "createdAt" | "updatedAt">;

export const addProduct = async (
    data: UpsertProductDTO
): Promise<Product> => {
    return handleResponse<Product>(
        axios.post(`${BASE_URL}/products`, data, {
            headers: {
                "Content-Type": "application/json",
            },
        })
    );
};

export const updateProduct = async (
    id: number,
    data: UpsertProductDTO
): Promise<Product> => {
    return handleResponse<Product>(
        axios.put(`${BASE_URL}/products/${id}`, data, {
            headers: {
                "Content-Type": "application/json",
            },
        })
    );
};

export const deleteProduct = async (id: number): Promise<void> => {
    await handleResponse(
        axios.delete(`${BASE_URL}/products/${id}`)
    );
};