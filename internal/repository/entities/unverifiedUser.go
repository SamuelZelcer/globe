package entities

import "time"

type UnverifiedUser struct {
	ID uint64 `gorm:"type:bigint;primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(60);not null"`
	Email string `gorm:"type:varchar(80);not null"`
	Password []byte `gorm:"type:varchar(1000);not null"`
	Expired time.Time `gorm:"index:idx_unverified_user_expiration;type:timestamp;not null"`
	Code string `gorm:"type:varchar(6)"`
}