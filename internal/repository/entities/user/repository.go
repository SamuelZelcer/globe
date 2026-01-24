package user

import (
	"globe/internal/repository/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user *entities.User) error
	FindByEmail(email string, user *entities.User) error
	FindByID(ID uint64, user *entities.User) error
	IsUsernameAlreadyInUse(username string) (bool, error)
	IsEmailAlreadyInUse(email string) (bool, error)
	FindUserByIDWithAllHisProducts(ID uint64, user *entities.User) error
	UpdateEmailByID(ID uint64, email string) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Save(user *entities.User) error {
	return r.DB.Save(user).Error
}

func (r *repository) FindByEmail(email string, user *entities.User) error {
	return r.DB.First(user, "email = ?", email).Error
}

func (r *repository) FindByID(ID uint64, user *entities.User) error {
	return r.DB.First(user, ID).Error
}

func (r *repository) IsUsernameAlreadyInUse(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *repository) IsEmailAlreadyInUse(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *repository) FindUserByIDWithAllHisProducts(ID uint64, user *entities.User) error {
	return r.DB.Preload("Products").First(user, ID).Error
}

func (r *repository) UpdateEmailByID(ID uint64, email string) error {
	return r.DB.Model(&entities.User{}).Where(ID).Update("email", email).Error
}