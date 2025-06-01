import React, {useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";

import styles from "./ProductList.module.scss";
import {deleteProduct, getProducts} from "../../api/ProductApi.ts";
import type {Product} from "../../types/Product.ts";

const ProductList: React.FC = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const navigate = useNavigate();

    const fetchProducts = async () => {
        const data = await getProducts();
        setProducts(data);
    };

    useEffect(() => {
        fetchProducts();
    }, []);

    const handleDelete = async (id: number) => {
        if (!confirm("Вы уверены, что хотите удалить этот товар?")) return;
        await deleteProduct(id);
        fetchProducts();
    };

    return (
        <div className={styles["products-container"]}>
            <div className={styles["header-row"]}>
                <h2 className={styles.title}>Управление товарами</h2>
                <button
                    className={styles["add-button"]}
                    onClick={() => navigate("/products/add")}
                >
                    Добавить товар
                </button>
            </div>

            <table className={styles.table}>
                <thead>
                <tr>
                    <th className={styles.th}>ID</th>
                    <th className={styles.th}>Название</th>
                    <th className={styles.th}>Цена (₽)</th>
                    <th className={styles.th}>Описание</th>
                    <th className={styles.th}>Мощность (Вт)</th>
                    <th className={styles.th}>Цвет</th>
                    <th className={styles.th}>Температура (K)</th>
                    <th className={styles.th}>Тип цоколя</th>
                    <th className={styles.th}>Картинка</th>
                    <th className={styles.th}>Действия</th>
                </tr>
                </thead>
                <tbody>
                {products.map((prod) => (
                    <tr key={prod.id}>
                        <td className={styles.td}>{prod.id}</td>
                        <td className={styles.td}>{prod.name}</td>
                        <td className={styles.td}>{prod.price}</td>
                        <td className={styles.td}>{prod.description}</td>
                        <td className={styles.td}>{prod.power}</td>
                        <td className={styles.td}>{prod.color}</td>
                        <td className={styles.td}>{prod.temperature}</td>
                        <td className={styles.td}>{prod.socketType}</td>
                        <td className={styles.td}>
                            <img
                                src={prod.imageUrl}
                                alt={prod.name}
                                width={60}
                                height={60}
                            />
                        </td>
                        <td className={styles.td}>
                            <div className={styles["action-buttons"]}>
                                <button
                                    className={styles["edit-button"]}
                                    onClick={() =>
                                        navigate(`/products/edit/${prod.id}`)
                                    }
                                >
                                    Редактировать
                                </button>
                                <button
                                    className={styles["delete-button"]}
                                    onClick={() => handleDelete(prod.id)}
                                >
                                    Удалить
                                </button>
                            </div>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default ProductList;
