import React, {useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";

import styles from "./ProductList.module.scss";
import {deleteProduct, getProducts} from "../../api/ProductApi.ts";
import type {Category, Product} from "../../types/Product.ts";
import {getCategories} from "../../api/OrderApi.ts";

const ProductList: React.FC = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const navigate = useNavigate();

    const fetchProducts = async () => {
        try {
            const data = await getProducts();
            setProducts(data);
        } catch (err) {
            console.error("Error fetching products:", err);
        }
    };

    const fetchCategories = async () => {
        try {
            const cats = await getCategories();
            setCategories(cats);
        } catch (err) {
            console.error("Error fetching categories:", err);
        }
    };

    useEffect(() => {
        fetchCategories();
        fetchProducts();
    }, []);

    const handleDelete = async (id: number) => {
        if (!confirm("Вы уверены, что хотите удалить этот товар?")) return;

        try {
            await deleteProduct(id);
            fetchProducts();
        } catch (err) {
            console.error("Error deleting product:", err);
            alert("Не удалось удалить товар");
        }
    };

    const getCategoryName = (category_id: string): string => {
        const cat = categories.find((c) => c.id === category_id);
        return cat ? cat.name : category_id;
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
                    <th className={styles.th}>SKU</th>
                    <th className={styles.th}>Название</th>
                    <th className={styles.th}>Категория</th>
                    <th className={styles.th}>Цена (₽)</th>
                    <th className={styles.th}>В наличии</th>
                    <th className={styles.th}>Изображение</th>
                    <th className={styles.th}>Атрибуты (JSON)</th>
                    <th className={styles.th}>Действия</th>
                </tr>
                </thead>
                <tbody>
                {products.map((prod) => (
                    <tr key={prod.id}>
                        <td className={styles.td}>{prod.id}</td>
                        <td className={styles.td}>{prod.sku}</td>
                        <td className={styles.td}>{prod.name}</td>
                        <td className={styles.td}>
                            {getCategoryName(prod.category_id)}
                        </td>
                        <td className={styles.td}>{prod.price}</td>
                        <td className={styles.td}>{prod.stock_qty}</td>
                        <td className={styles.td}>
                            <img
                                src={'https://santhimetaleshop.in/cdn/shop/files/Untitleddesign_26a5d7f4-82b7-4e7a-ac43-068a31086beb.png?v=1694498000\u0026width=1445'}
                                alt={prod.name}
                                width={60}
                                height={60}
                            />
                        </td>
                        <td className={styles.td}>
                <pre style={{fontSize: 12, whiteSpace: "pre-wrap"}}>
                  {JSON.stringify(prod.attributes, null, 2)}
                </pre>
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
