package category

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	UserID *uint  `json:"user_id" gorm:"constraint:OnDelete:CASCADE"`
	Title  string `json:"title"`
}
