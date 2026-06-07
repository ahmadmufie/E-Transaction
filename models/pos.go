package models

// POSItem merepresentasikan satu baris barang yang di-scan oleh kasir
type POSItem struct {
	EquipmentID string `json:"equipment_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

// POSRequest adalah keranjang belanja keseluruhan yang dikirim dari Frontend POS
type POSRequest struct {
	CashierID string    `json:"cashier_id" binding:"required"`
	PayMethod string    `json:"pay_method" binding:"required"` // TUNAI, DEBIT,EWALLET
	Items     []POSItem `json:"items" binding:"required,min=1"`
}
