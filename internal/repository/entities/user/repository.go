package user

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *entities.User) error
	FindByEmail(email *string, user *entities.User) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(user *entities.User) error {
	return r.DB.Create(user).Error
}

func (r *repository) FindByEmail(email *string, user *entities.User) error {
	return r.DB.Find(user, "email = ?", email).Error
}