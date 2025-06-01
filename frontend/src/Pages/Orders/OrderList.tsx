import React, { useEffect, useState } from "react"
import styles from "./OrderList.module.scss";
import {getOrders, updateOrderStatus} from "../../api/OrderApi.ts";
import type {Order, OrderStatus} from "../../types/Order.ts";

const statuses: OrderStatus[] = [
    "Pending",
    "Processing",
    "Shipped",
    "Delivered"
];

const OrderList: React.FC = () => {
    const [orders, setOrders] = useState<Order[]>([]);

    const fetchOrders = async () => {
        const data = await getOrders();
        setOrders(data);
    };

    useEffect(() => {
        fetchOrders();
    }, []);

    const handleStatusChange = async (
        orderId: number,
        newStatus: OrderStatus
    ) => {
        await updateOrderStatus(orderId, newStatus);
        await fetchOrders();
    };

    return (
        <div className={styles["orders-container"]}>
            <h2 className={styles.title}>Просмотр заказов</h2>
            <table className={styles.table}>
                <thead>
                <tr>
                    <th className={styles.th}>ID заказа</th>
                    <th className={styles.th}>Клиент</th>
                    <th className={styles.th}>Товары (ID:кол-во)</th>
                    <th className={styles.th}>Сумма (₽)</th>
                    <th className={styles.th}>Статус</th>
                </tr>
                </thead>
                <tbody>
                {orders.map((order) => (
                    <tr key={order.id}>
                        <td className={styles.td}>{order.id}</td>
                        <td className={styles.td}>{order.customerName}</td>
                        <td className={styles.td}>
                            {order.items
                                .map((item) => `${item.productId}:${item.quantity}`)
                                .join(", ")}
                        </td>
                        <td className={styles.td}>{order.totalPrice}</td>
                        <td className={styles.td}>
                            <select
                                className={styles["status-select"]}
                                value={order.status}
                                onChange={(e) =>
                                    handleStatusChange(
                                        order.id,
                                        e.target.value as OrderStatus
                                    )
                                }
                            >
                                {statuses.map((s) => (
                                    <option key={s} value={s}>
                                        {s}
                                    </option>
                                ))}
                            </select>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default OrderList;
