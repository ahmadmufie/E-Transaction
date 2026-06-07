package controllers

import (
	"net/http"

	"etransact-backend/config"
	"etransact-backend/models"

	"github.com/gin-gonic/gin"
)

// CreateEquipment - Menambahkan alat berat baru ke database
func CreateEquipment(c *gin.Context) {
	var input models.Equipment

	// 1. Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid atau kurang lengkap!"})
		return
	}

	// 2. Eksekusi Query SQL untuk Insert Data
	// Kita menggunakan $1, $2 (Parameterized Query) untuk mencegah SQL Injection
	sqlStatement := `
		INSERT INTO equipments (name, type, stock, base_price, retail_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id string
	err := config.DB.QueryRow(
		sqlStatement,
		input.Name,
		input.Type,
		input.Stock,
		input.BasePrice,
		input.RetailPrice,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database: " + err.Error()})
		return
	}

	// 3. Kembalikan respons sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Alat berat berhasil ditambahkan!",
		"id":      id,
	})
}

// GetEquipments - Mengambil semua data alat berat
func GetEquipments(c *gin.Context) {
	// 1. Eksekusi Query SQL untuk Select Data
	rows, err := config.DB.Query("SELECT id, name, type, stock, base_price, retail_price FROM equipments")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
		return
	}
	defer rows.Close()

	var equipments []models.Equipment

	// 2. Looping hasil query dan masukkan ke dalam array (slice)
	for rows.Next() {
		var eq models.Equipment
		if err := rows.Scan(&eq.ID, &eq.Name, &eq.Type, &eq.Stock, &eq.BasePrice, &eq.RetailPrice); err != nil {
			continue
		}
		equipments = append(equipments, eq)
	}

	// 3. Kembalikan respons JSON
	c.JSON(http.StatusOK, gin.H{
		"data": equipments,
	})
}

// UpdateEquipment - Mengubah data alat berat (PUT /api/equipments/:id)
func UpdateEquipment(c *gin.Context) {
	id := c.Param("id")
	var eq models.Equipment
	if err := c.ShouldBindJSON(&eq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	_, err := config.DB.Exec(
		"UPDATE equipments SET name=$1, type=$2, base_price=$3, retail_price=$4, stock=$5 WHERE id=$6",
		eq.Name, eq.Type, eq.BasePrice, eq.RetailPrice, eq.Stock, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Armada berhasil diupdate"})
}

// DeleteEquipment - Menghapus data alat berat (DELETE /api/equipments/:id)
func DeleteEquipment(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM equipments WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Armada berhasil dihapus"})
}
