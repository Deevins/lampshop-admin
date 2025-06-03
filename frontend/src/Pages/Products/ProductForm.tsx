import React, {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";

import styles from "./ProductForm.module.scss";
import {addProduct, getAttributeOptions, getProductById, updateProduct} from "../../api/ProductApi.ts";
import type {AttributeOption, Category} from "../../types/Product.ts";
import { getCategories} from "../../api/OrderApi.ts";


interface FormState {
    sku: string;
    name: string;
    description: string;
    category_id: string;
    price: number;
    stock_qty: number;
    image_url: string;
    is_active: boolean;
    attributes: Record<string, string | number>;
}

const ProductForm: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const isEditMode = Boolean(id);
    const navigate = useNavigate();

    const [categories, setCategories] = useState<Category[]>([]);
    const [attributeOptions, setAttributeOptions] = useState<
        AttributeOption[]
    >([]);

    const [formState, setFormState] = useState<FormState>({
        sku: "",
        name: "",
        description: "",
        category_id: "",
        price: 0,
        stock_qty: 0,
        image_url: "",
        is_active: true,
        attributes: {},
    });

    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [productLoaded, setProductLoaded] = useState<boolean>(false);

    useEffect(() => {
        const fetchCats = async () => {
            try {
                const cats = await getCategories();
                setCategories(cats);
            } catch (err) {
                console.error("Error fetching categories:", err);
            }
        };
        fetchCats();
    }, []);

    useEffect(() => {
        if (isEditMode && id) {
            setLoading(true);
            getProductById(id)
                .then((prod) => {
                    setFormState({
                        sku: prod.sku,
                        name: prod.name,
                        description: prod.description,
                        category_id: prod.category_id,
                        price: prod.price,
                        stock_qty: prod.stock_qty,
                        image_url: prod.image_url,
                        is_active: prod.is_active,
                        attributes: {...prod.attributes},
                    });
                    setProductLoaded(true);
                    return getAttributeOptions(prod.category_id);
                })
                .then((opts) => {
                    setAttributeOptions(opts);
                })
                .catch((err) => {
                    console.error("Error loading product:", err);
                    alert("Не удалось загрузить данные товара");
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [id, isEditMode]);

    useEffect(() => {
        // Если режим редактирования, но товар ещё не загружен (productLoaded=false),
        // то нужно всё равно подгрузить атрибуты для выбранной категории (если юзер сменил dropdown раньше, чем загрузился товар)
        if (isEditMode && !productLoaded && formState.category_id) {
            const cid = formState.category_id;
            getAttributeOptions(cid)
                .then((opts) => {
                    setAttributeOptions(opts);
                    const newAttrs: Record<string, string | number> = {};
                    opts.forEach((attr) => {
                        newAttrs[attr.Key] = attr.Type === "number" ? 0 : "";
                    });
                    setFormState((prev) => ({
                        ...prev,
                        attributes: newAttrs,
                    }));
                })
                .catch((err) => {
                    console.error("Error fetching attribute options:", err);
                    setAttributeOptions([]);
                    setFormState((prev) => ({
                        ...prev,
                        attributes: {},
                    }));
                });
            return;
        }

        if (!isEditMode && formState.category_id) {
            setLoading(true);
            getAttributeOptions(formState.category_id)
                .then((opts) => {
                    setAttributeOptions(opts);
                    const newAttrs: Record<string, string | number> = {};
                    opts.forEach((attr) => {
                        newAttrs[attr.Key] = attr.Type === "number" ? 0 : "";
                    });
                    setFormState((prev) => ({
                        ...prev,
                        attributes: newAttrs,
                    }));
                })
                .catch((err) => {
                    console.error("Error fetching attribute options:", err);
                    setAttributeOptions([]);
                    setFormState((prev) => ({
                        ...prev,
                        attributes: {},
                    }));
                })
                .finally(() => setLoading(false));
        }
    }, [formState.category_id, isEditMode, productLoaded]);

    const handleChangeBasic = (
        e: React.ChangeEvent<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >
    ) => {
        const target = e.target as HTMLInputElement;
        const {name, type, value, checked} = target;
        setFormState((prev) => ({
            ...prev,
            [name]:
                type === "checkbox"
                    ? checked
                    : type === "number"
                        ? Number(value)
                        : value,
        }));
    };

    const handleChangeAttribute = (key: string, value: string | number) => {
        setFormState((prev) => ({
            ...prev,
            attributes: {
                ...prev.attributes,
                [key]: value,
            },
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        if (!formState.category_id) {
            setError("Пожалуйста, выберите категорию");
            return;
        }

        const dto = {
            sku: formState.sku,
            name: formState.name,
            description: formState.description,
            category_id: formState.category_id,
            price: formState.price,
            stock_qty: formState.stock_qty,
            image_url: formState.image_url,
            is_active: formState.is_active,
            attributes: formState.attributes,
        };

        setLoading(true);
        try {
            if (isEditMode && id) {
                await updateProduct(id, dto);
            } else {
                await addProduct(dto);
            }
            navigate("/products");
        } catch (err) {
            console.error("Error saving product:", err);
            setError("Ошибка при сохранении товара");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={styles["form-container"]}>
            <h2 className={styles.title}>
                {isEditMode ? "Редактировать товар" : "Добавить товар"}
            </h2>

            {loading && <p>Загрузка данных...</p>}

            {!loading && (
                <form onSubmit={handleSubmit}>
                    {/* SKU */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="sku">
                            Артикул (SKU)
                        </label>
                        <input
                            id="sku"
                            name="sku"
                            type="text"
                            value={formState.sku}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            required
                        />
                    </div>

                    {/* Название */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="name">
                            Название
                        </label>
                        <input
                            id="name"
                            name="name"
                            type="text"
                            value={formState.name}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            required
                        />
                    </div>

                    {/* Описание */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="description">
                            Описание
                        </label>
                        <textarea
                            id="description"
                            name="description"
                            value={formState.description}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            rows={2}
                        />
                    </div>

                    {/* Категория */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="category_id">
                            Категория
                        </label>
                        <select
                            id="category_id"
                            name="category_id"
                            value={formState.category_id}
                            onChange={handleChangeBasic}
                            className={styles["form-select"]}
                            required
                        >
                            <option value="">— Выберите категорию —</option>
                            {categories.map((cat) => (
                                <option key={cat.ID} value={cat.ID}>
                                    {cat.Name}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Цена */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="price">
                            Цена (₽)
                        </label>
                        <input
                            id="price"
                            name="price"
                            type="number"
                            value={formState.price}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            min={0}
                            step={0.01}
                            required
                        />
                    </div>

                    {/* Количество на складе */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="stock_qty">
                            Количество на складе
                        </label>
                        <input
                            id="stock_qty"
                            name="stock_qty"
                            type="number"
                            value={formState.stock_qty}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            min={0}
                            required
                        />
                    </div>

                    {/* Ссылка на изображение */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="image_url">
                            Ссылка на изображение (URL)
                        </label>
                        <input
                            id="image_url"
                            name="image_url"
                            type="text"
                            value={formState.image_url}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                        />
                    </div>

                    {/* Активность товара */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]}>
                            <input
                                type="checkbox"
                                name="is_active"
                                checked={formState.is_active}
                                onChange={handleChangeBasic}
                            />{" "}
                            Активный товар
                        </label>
                    </div>

                    {/* Динамические атрибуты */}
                    {formState.category_id && (
                        <div className={styles["attributes-section"]}>
                            <h3>Атрибуты для категории:</h3>
                            {attributeOptions.length > 0 ? (
                                attributeOptions.map((opt) => (
                                    <div
                                        key={opt.Key}
                                        className={styles["attribute-row"]}
                                    >
                                        <label htmlFor={opt.Key}>{opt.Label}</label>
                                        <input
                                            id={opt.Key}
                                            type={opt.Type === "number" ? "number" : "text"}
                                            value={
                                                formState.attributes[opt.Key] as string | number
                                            }
                                            onChange={(e) => {
                                                const val =
                                                    opt.Type === "number"
                                                        ? Number(e.target.value)
                                                        : e.target.value;
                                                handleChangeAttribute(opt.Key, val);
                                            }}
                                            className={styles["form-input"]}
                                        />
                                    </div>
                                ))
                            ) : (
                                <p>У этой категории нет дополнительных атрибутов</p>
                            )}
                        </div>
                    )}

                    {error && (
                        <p style={{color: "#d9534f", marginTop: 10}}>
                            {error}
                        </p>
                    )}

                    <button
                        type="submit"
                        className={styles["submit-button"]}
                        disabled={loading}
                    >
                        {isEditMode ? "Сохранить изменения" : "Добавить товар"}
                    </button>
                </form>
            )}
        </div>
    );
};

export default ProductForm;
