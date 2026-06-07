package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"etransact-backend/config"
	"etransact-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateContractDraft(c *gin.Context) {
	var req models.ContractRequest

	// 1. Validasi Input dari Client
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah atau tidak lengkap!"})
		return
	}

	// 2. ANTI-TAMPERING: Ambil harga ASLI dari Database berdasarkan EquipmentID
	var basePrice float64
	err := config.DB.QueryRow("SELECT base_price FROM equipments WHERE id = $1", req.EquipmentID).Scan(&basePrice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alat berat tidak ditemukan di sistem."})
		return
	}

	// 3. Kalkulasi Logika Bisnis (Harga * Durasi)
	totalValue := basePrice * float64(req.RentDuration)

	// 4. DIGITAL SIGNATURE (Hashing dengan SHA-256)
	// Kita gabungkan data mentah menjadi satu string panjang
	timestamp := time.Now().Format(time.RFC3339)
	rawDocumentString := fmt.Sprintf("%s|%s|%d|%.2f|%s", req.ResellerID, req.EquipmentID, req.RentDuration, totalValue, timestamp)

	// Proses Hashing
	hashFunc := sha256.New()
	hashFunc.Write([]byte(rawDocumentString))
	hashedResult := hex.EncodeToString(hashFunc.Sum(nil))

	// 5. Simpan ke Database
	sqlStatement := `
		INSERT INTO contracts (reseller_id, equipment_id, rent_duration_months, total_value, status, pdf_hash)
		VALUES ($1, $2, $3, $4, 'DRAFT', $5)
		RETURNING id`

	var contractID string
	err = config.DB.QueryRow(
		sqlStatement,
		req.ResellerID,
		req.EquipmentID,
		req.RentDuration,
		totalValue,
		hashedResult,
	).Scan(&contractID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kontrak: " + err.Error()})
		return
	}

	// 6. Kembalikan Respons Sukses
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Draft Kontrak Berhasil Dibuat!",
		"contract_id":  contractID,
		"total_value":  totalValue,
		"digital_hash": hashedResult,
	})
}
