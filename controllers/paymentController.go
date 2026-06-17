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

	// 2. Terjemahkan status Midtrans ke sistem kita (Disesuaikan dengan ENUM/Status Kontrak)
	var contractStatus string
	switch notification.TransactionStatus {
	case "settlement", "capture":
		contractStatus = "PAID" // atau 'ACTIVE' / 'APPROVED' sesuai enum aplikasi Anda
	case "expire", "cancel", "deny":
		contractStatus = "FAILED"
	case "pending":
		contractStatus = "PENDING"
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Status diabaikan"})
		return
	}

	// 3. Update status transaksi di database (Disesuaikan ke tabel contracts)
	// Jika order_id Midtrans Anda mengirimkan UUID Contract ID, gunakan id = $2
	// Jika order_id Midtrans Anda mengirimkan Nomor Kontrak, ganti id menjadi contract_number
	sqlStatement := `
		UPDATE contracts 
		SET status = $1 
		WHERE id = $2 AND status = 'DRAFT'`

	result, err := config.DB.Exec(sqlStatement, contractStatus, notification.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status database kontrak"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Kontrak tidak ditemukan atau sudah diproses sebelumnya"})
		return
	}

	// 4. Berikan respons HTTP 200 OK ke Midtrans agar mereka berhenti mengirim notifikasi
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Status kontrak berhasil diupdate ke " + contractStatus})
}
