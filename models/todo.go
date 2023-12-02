package models

type TODO struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	UserID      uint   `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Finished    bool   `json:"finished"`
}
