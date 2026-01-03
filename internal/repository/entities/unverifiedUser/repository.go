package unverifiedUser

import (
	"globe/internal/repository/entities/user"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	DeleteExpiredUsers(now time.Time)
	Create(user *UnverifiedUser) (*uint32, error)
	IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool
	FindByID(ID *uint32, user *UnverifiedUser) error
	DeleteByID(ID *uint32) error
	UpdateVerificationCode(ID *uint32, code *string) error
	FindCodeByID(ID *uint32, code *string) error
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

func (r *repository) Create(user *UnverifiedUser) (*uint32, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return &user.ID, nil
}

func (r *repository) IsUsernameOrEmailAlreadyInUse(username *string, email *string) bool {
	var count int64
	err := r.DB.Model(&user.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (r *repository) FindByID(ID *uint32, user *UnverifiedUser) error {
	return r.DB.Find(user, "id = ?", ID).Error
}

func (r *repository) DeleteByID(ID *uint32) error {
	return r.DB.Delete(&UnverifiedUser{}, " id = ?", ID).Error
}

func (r *repository) UpdateVerificationCode(ID *uint32, newCode *string) error {
	return r.DB.Model(&UnverifiedUser{}).Where("id = ?", ID).Update("code", &newCode).Error
}

func (r *repository) FindCodeByID(ID *uint32, code *string) error {
	return r.DB.Model(&UnverifiedUser{}).Where("id = ?", ID).Select("code").Scan(code).Error
}