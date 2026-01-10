package dtos

type CreateProductRequest struct {
	Name *string `json:"name"`
	Price *string `json:"price"`
	Description *string `json:"description"`
}

type UpdateProductRequest struct {
	ProductID *uint32 `json:"productID"`
	Name *string `json:"name"`
	Price *string `json:"price"`
	Description *string `json:"description"`
}

type UpdateProductResponse struct {
	Name *string
	Price *string
	Description *string
}