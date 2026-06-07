package models

// MidtransNotification merepresentasikan format data yang dikirim oleh Payment Gateway
type MidtransNotification struct {
	OrderID           string `json:"order_id" binding:"required"`           // Ini adalah ID Transaksi kita
	TransactionStatus string `json:"transaction_status" binding:"required"` // "settlement", "pending", "expire", "cancel"
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"` // Untuk keamanan aslinya
}
