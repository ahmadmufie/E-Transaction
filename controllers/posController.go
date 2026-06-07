package controllers

import (
	"fmt"
	"net/http"

	"etransact-backend/config"
	"etransact-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CheckoutPOS(c *gin.Context) {
	var req models.POSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang tidak valid!"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi DB"})
		return
	}

	var totalAmount float64
	var transactionID string

	// Penyesuaian ENUM Database
	dbPayMethod := req.PayMethod
	if dbPayMethod == "TRANSFER" {
		dbPayMethod = "DEBIT"
	}

	err = tx.QueryRow(
		"INSERT INTO pos_transactions (cashier_id, total_amount, pay_method, payment_status) VALUES ($1, $2, $3, 'PENDING') RETURNING id",
		req.CashierID, 0, dbPayMethod,
	).Scan(&transactionID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat header transaksi"})
		return
	}

	for _, item := range req.Items {
		var retailPrice float64
		var currentStock int

		err := tx.QueryRow("SELECT retail_price, stock FROM equipments WHERE id = $1", item.EquipmentID).Scan(&retailPrice, &currentStock)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
			return
		}

		if currentStock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi"})
			return
		}

		subTotal := retailPrice * float64(item.Quantity)
		totalAmount += subTotal

		_, err = tx.Exec(
			"INSERT INTO pos_transaction_details (transaction_id, equipment_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, item.EquipmentID, item.Quantity, subTotal,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail transaksi"})
			return
		}

		_, err = tx.Exec("UPDATE equipments SET stock = stock - $1 WHERE id = $2", item.Quantity, item.EquipmentID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate stok"})
			return
		}
	}

	_, err = tx.Exec("UPDATE pos_transactions SET total_amount = $1 WHERE id = $2", totalAmount, transactionID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengkalkulasi total"})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi permanen"})
		return
	}

	responseData := gin.H{
		"message":        "Transaksi POS Berhasil!",
		"transaction_id": transactionID,
		"total_amount":   totalAmount,
		"pay_method":     req.PayMethod,
	}

	// Integrasi Midtrans
	if req.PayMethod == "TRANSFER" || req.PayMethod == "EWALLET" {
		snapResp, errSnap := snap.CreateTransaction(&snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transactionID,
				GrossAmt: int64(totalAmount),
			},
		})

		if errSnap != nil {
			fmt.Println("❌ ERROR MIDTRANS:", errSnap.GetMessage())
			responseData["midtrans_error"] = errSnap.GetMessage()
		} else {
			responseData["snap_token"] = snapResp.Token
			responseData["redirect_url"] = snapResp.RedirectURL
		}
	}

	c.JSON(http.StatusOK, responseData)
}

// ... GetTransactionHistory tetap sama seperti sebelumnya
func GetTransactionHistory(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, cashier_id, total_amount, pay_method, payment_status, created_at 
		FROM pos_transactions 
		ORDER BY created_at DESC`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat"})
		return
	}
	defer rows.Close()

	var history []gin.H
	for rows.Next() {
		var id, cashierID, payMethod, status, createdAt string
		var total float64
		err := rows.Scan(&id, &cashierID, &total, &payMethod, &status, &createdAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data riwayat transaksi"})
			return
		}

		history = append(history, gin.H{
			"transaction_id": id,
			"cashier":        cashierID,
			"total":          total,
			"method":         payMethod,
			"status":         status,
			"date":           createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
