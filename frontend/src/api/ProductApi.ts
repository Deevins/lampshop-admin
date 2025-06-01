import type {Product} from "../types/Product.ts";
import {delay} from "../utils/delay.ts";


let products: Product[] = [
    {
        id: 1,
        sku: "BULB-007",
        name: "EcoBright 7W",
        description: "Энергоэффективная лампочка для дома.",
        categoryId: "bulb",
        isActive: true,
        imageUrl: "https://via.placeholder.com/100?text=EcoBright+7W",
        price: 500,
        stockQty: 20,
        attributes: {
            power: 7,
            color: "Тёплый белый",
            temperature: 2700,
            socketType: "E27"
        },
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
    }
];

export const getProducts = async (): Promise<Product[]> => {
    await delay(200);
    return products.map((p) => ({...p, attributes: {...p.attributes}}));
};

export const getProductById = async (id: number): Promise<Product | null> => {
    await delay(150);
    const found = products.find((p) => p.id === id);
    return found ? {...found, attributes: {...found.attributes}} : null;
};

type UpsertProductDTO = Omit<Product, "id" | "createdAt" | "updatedAt">;

export const addProduct = async (
    data: UpsertProductDTO
): Promise<Product> => {
    await delay(200);
    const newId =
        products.length > 0 ? Math.max(...products.map((p) => p.id)) + 1 : 1;
    const now = new Date().toISOString();
    const newProd: Product = {
        id: newId,
        ...data,
        createdAt: now,
        updatedAt: now
    };
    products.push(newProd);
    return {...newProd, attributes: {...newProd.attributes}};
};

export const updateProduct = async (
    id: number,
    data: UpsertProductDTO
): Promise<Product | null> => {
    await delay(200);
    const idx = products.findIndex((p) => p.id === id);
    if (idx === -1) return null;
    products[idx] = {
        ...products[idx],
        ...data,
        updatedAt: new Date().toISOString()
    };
    return {...products[idx], attributes: {...products[idx].attributes}};
};

export const deleteProduct = async (id: number): Promise<boolean> => {
    await delay(150);
    const idx = products.findIndex((p) => p.id === id);
    if (idx === -1) return false;
    products.splice(idx, 1);
    return true;
};