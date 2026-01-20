package product

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Save(product *entities.Product) (*uint64, error)
	FindByID(ID *uint64, product *entities.Product) error
	DeleteByID(ID *uint64) error
	FindProductsForSearch(name *string, offset int, products *[]entities.Product) error
	CountProducts(name *string, amount *int64) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Save(product *entities.Product) (*uint64, error) {
	if err := r.DB.Save(product).Error; err != nil {
		return nil, err
	}
	return &product.ID, nil
}

func (r *repository) FindByID(ID *uint64, product *entities.Product) error {
	return r.DB.First(product, "id = ?", ID).Preload("User").Error
}

func (r *repository) DeleteByID(ID *uint64) error {
	return r.DB.Unscoped().Delete(&entities.Product{}, "id = ?", ID).Error
}

func (r *repository) FindProductsForSearch(name *string, offset int, products *[]entities.Product) error {
	return r.DB.Where("name % ?", name).
	Preload("User").
	Order("id desc").
    Limit(15).
	Offset(offset).
    Find(products).
	Error
}

func (r *repository) CountProducts(name *string, amount *int64) error {
	return r.DB.Table("products").
        Where("name % ?", name).
        Count(amount).
		Error
}