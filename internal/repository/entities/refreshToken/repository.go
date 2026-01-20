package refreshToken

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Save(refreshtoken *entities.RefreshToken) error
	FindByID(ID *uint64, refreshToken *entities.RefreshToken) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Save(refreshtoken *entities.RefreshToken) error {
	return r.DB.Save(refreshtoken).Error
}

func (r *repository) FindByID(ID *uint64, refreshToken *entities.RefreshToken) error {
	return r.DB.First(refreshToken, "id = ?", ID).Error
}