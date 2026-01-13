package entities

import "time"

type RefreshToken struct {
	ID uint64 `gorm:"type:bigint;primaryKey;not null"`
	Token string `gorm:"type:varchar(1000);not null"`
	Expired time.Time `gorm:"index:idx_refresh_token_expiration;type:timestamp;not null"`
}