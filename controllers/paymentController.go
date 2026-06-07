package controllers

import (
	"net/http"

	"etransact-backend/config"
	"etransact-backend/models"

	"github.com/gin-gonic/gin"
)

// HandleMidtransWebhook memproses notifikasi dari Payment Gateway
func HandleMidtransWebhook(c *gin.Context) {
	var notification models.MidtransNotification

	// 1. Terima dan validasi JSON dari Midtrans
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format notifikasi tidak valid"})
		return
	}

	// KEAMANAN: Di standar Enterprise, di sini kita WAJIB
	// memvalidasi notification.SignatureKey menggunakan SHA512 dan Server Key Midtrans kita
	// untuk memastikan request ini benar-benar dari Midtrans, bukan hacker.
	// (Untuk simulasi ini, kita anggap valid).

	// 2. Terjemahkan status Midtrans ke sistem kita
	var paymentStatus string
	switch notification.TransactionStatus {
	case "settlement", "capture":
		paymentStatus = "PAID"
	case "expire", "cancel", "deny":
		paymentStatus = "FAILED"
	case "pending":
		paymentStatus = "PENDING"
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Status diabaikan"})
		return
	}

	// 3. Update status transaksi di database
	sqlStatement := `
		UPDATE pos_transactions 
		SET payment_status = $1 
		WHERE id = $2 AND payment_status = 'PENDING'`

	result, err := config.DB.Exec(sqlStatement, paymentStatus, notification.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status database"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Transaksi tidak ditemukan atau sudah diproses sebelumnya"})
		return
	}

	// 4. Berikan respons HTTP 200 OK ke Midtrans agar mereka berhenti mengirim notifikasi
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Status pembayaran berhasil diupdate ke " + paymentStatus})
}
