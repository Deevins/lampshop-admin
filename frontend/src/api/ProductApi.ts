import type {Product} from "../types/Product.ts";
import {delay} from "../utils/delay.ts";


let products: Product[] = [
    {
        id: 1,
        name: "EcoBright 7W",
        price: 500,
        description: "Энергоэффективная лампочка для дома.",
        power: 7,
        color: "Тёплый белый",
        temperature: 2700,
        socketType: "E27",
        imageUrl:
            "https://via.placeholder.com/100?text=EcoBright+7W"
    }
    // You can prefill more if you want.
];

export const getProducts = async (): Promise<Product[]> => {
    await delay(300);
    return [...products];
};

export const getProductById = async (id: number): Promise<Product | null> => {
    await delay(200);
    const p = products.find((prod) => prod.id === id);
    return p ? {...p} : null;
};

export const addProduct = async (data: Omit<Product, "id">): Promise<Product> => {
    await delay(300);
    const newId = products.length > 0 ? Math.max(...products.map((p) => p.id)) + 1 : 1;
    const newProduct: Product = {id: newId, ...data};
    products.push(newProduct);
    return {...newProduct};
};

export const updateProduct = async (id: number, data: Omit<Product, "id">): Promise<Product | null> => {
    await delay(300);
    const idx = products.findIndex((p) => p.id === id);
    if (idx === -1) return null;
    products[idx] = {id, ...data};
    return {...products[idx]};
};

export const deleteProduct = async (id: number): Promise<boolean> => {
    await delay(200);
    const idx = products.findIndex((p) => p.id === id);
    if (idx === -1) return false;
    products.splice(idx, 1);
    return true;
};