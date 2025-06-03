export interface Category {
    ID: string;
    Name: string;
}

export interface AttributeOption {
    Key: string;
    Label: string;
    Type: "text" | "number";
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