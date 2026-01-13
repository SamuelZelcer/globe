package dtos

type CreateProductRequest struct {
	Name *string `json:"name"`
	Price *string `json:"price"`
	Description *string `json:"description"`
	RefreshToken *string
}

type UpdateProductRequest struct {
	ProductID *uint64 `json:"productID"`
	Name *string `json:"name"`
	Price *string `json:"price"`
	Description *string `json:"description"`
	RefreshToken *string
}

type UpdateProductResponse struct {
	Name *string
	Price *string
	Description *string
}

type DeleteProductRequest struct {
	ProductID *uint64 `json:"productID"`
	RefreshToken *string
}