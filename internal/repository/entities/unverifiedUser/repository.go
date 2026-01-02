package unverifiedUser

import (
	"globe/internal/repository/entities/user"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	DeleteExpiredUsers(now time.Time)
	Create(user *UnverifiedUser) (uint32, error)
	IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) DeleteExpiredUsers(now time.Time) {
	r.DB.Delete(&UnverifiedUser{}, "expired < ?", now)
}

func (r *repository) Create(user *UnverifiedUser) (uint32, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repository) IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool {
	var count int64
	err := r.DB.Model(&user.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}