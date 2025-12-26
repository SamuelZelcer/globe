package user

type User struct {
	ID uint32 `gorm:"type:bigint;primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(60);unique;not null"`
	Email string `gorm:"type:varchar(120);unique;not null"`
	Password []byte `gorm:"type:varchar(1000);not null"`
}