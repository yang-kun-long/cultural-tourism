// File: services/product_service.go
package services

import (
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

// Collection name for Products
const CollectionProduct = "product"

// CreateProduct creates a new product
func CreateProduct(product *models.Product) (map[string]interface{}, error) {
	// [Security] 强制初始化字段，防止恶意篡改
	product.ID = ""

	product.CreatedAt = time.Now().Format(time.RFC3339)
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	// 防止恶意写入非预期字段（业务兜底）
	if product.Price <= 0 {
		product.Price = 0
	}

	return tcb.Client.CreateData(CollectionProduct, product)
}

// ListProducts retrieves product list with pagination
func ListProducts(query models.ProductQuery) (map[string]interface{}, error) {
	where := make(map[string]interface{})

	// 状态筛选 - 默认只返回上线状态
	where["status"] = map[string]interface{}{"$eq": 1}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionProduct, filter, query.Page, query.Size)
}

// GetProductDetail retrieves a single product by ID
func GetProductDetail(id string) (map[string]interface{}, error) {
	return tcb.Client.GetDetail(CollectionProduct, id)
}

// UpdateProduct updates an existing product
func UpdateProduct(id string, product *models.Product) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if product.Name != "" {
		updateData["name"] = product.Name
	}
	if product.Image != "" {
		updateData["image"] = product.Image
	}
	if product.Price != 0 {
		updateData["price"] = product.Price
	}
	if product.JumpAppID != "" {
		updateData["jump_app_id"] = product.JumpAppID
	}
	if product.JumpPath != "" {
		updateData["jump_path"] = product.JumpPath
	}


	return tcb.Client.UpdateData(CollectionProduct, id, updateData)
}

// DeleteProduct deletes a product by ID
func DeleteProduct(id string) error {
	return tcb.Client.DeleteData(CollectionProduct, id)
}

// BatchUpdateProductStatus batch updates product status (for admin operations)
func BatchUpdateProductStatus(ids []string, status int) error {
	for _, id := range ids {
		updateData := map[string]interface{}{

			"updated_at": time.Now().Format(time.RFC3339),
		}
		if err := tcb.Client.UpdateData(CollectionProduct, id, updateData); err != nil {
			return err
		}
	}
	return nil
}

// CountProducts counts products (for statistics)
func CountProducts() (map[string]interface{}, error) {
	where := map[string]interface{}{
		"status": map[string]interface{}{"$eq": 1},
	}

	filter := map[string]interface{}{
		"where": where,
		"count": true,
	}

	return tcb.Client.ListData(CollectionProduct, filter, 1, 1)
}