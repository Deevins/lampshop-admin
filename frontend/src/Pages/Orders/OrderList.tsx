import React, {useEffect, useState} from "react"
import styles from "./OrderList.module.scss";
import {getOrders, updateOrderStatus} from "../../api/OrderApi.ts";
import type {Order, OrderStatus} from "../../types/Order.ts";

const statuses: OrderStatus[] = [
    "Pending",
    "Processing",
    "Shipped",
    "Delivered",
];

interface Notification {
    message: string;
    type: "success" | "error";
    visible: boolean;
}

const OrderList: React.FC = () => {
    const [orders, setOrders] = useState<Order[]>([]);
    const [notification, setNotification] = useState<Notification>({
        message: "",
        type: "success",
        visible: false,
    });

    const fetchOrders = async () => {
        try {
            const data = await getOrders();
            setOrders(data);
        } catch (err) {
            console.error("Error fetching orders:", err);
            setNotification({
                message: "Ошибка при загрузке заказов",
                type: "error",
                visible: true,
            });
        }
    };

    useEffect(() => {
        fetchOrders();
    }, []);

    useEffect(() => {
        if (notification.visible) {
            const timer = setTimeout(() => {
                setNotification((prev) => ({...prev, visible: false}));
            }, 3000);
            return () => clearTimeout(timer);
        }
    }, [notification.visible]);

    const handleStatusChange = async (
        orderId: number,
        newStatus: OrderStatus
    ) => {
        try {
            const updated = await updateOrderStatus(orderId, newStatus);
            if (updated) {
                setNotification({
                    message: `Статус заказа #${orderId} обновлён: ${newStatus}`,
                    type: "success",
                    visible: true,
                });
                fetchOrders();
            } else {
                setNotification({
                    message: `Не удалось обновить статус заказа #${orderId}`,
                    type: "error",
                    visible: true,
                });
            }
        } catch (err) {
            console.error("Error updating order status:", err);
            setNotification({
                message: `Ошибка при обновлении статуса #${orderId}`,
                type: "error",
                visible: true,
            });
        }
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

            {/* Всплывающее уведомление (тост) */}
            {notification.visible && (
                <div
                    className={`${styles.notification} ${
                        notification.type === "success"
                            ? styles.success
                            : styles.error
                    }`}
                >
                    {notification.message}
                </div>
            )}
        </div>
    );
};

export default OrderList;
