package dtos

type CreateProductRequest struct {
	Name *string `json:"name"`
	Price *string `json:"price"`
	Description *string `json:"description"`
}