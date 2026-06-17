package models

import "gorm.io/gorm"

// Contract adalah model database GORM untuk menyimpan data kontrak yang terisolasi
type Contract struct {
	gorm.Model
	ContractNumber string `gorm:"type:varchar(50);unique;not null" json:"contract_number"` // Contoh format: SUPP_TechAsus.DK01
	Title          string `gorm:"type:varchar(255);not null" json:"title"`
	Content        string `gorm:"type:longtext" json:"content"`
	Status         string `gorm:"type:varchar(20);default:'draft'" json:"status"` // "draft", "active", "terminated"

	// STRATEGI ISOLASI DATA (Multi-Tenancy)
	OwnerID   uint `gorm:"not null" json:"owner_id"`   // ID Pembuat Kontrak (User ID)
	PartnerID uint `gorm:"not null" json:"partner_id"` // ID Rekan Bisnis (User ID)
}

// ContractRequest adalah struktur data JSON yang dikirim dari Frontend (Kodingan Asli Anda)
type ContractRequest struct {
	ResellerID   string `json:"reseller_id" binding:"required"`
	EquipmentID  string `json:"equipment_id" binding:"required"`
	RentDuration int    `json:"rent_duration_months" binding:"required,min=1"`
}
