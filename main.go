package main

import (
	"etransact-backend/config"
	"etransact-backend/controllers"
	"etransact-backend/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Memulai Server E-Transact...")

	// 1. Koneksi ke Database
	config.ConnectDB()

	// 2. Inisialisasi Midtrans
	config.InitMidtrans()

	// 3. Inisialisasi Router Gin
	r := gin.Default()

	// 3. Middleware CORS (Untuk Keamanan Full-Stack)
	r.Use(CORSMiddleware())

	// 4. PUBLIC ROUTE
	// Siapa saja bisa akses tanpa token (untuk Login)
	r.POST("/api/login", controllers.Login)

	// Endpoint Webhook Payment Gateway (Harus Publik)
	r.POST("/api/webhooks/midtrans", controllers.HandleMidtransWebhook)
	// Pastikan letaknya berdekatan dengan rute GET dan POST equipments Anda
	r.PUT("/api/equipments/:id", controllers.UpdateEquipment)
	r.DELETE("/api/equipments/:id", controllers.DeleteEquipment)

	// 5. PRIVATE ROUTE
	// Semua rute di dalam grup ini akan melewati Satpam (Middleware)
	protected := r.Group("/api")
	protected.Use(middlewares.RequireAuth())
	{
		// Endpoint Tes Profil
		protected.GET("/profile", func(c *gin.Context) {
			email, _ := c.Get("userEmail")
			role, _ := c.Get("userRole")

			c.JSON(http.StatusOK, gin.H{
				"message": "Selamat datang di area rahasia!",
				"email":   email,
				"role":    role,
			})
		})

		// Rute untuk Modul Alat Berat (Equipments)
		protected.POST("/equipments", controllers.CreateEquipment) // Tambah Data
		protected.GET("/equipments", controllers.GetEquipments)    // Lihat Data

		// Rute untuk kontrak B2B
		protected.POST("/contracts", controllers.CreateContractDraft) // Buat Draft Kontrak

		// Rute baru
		protected.POST("/pos/checkout", controllers.CheckoutPOS)         // Proses Checkout POS
		protected.GET("/pos/history", controllers.GetTransactionHistory) // Riwayat Transaksi POS
	}

	// 6. Jalankan Server
	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}

// CORSMiddleware - Mengatur izin akses lintas port (Full-Stack Security)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Mengizinkan Next.js
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Jika browser hanya mengecek izin (OPTIONS), langsung jawab dengan 204 (No Content) atau 200 OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
