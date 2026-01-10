package entities

type Product struct {
    ID          uint32 `gorm:"primaryKey;autoIncrement"`
    Name        string `gorm:"type:varchar(100);not null"`
    Price       uint32 `gorm:"not null"`
    Description string `gorm:"type:varchar(800);not null"`
    Owner       uint32 
    User        User `gorm:"foreignKey:Owner;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}