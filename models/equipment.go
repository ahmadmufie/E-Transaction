package models

// Equipment merepresentasikan struktur data dari tabel equipments di database
type Equipment struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"` // 'HEAVY_MACHINERY' atau 'LIGHT_PARTS'
	Stock       int     `json:"stock" binding:"required,min=0"`
	BasePrice   float64 `json:"base_price" binding:"required,min=0"`
	RetailPrice float64 `json:"retail_price" binding:"required,min=0"`
}
