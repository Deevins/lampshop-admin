export interface Product {
    id: number;
    sku: string;
    name: string;
    description: string;
    categoryId: string;
    isActive: boolean;
    imageUrl: string;
    price: number;
    stockQty: number;
    attributes: Record<string, string | number>;
    createdAt?: string;
    updatedAt?: string;
}


export interface Category {
    id: string;
    name: string;
}