package user

import "gorm.io/gorm"

type Repository interface {

}

type repository struct {
	DB *gorm.DB
}

func InitRepository(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}