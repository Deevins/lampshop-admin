import React, {useEffect, useState} from "react";
import styles from "./OrderList.module.scss";
import {getOrders, updateOrderStatus} from "../../api/OrderApi.ts";
import type {Order, OrderStatus} from "../../types/Order.ts";

const statuses: OrderStatus[] = ["pending", "processing", "shipped", "delivered"];

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

    // 1) Функция для загрузки списка заказов
    const fetchOrders = async () => {
        try {
            const data = await getOrders();
            console.log("Получили с сервера (raw orders):", data);
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

    // 2) При маунте компонента загружаем заказы
    useEffect(() => {
        fetchOrders();
    }, []);

    // 3) Автоматически скрываем уведомление через 3 секунды
    useEffect(() => {
        if (notification.visible) {
            const timer = setTimeout(() => {
                setNotification((prev) => ({...prev, visible: false}));
            }, 3000);
            return () => clearTimeout(timer);
        }
    }, [notification.visible]);

    // 4) Обработчик смены статуса конкретного заказа
    const handleStatusChange = async (orderId: number, newStatus: OrderStatus) => {
        try {
            const updated = await updateOrderStatus(orderId, newStatus);
            if (updated) {
                setNotification({
                    message: `Статус заказа #${orderId} обновлён: ${newStatus}`,
                    type: "success",
                    visible: true,
                });
                fetchOrders(); // заново подтягиваем список
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
                {/* Если orders ещё не загружены или пустой массив, .map() отработает корректно (пустой tbody) */}
                {orders.map((order) => {
                    // На всякий случай возьмём items как пустой массив, если он undefined
                    const itemsArr = order.items ?? [];
                    console.log(itemsArr)
                    return (
                        <tr key={order.id}>
                            <td className={styles.td}>{order.id}</td>
                            <td className={styles.td}>{order.full_name}</td>
                            <td className={styles.td}>
                                {itemsArr.length > 0
                                    ? itemsArr.map((item) => `${item.product_id}:${item.quantity}`).join(", ")
                                    : "-"}
                            </td>
                            <td className={styles.td}>{order.total}</td>
                            <td className={styles.td}>
                                <select
                                    className={styles["status-select"]}
                                    value={order.status}
                                    onChange={(e) =>
                                        handleStatusChange(order.id, e.target.value as OrderStatus)
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
                    );
                })}
                </tbody>
            </table>

            {notification.visible && (
                <div
                    className={`${styles.notification} ${
                        notification.type === "success" ? styles.success : styles.error
                    }`}
                >
                    {notification.message}
                </div>
            )}
        </div>
    );
};

export default OrderList;
