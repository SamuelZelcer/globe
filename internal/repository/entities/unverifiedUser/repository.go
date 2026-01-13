package unverifiedUser

import (
	"globe/internal/repository/entities"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	DeleteExpiredUsers(now time.Time)
	Create(user *entities.UnverifiedUser) (*uint64, error)
	IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool
	FindByID(ID *uint64, user *entities.UnverifiedUser) error
	DeleteByID(ID *uint64) error
	UpdateVerificationCode(ID *uint64, code *string) error
	FindCodeByID(ID *uint64, code *string) error
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

func (r *repository) Create(user *entities.UnverifiedUser) (*uint64, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return &user.ID, nil
}

func (r *repository) IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool {
	var count int64
	err := r.DB.Model(&entities.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (r *repository) FindByID(ID *uint64, user *entities.UnverifiedUser) error {
	return r.DB.Find(user, "id = ?", ID).Error
}

func (r *repository) DeleteByID(ID *uint64) error {
	return r.DB.Delete(&entities.UnverifiedUser{}, " id = ?", ID).Error
}

func (r *repository) UpdateVerificationCode(ID *uint64, newCode *string) error {
	return r.DB.Model(&entities.UnverifiedUser{}).Where("id = ?", ID).Update("code", &newCode).Error
}

func (r *repository) FindCodeByID(ID *uint64, code *string) error {
	return r.DB.Model(&entities.UnverifiedUser{}).Where("id = ?", ID).Select("code").Scan(code).Error
}