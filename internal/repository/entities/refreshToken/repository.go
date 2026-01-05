package refreshToken

import "gorm.io/gorm"

type Repository interface {
	Create(refreshtoken *RefreshToken) error
	FindByID(ID *uint32, refreshToken *RefreshToken) error
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

func (r *repository) FindByID(ID *uint32, refreshToken *RefreshToken) error {
	return r.DB.Find(refreshToken, "id = ?", ID).Error
}