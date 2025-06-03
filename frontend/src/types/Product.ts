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
    category_id: string;
    is_active: boolean;
    image_url: string;
    price: number;
    stock_qty: number;
    attributes: Record<string, string | number>;
    createdAt?: string;
    updatedAt?: string;
}