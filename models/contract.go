package models

// ContractRequest adalah struktur data JSON yang dikirim oleh Klien/Reseller dari Frontend
type ContractRequest struct {
	ResellerID   string `json:"reseller_id" binding:"required"`
	EquipmentID  string `json:"equipment_id" binding:"required"`
	RentDuration int    `json:"rent_duration_months" binding:"required,min=1"`
}
