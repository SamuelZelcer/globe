package unverifiedUser

import (
	"globe/internal/repository/entities"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	DeleteExpiredUsers(now time.Time)
	Create(user *entities.UnverifiedUser) (uint64, error)
	FindByID(ID uint64, user *entities.UnverifiedUser) error
	DeleteByID(ID uint64) error
	UpdateVerificationCode(ID uint64, code *string) error
	FindCodeByID(ID uint64, code *string) error
}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) DeleteExpiredUsers(now time.Time) {
	r.DB.Delete(&entities.UnverifiedUser{}, "expired < ?", now)
}

func (r *repository) Create(user *entities.UnverifiedUser) (uint64, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repository) FindByID(ID uint64, user *entities.UnverifiedUser) error {
	return r.DB.First(user, ID).Error
}

func (r *repository) DeleteByID(ID uint64) error {
	return r.DB.Delete(&entities.UnverifiedUser{}, ID).Error
}

func (r *repository) UpdateVerificationCode(ID uint64, newCode *string) error {
	return r.DB.Model(&entities.UnverifiedUser{}).Where(ID).Update("code", newCode).Error
}

func (r *repository) FindCodeByID(ID uint64, code *string) error {
	return r.DB.Model(&entities.UnverifiedUser{}).Where(ID).Select("code").Scan(code).Error
}