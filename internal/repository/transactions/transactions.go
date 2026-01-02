package transactions

import "gorm.io/gorm"

type Transactions interface {
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) error
}

type transactions struct {
	DB *gorm.DB
}

func Init(DB *gorm.DB) Transactions {
	return &transactions{DB: DB}
}

func (t *transactions) BeginTransaction() *gorm.DB {
	return t.DB.Begin()
}

func (t *transactions) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}