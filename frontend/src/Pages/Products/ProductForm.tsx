import React, {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";

import styles from "./ProductForm.module.scss";
import {addProduct, getProductById, updateProduct} from "../../api/ProductApi.ts";
import type {Category} from "../../types/Product.ts";
import type {AttributeOption} from "../../api/api.ts";
import {getAttributeOptions, getCategories} from "../../api/OrderApi.ts";


interface FormState {
    sku: string;
    name: string;
    description: string;
    categoryId: string;
    price: number;
    stockQty: number;
    imageUrl: string;
    isActive: boolean;
    attributes: Record<string, string | number>;
}

const ProductForm: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const isEditMode = Boolean(id);
    const navigate = useNavigate();

    // Состояния для динамических данных
    const [categories, setCategories] = useState<Category[]>([]);
    const [attributeOptions, setAttributeOptions] = useState<AttributeOption[]>([]);

    // Состояние формы
    const [formState, setFormState] = useState<FormState>({
        sku: "",
        name: "",
        description: "",
        categoryId: "",
        price: 0,
        stockQty: 0,
        imageUrl: "",
        isActive: true,
        attributes: {},
    });

    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    // Флаг, что данные товара (в режиме редактирования) уже загружены
    const [productLoaded, setProductLoaded] = useState<boolean>(false);

    // 1) Загрузить категории при монтировании
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

    // 2) Если в режиме редактирования, загрузить товар по ID
    useEffect(() => {
        if (isEditMode && id) {
            setLoading(true);
            getProductById(Number(id))
                .then((prod) => {
                    // Заполняем все поля формы
                    setFormState({
                        sku: prod.sku,
                        name: prod.name,
                        description: prod.description,
                        categoryId: prod.categoryId,
                        price: prod.price,
                        stockQty: prod.stockQty,
                        imageUrl: prod.imageUrl,
                        isActive: prod.isActive,
                        attributes: { ...prod.attributes }, // сохраняем атрибуты из ответа
                    });
                    setProductLoaded(true);
                    // Подгружаем опции атрибутов для данной категории
                    return getAttributeOptions(prod.categoryId);
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

    // 3) Если не в режиме редактирования, или если пользователь сменил категорию после создания,
    //    нужно подгрузить опции и инициализировать атрибуты по умолчанию.
    useEffect(() => {
        // Случай: режим редактирования. Тогда мы не хотим перезаписывать уже загруженные атрибуты.
        // Но если товар ещё не загружен, это значит, что категория может меняться (нет исходных attrs) —
        // тогда мы выполним инициализацию.
        if (isEditMode) {
            if (!productLoaded && formState.categoryId) {
                // Пока товар не загружен, но categoryId мог быть из начального состояния (""),
                // и юзер может его менять до загрузки. Впрочем, обычно forced load происходит сразу.
                const cid = formState.categoryId;
                getAttributeOptions(cid)
                    .then((opts) => {
                        setAttributeOptions(opts);
                        // Инициализируем пустые значения (пока)
                        const newAttrs: Record<string, string | number> = {};
                        opts.forEach((attr) => {
                            newAttrs[attr.key] = attr.type === "number" ? 0 : "";
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
            }
            return;
        }

        // Случай: режим "создания" (не редактирование) или товар загружен и пользователь изменил категорию вручную.
        if (formState.categoryId) {
            setLoading(true);
            getAttributeOptions(formState.categoryId)
                .then((opts) => {
                    setAttributeOptions(opts);
                    // Инициализируем пустые значения атрибутов
                    const newAttrs: Record<string, string | number> = {};
                    opts.forEach((attr) => {
                        newAttrs[attr.key] = attr.type === "number" ? 0 : "";
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
    }, [formState.categoryId, isEditMode, productLoaded]);

    // 4) Обработчик для простых полей: <input>, <textarea>, <select>
    const handleChangeBasic = (
        e: React.ChangeEvent<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >
    ) => {
        const target = e.target as HTMLInputElement;
        const { name, type, value, checked } = target;

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

    // 5) Обработчик изменения динамических атрибутов
    const handleChangeAttribute = (key: string, value: string | number) => {
        setFormState((prev) => ({
            ...prev,
            attributes: {
                ...prev.attributes,
                [key]: value,
            },
        }));
    };

    // 6) Отправка формы (создание или обновление)
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        if (!formState.categoryId) {
            setError("Пожалуйста, выберите категорию");
            return;
        }

        const dto = {
            sku: formState.sku,
            name: formState.name,
            description: formState.description,
            categoryId: formState.categoryId,
            price: formState.price,
            stockQty: formState.stockQty,
            imageUrl: formState.imageUrl,
            isActive: formState.isActive,
            attributes: formState.attributes,
        };

        setLoading(true);
        try {
            if (isEditMode && id) {
                await updateProduct(Number(id), dto);
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
                    {/* ===== Артикул (SKU) ===== */}
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

                    {/* ===== Название ===== */}
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

                    {/* ===== Описание ===== */}
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

                    {/* ===== Категория ===== */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="categoryId">
                            Категория
                        </label>
                        <select
                            id="categoryId"
                            name="categoryId"
                            value={formState.categoryId}
                            onChange={handleChangeBasic}
                            className={styles["form-select"]}
                            required
                        >
                            <option value="">— Выберите категорию —</option>
                            {categories.map((cat) => (
                                <option key={cat.id} value={cat.id}>
                                    {cat.name}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* ===== Цена ===== */}
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

                    {/* ===== Количество на складе ===== */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="stockQty">
                            Количество на складе
                        </label>
                        <input
                            id="stockQty"
                            name="stockQty"
                            type="number"
                            value={formState.stockQty}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                            min={0}
                            required
                        />
                    </div>

                    {/* ===== Ссылка на изображение ===== */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]} htmlFor="imageUrl">
                            Ссылка на изображение (URL)
                        </label>
                        <input
                            id="imageUrl"
                            name="imageUrl"
                            type="text"
                            value={formState.imageUrl}
                            onChange={handleChangeBasic}
                            className={styles["form-input"]}
                        />
                    </div>

                    {/* ===== Активность товара ===== */}
                    <div className={styles["form-group"]}>
                        <label className={styles["form-label"]}>
                            <input
                                type="checkbox"
                                name="isActive"
                                checked={formState.isActive}
                                onChange={handleChangeBasic}
                            />{" "}
                            Активный товар
                        </label>
                    </div>

                    {/* ===== Динамические атрибуты ===== */}
                    {formState.categoryId && (
                        <div className={styles["attributes-section"]}>
                            <h3>Атрибуты для категории:</h3>
                            {attributeOptions.length > 0 ? (
                                attributeOptions.map((opt) => (
                                    <div
                                        key={opt.key}
                                        className={styles["attribute-row"]}
                                    >
                                        <label htmlFor={opt.key}>{opt.label}</label>
                                        <input
                                            id={opt.key}
                                            type={opt.type === "number" ? "number" : "text"}
                                            value={
                                                formState.attributes[opt.key] as string | number
                                            }
                                            onChange={(e) => {
                                                const val =
                                                    opt.type === "number"
                                                        ? Number(e.target.value)
                                                        : e.target.value;
                                                handleChangeAttribute(opt.key, val);
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
                        <p style={{ color: "#d9534f", marginTop: 10 }}>
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
