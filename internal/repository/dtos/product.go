package dtos

type CreateProductRequest struct {
	Name string `json:"name"`
	Price string `json:"price"`
	Description string `json:"description"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateProductRequest struct {
	ProductID uint64 `json:"productID"`
	Name string `json:"name"`
	Price string `json:"price"`
	Description string `json:"description"`
	RefreshToken string `json:"refreshToken"`
}

type DeleteProductRequest struct {
	ProductID uint64 `json:"productID"`
	RefreshToken string `json:"refreshToken"`
}

type SearchRequest struct {
	Name string `json:"name"`
	Page uint32 `json:"page"`
}

type UpdateProductResponce struct {
	Name string `json:"name"`
	Price string `json:"price"`
	Description string `json:"description"`
}

type SearchProduct struct {
	ProductID uint64 `json:"productID"`
	Name string `json:"name"`
	Price string `json:"price"`
}

type SearchProductResponse struct {
	TotalAmountOfProducts int64 `json:"totalAmountOfProducts"`
	TotalAmountOfPages int64 `json:"totalAmountOfPages"`
	CurrentPage uint32 `json:"currentPage"`
	Products *[]SearchProduct `json:"products"`
}