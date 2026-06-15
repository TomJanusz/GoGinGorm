package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name   string `json:"name" form:"name"`
	Statut string `json:"status" form:"statut"`
}
