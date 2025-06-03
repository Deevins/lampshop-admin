-- +goose NO TRANSACTION
-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
                                     username    TEXT        NOT NULL PRIMARY KEY,
                                     hashed_password TEXT   NOT NULL,
                                     created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO users (username, hashed_password)
VALUES
    ('admin', 'qwe')
ON CONFLICT DO NOTHING;


-- ====================
-- Таблица опций атрибутов (attribute_options)
-- Каждая категория может иметь несколько опций.
-- key — внутреннее имя (например, power, color и т. д.)
-- label — человекочитаемое название (Мощность (Вт), Цвет и т. п.)
-- type — "text" или "number"
-- ====================

CREATE TABLE IF NOT EXISTS categories (
                                          id   UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
                                          name TEXT        NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS attribute_options (
                                                 id          SERIAL PRIMARY KEY,
                                                 category_id UUID NOT NULL,
                                                 key         TEXT NOT NULL,
                                                 label       TEXT NOT NULL,
                                                 type        TEXT NOT NULL CHECK (type IN ('text', 'number')),
                                                 UNIQUE (category_id, key)
);

-- Примеры initial-данных:
-- INSERT INTO attribute_options (category_id, key, label, type)
-- uuid.New, 'power', 'Мощность (Вт)', 'number'  FROM categories c WHERE c.name = 'Лампочки' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'color', 'Цвет', 'text'        FROM categories c WHERE c.name = 'Лампочки' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'temperature', 'Температура (K)', 'number' FROM categories c WHERE c.name = 'Лампочки' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'socketType', 'Тип цоколя', 'text' FROM categories c WHERE c.name = 'Лампочки' ON CONFLICT DO NOTHING;
--
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'length', 'Длина (м)', 'number'    FROM categories c WHERE c.name = 'Кабели' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'material', 'Материал', 'text'     FROM categories c WHERE c.name = 'Кабели' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'color', 'Цвет', 'text'            FROM categories c WHERE c.name = 'Кабели' ON CONFLICT DO NOTHING;
--
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'manufacturer', 'Производитель', 'text' FROM categories c WHERE c.name = 'Оборудование' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'model', 'Модель', 'text'            FROM categories c WHERE c.name = 'Оборудование' ON CONFLICT DO NOTHING;
-- INSERT INTO attribute_options (category_id, key, label, type)
-- SELECT c.id, 'warranty', 'Гарантия (мес.)', 'number' FROM categories c WHERE c.name = 'Оборудование' ON CONFLICT DO NOTHING;



-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS attribute_options;