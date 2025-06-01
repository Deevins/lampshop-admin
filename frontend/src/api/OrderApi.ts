import type {Order, OrderStatus} from "../types/Order.ts";
import {delay} from "../utils/delay.ts";



let orders: Order[] = [
    {
        id: 1,
        customerName: "Иван Иванов",
        items: [
            { productId: 1, quantity: 2 }
            // You can expand as you wish
        ],
        totalPrice: 1000,
        status: "Pending"
    }
    // More mock orders can be added.
];


export const getOrders = async (): Promise<Order[]> => {
    await delay(300);
    return [...orders];
};

export const updateOrderStatus = async (
    id: number,
    status: OrderStatus
): Promise<Order | null> => {
    await delay(200);
    const idx = orders.findIndex((o) => o.id === id);
    if (idx === -1) return null;
    orders[idx].status = status;
    return { ...orders[idx] };
};
