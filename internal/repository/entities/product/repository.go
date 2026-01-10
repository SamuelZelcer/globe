package product

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Save(product *entities.Product) error
	FindByID(ID *uint32, product *entities.Product) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Save(product *entities.Product) error {
	return r.DB.Save(product).Error
}

func (r *repository) FindByID(ID *uint32, product *entities.Product) error {
	return r.DB.Find(product, "id = ?", ID).Error
}