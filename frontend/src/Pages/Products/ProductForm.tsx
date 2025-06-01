import React, {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import styles from "./ProductForm.module.scss";
import type {Product} from "../../types/Product.ts";
import {addProduct, getProductById, updateProduct} from "../../api/ProductApi.ts";

const ProductForm: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const isEditMode = Boolean(id);
    const navigate = useNavigate();

    const [formState, setFormState] = useState<Omit<Product, "id">>({
        name: "",
        price: 0,
        description: "",
        power: 0,
        color: "",
        temperature: 0,
        socketType: "",
        imageUrl: ""
    });

    useEffect(() => {
        if (isEditMode && id) {
            (async () => {
                const prod = await getProductById(Number(id));
                if (prod) {
                    const {id: _, ...rest} = prod;
                    setFormState(rest);
                }
            })();
        }
    }, [id]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => {
        const {name, value} = e.target;
        setFormState((prev) => ({
            ...prev,
            [name]:
                e.target.type === "number" ? Number(value) : value
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (isEditMode && id) {
            await updateProduct(Number(id), formState);
        } else {
            await addProduct(formState);
        }
        navigate("/products");
    };

    return (
        <div className={styles["form-container"]}>
            <h2 className={styles.title}>
                {isEditMode ? "Редактировать товар" : "Добавить товар"}
            </h2>
            <form onSubmit={handleSubmit}>
                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="name">
                        Название
                    </label>
                    <input
                        id="name"
                        type="text"
                        name="name"
                        value={formState.name}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="price">
                        Цена (₽)
                    </label>
                    <input
                        id="price"
                        type="number"
                        name="price"
                        value={formState.price}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="description">
                        Описание
                    </label>
                    <textarea
                        id="description"
                        name="description"
                        value={formState.description}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        rows={2}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="power">
                        Мощность (Вт)
                    </label>
                    <input
                        id="power"
                        type="number"
                        name="power"
                        value={formState.power}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="color">
                        Цвет
                    </label>
                    <input
                        id="color"
                        type="text"
                        name="color"
                        value={formState.color}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="temperature">
                        Температура (K)
                    </label>
                    <input
                        id="temperature"
                        type="number"
                        name="temperature"
                        value={formState.temperature}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="socketType">
                        Тип цоколя
                    </label>
                    <input
                        id="socketType"
                        type="text"
                        name="socketType"
                        value={formState.socketType}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <div className={styles["form-group"]}>
                    <label className={styles["form-label"]} htmlFor="imageUrl">
                        Ссылка на картинку (URL)
                    </label>
                    <input
                        id="imageUrl"
                        type="text"
                        name="imageUrl"
                        value={formState.imageUrl}
                        onChange={handleChange}
                        className={styles["form-input"]}
                        required
                    />
                </div>

                <button type="submit" className={styles["submit-button"]}>
                    {isEditMode ? "Сохранить изменения" : "Добавить товар"}
                </button>
            </form>
        </div>
    );
};

export default ProductForm;
