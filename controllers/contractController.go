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

// 1. FUNGSI MEMBUAT KONTRAK (Sudah Ditambah Format Nomor Kontrak DK)
func CreateContractDraft(c *gin.Context) {
	var req models.ContractRequest

	// Validasi Input dari Client
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah atau tidak lengkap!"})
		return
	}

	// ANTI-TAMPERING: Ambil harga ASLI dari Database berdasarkan EquipmentID
	var basePrice float64
	err := config.DB.QueryRow("SELECT base_price FROM equipments WHERE id = $1", req.EquipmentID).Scan(&basePrice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alat berat tidak ditemukan di sistem."})
		return
	}

	// Kalkulasi Logika Bisnis (Harga * Durasi)
	totalValue := basePrice * float64(req.RentDuration)

	// GENERATE NOMOR KONTRAK SESUAI REQ DOSEN (Contoh: RESELLER01.DK01-timestamp)
	// Kita buat keunikan agar format ***.DK01 tetap terjaga
	contractNumber := fmt.Sprintf("%s.DK%s", req.ResellerID, time.Now().Format("02150405"))

	// DIGITAL SIGNATURE (Hashing dengan SHA-256)
	timestamp := time.Now().Format(time.RFC3339)
	rawDocumentString := fmt.Sprintf("%s|%s|%s|%d|%.2f|%s", contractNumber, req.ResellerID, req.EquipmentID, req.RentDuration, totalValue, timestamp)

	hashFunc := sha256.New()
	hashFunc.Write([]byte(rawDocumentString))
	hashedResult := hex.EncodeToString(hashFunc.Sum(nil))

	// Simpan ke Database (Pastikan tabel contracts Anda memiliki kolom contract_number)
	sqlStatement := `
		INSERT INTO contracts (contract_number, reseller_id, equipment_id, rent_duration_months, total_value, status, pdf_hash)
		VALUES ($1, $2, $3, $4, $5, 'DRAFT', $6)
		RETURNING id`

	var contractID string
	err = config.DB.QueryRow(
		sqlStatement,
		contractNumber,
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

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Draft Kontrak Berhasil Dibuat!",
		"contract_id":     contractID,
		"contract_number": contractNumber,
		"total_value":     totalValue,
		"digital_hash":    hashedResult,
	})
}

// 2. STRATEGI PENGAMBILAN DATA TERISOLASI (Mencegah Kebocoran Data antar Supplier/Distributor)
func GetMyContracts(c *gin.Context) {
	// Ambil userID dari token JWT yang lolos dari Middleware
	// Di jwtMiddleware Anda menggunakan c.Set("userID", claims["user_id"])
	loggedInUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid, silakan login kembali!"})
		return
	}

	// Jalankan Query dengan memfilter BERDASARKAN ID USER YANG LOGIN ($1)
	// Ini menjamin User A tidak akan bisa melihat data User B meskipun mereka menembak API yang sama
	rows, err := config.DB.Query(
		"SELECT id, contract_number, reseller_id, equipment_id, rent_duration_months, total_value, status, pdf_hash FROM contracts WHERE reseller_id = $1",
		loggedInUserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data database: " + err.Error()})
		return
	}
	defer rows.Close()

	// Struct lokal untuk mapping data rows
	type ContractResult struct {
		ID                 string  `json:"id"`
		ContractNumber     string  `json:"contract_number"`
		ResellerID         string  `json:"reseller_id"`
		EquipmentID        string  `json:"equipment_id"`
		RentDurationMonths int     `json:"rent_duration_months"`
		TotalValue         float64 `json:"total_value"`
		Status             string  `json:"status"`
		PdfHash            string  `json:"pdf_hash"`
	}

	var contracts []ContractResult
	for rows.Next() {
		var con ContractResult
		err := rows.Scan(&con.ID, &con.ContractNumber, &con.ResellerID, &con.EquipmentID, &con.RentDurationMonths, &con.TotalValue, &con.Status, &con.PdfHash)
		if err != nil {
			continue
		}
		contracts = append(contracts, con)
	}

	// Jika data tidak ditemukan atau kosong
	if len(contracts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Anda tidak memiliki kontrak yang terafiliasi.",
			"data":    []ContractResult{},
		})
		return
	}

	// Berhasil mengembalikan data yang khusus milik user tersebut
	c.JSON(http.StatusOK, gin.H{
		"message": "Data kontrak berhasil diambil (Terisolasi Aman)",
		"data":    contracts,
	})
}
