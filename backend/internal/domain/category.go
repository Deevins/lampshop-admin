package domain

// Category представляет запись о категории товаров.
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AttributeOption описывает один атрибут, который динамически выбирается
// для товара в зависимости от категории.
type AttributeOption struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"` // "text" или "number"
}
