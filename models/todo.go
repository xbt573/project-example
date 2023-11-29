package models

type TODO struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Finished    bool   `json:"finished,omitempty"`
}
