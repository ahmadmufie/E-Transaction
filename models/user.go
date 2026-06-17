package models

import "gorm.io/gorm"

// User adalah model database untuk menyimpan data pengguna (Supplier & Distributor)
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`   // "-" agar password tidak ikut bocor saat kirim JSON response
	Role     string `gorm:"type:varchar(20);not null" json:"role"` // Nilainya: "supplier" atau "distributor"
}
