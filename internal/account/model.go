package account

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID uint   `json:"user_id" gorm:"constraint:OnDelete:CASCADE"`
	Type   string `json:"type"`
	Title  string `json:"title"`
}
