import type {Category} from "../types/Product.ts";

export const categories: Category[] = [
    {id: "bulb", name: "Лампочки"},
    {id: "cable", name: "Кабели"},
    {id: "equipment", name: "Оборудование"}
];

export interface AttributeOption {
    key: string;
    label: string;
    type: "text" | "number";
}

export const attributeOptions: Record<string, AttributeOption[]> = {
    bulb: [
        {key: "power", label: "Мощность (Вт)", type: "number"},
        {key: "color", label: "Цвет", type: "text"},
        {key: "temperature", label: "Температура (K)", type: "number"},
        {key: "socketType", label: "Тип цоколя", type: "text"}
    ],
    cable: [
        {key: "length", label: "Длина (м)", type: "number"},
        {key: "material", label: "Материал", type: "text"},
        {key: "color", label: "Цвет", type: "text"}
    ],
    equipment: [
        {key: "manufacturer", label: "Производитель", type: "text"},
        {key: "model", label: "Модель", type: "text"},
        {key: "warranty", label: "Гарантия (мес.)", type: "number"}
    ]
};

