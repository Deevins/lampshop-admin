export interface Category {
    id: string;
    name: string;
}

export interface AttributeOption {
    key: string;
    label: string;
    type: "text" | "number";
}

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