package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	IsUsernameOrEmailAlreadyInUse(username *string, email *string) (bool, error, error)
	FindUserByEmail(email *string, user *User) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *repository) IsUsernameOrEmailAlreadyInUse(username *string, email *string) (bool, error, error) {
	errUsername := r.DB.Where("username = ?", *username).First(&User{}).Error
	errEmail := r.DB.Where("email = ?", *email).First(&User{}).Error

	if errUsername == gorm.ErrRecordNotFound && errEmail == gorm.ErrRecordNotFound {
		return false, errUsername, errEmail
	}
	return true, errUsername, errEmail
}

func (r *repository) FindUserByEmail(email *string, user *User) error {
	return r.DB.Find(user, "email = ?", *email).Error
}