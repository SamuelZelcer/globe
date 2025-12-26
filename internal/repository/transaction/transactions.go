package transaction

import "gorm.io/gorm"


type Transactions interface {
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) error
}

type transactions struct {
	database *gorm.DB
}

func InitTransactions(database *gorm.DB) Transactions {
	return &transactions{database: database}
}

func (t *transactions) BeginTransaction() *gorm.DB {
	return t.database.Begin()
}

func (t *transactions) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}