package product

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Create(product *entities.Product) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(product *entities.Product) error {
	return r.DB.Create(product).Error
}