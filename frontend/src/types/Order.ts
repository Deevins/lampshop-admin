export interface OrderItem {
    productId: number;
    quantity: number;
}

export type OrderStatus = "Pending" | "Processing" | "Shipped" | "Delivered";

export interface Order {
    id: number;
    customerName: string;
    items: OrderItem[];
    totalPrice: number;
    status: OrderStatus;
}