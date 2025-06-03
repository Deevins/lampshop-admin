package infra

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// OrderServiceClient отвечает за взаимодействие с сервисом заказов.
type OrderServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewOrderServiceClient создает клиента для сервиса заказов.
// baseURL — это адрес, по которому доступен order-service, например "http://localhost:9000".
func NewOrderServiceClient(baseURL string) *OrderServiceClient {
	return &OrderServiceClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *OrderServiceClient) GetAllOrders(token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/orders", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *OrderServiceClient) GetActiveOrders(token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/orders/active", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *OrderServiceClient) GetOrderStatus(token, id string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/orders/"+id+"/status", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *OrderServiceClient) CreateOrder(token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+"/checkout", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *OrderServiceClient) UpdateOrderStatus(token string, id string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", c.BaseURL+"/orders/"+id+"/status", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

// ProductServiceClient отвечает за взаимодействие с сервисом товаров.
type ProductServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewProductServiceClient создает клиента для сервиса товаров.
// baseURL — это адрес, по которому доступен product-service, например "http://localhost:9100".
func NewProductServiceClient(baseURL string) *ProductServiceClient {
	return &ProductServiceClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *ProductServiceClient) ListProducts(token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/products", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) CreateProduct(token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+"/products", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) GetProductByID(token, id string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/products/"+id, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) UpdateProduct(token, id string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", c.BaseURL+"/products/"+id, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) DeleteProduct(token, id string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+"/products/"+id, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) GetUploadURL(token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/products/upload-url", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

func (c *ProductServiceClient) NotifyUpload(token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+"/products/notify-upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTPClient.Do(req)
}

// Помощник: копирует статус, заголовки и тело ответа от внешнего сервиса в Gin-контекст.
func ProxyResponse(c *gin.Context, resp *http.Response, err error) {
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Копируем статус-код
	c.Status(resp.StatusCode)
	// Копируем заголовки
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}
	// Копируем тело
	bodyBuf := new(bytes.Buffer)
	_, copyErr := io.Copy(bodyBuf, resp.Body)
	if copyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": copyErr.Error()})
		return
	}
	c.Writer.Write(bodyBuf.Bytes())
}
