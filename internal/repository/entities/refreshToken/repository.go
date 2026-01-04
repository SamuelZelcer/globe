package refreshToken

import "gorm.io/gorm"

type Repository interface {
	Create(refreshtoken *RefreshToken) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(refreshtoken *RefreshToken) error {
	return r.DB.Save(refreshtoken).Error
}