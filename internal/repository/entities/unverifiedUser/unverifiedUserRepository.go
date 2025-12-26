package unverifiedUser

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *UnverifiedUser) (uint32, error)
	FindByID(id *uint32, user *UnverifiedUser) error
	DeleteByID(id *uint32) error
	DeleteExpired(now time.Time) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(user *UnverifiedUser) (uint32, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repository) FindByID(id *uint32, user *UnverifiedUser) error {
	return r.DB.Find(user, "id = ?", *id).Error
}

func (r *repository) DeleteByID(id *uint32) error {
	return r.DB.Delete(&UnverifiedUser{}, "id = ?", *id).Error
}

func (r *repository) DeleteExpired(now time.Time) error {
	return r.DB.Delete(&UnverifiedUser{}, "expired < ?", now).Error
}