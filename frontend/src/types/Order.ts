export interface OrderItem {
    product_id: number;
    quantity: number;
}

export type OrderStatus = "pending" | "processing" | "shipped" | "delivered";

export interface Order {
    id: number;
    full_name: string;
    items: OrderItem[];
    total: number;
    status: OrderStatus;
    createdAt?: string;
    updatedAt?: string;
}