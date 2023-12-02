package models

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Login        string `json:"login" validate:"required,min=1"`
	Password     string `json:"password" validate:"required,min=1" gorm:"-:all"` // last tag ignores read, write and migrations
	PasswordHash string
}
