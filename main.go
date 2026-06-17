package main

import (
	"etransact-backend/config"
	"etransact-backend/controllers"
	"etransact-backend/middlewares"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // 1. Tambahkan import ini
)

func main() {
	fmt.Println("Memulai Server E-Transact...")

	// 2. Load file .env di awal aplikasi sebelum memanggil config lain
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal membaca file .env! Pastikan file .env sudah dibuat di root folder.")
	}

	// 3. Koneksi ke Database
	config.ConnectDB()

	// 4. Inisialisasi Midtrans
	config.InitMidtrans()

	// 5. Inisialisasi Router Gin
	r := gin.Default()

	// Middleware CORS (Untuk Keamanan Full-Stack)
	r.Use(CORSMiddleware())

	// PUBLIC ROUTE
	r.POST("/api/login", controllers.Login)
	r.POST("/api/webhooks/midtrans", controllers.HandleMidtransWebhook)
	r.PUT("/api/equipments/:id", controllers.UpdateEquipment)
	r.DELETE("/api/equipments/:id", controllers.DeleteEquipment)

	// PRIVATE ROUTE
	protected := r.Group("/api")
	protected.Use(middlewares.RequireAuth())
	{
		protected.GET("/profile", func(c *gin.Context) {
			email, _ := c.Get("userEmail")
			role, _ := c.Get("userRole")

			c.JSON(http.StatusOK, gin.H{
				"message": "Selamat datang di area rahasia!",
				"email":   email,
				"role":    role,
			})
		})

		protected.POST("/equipments", controllers.CreateEquipment)
		protected.GET("/equipments", controllers.GetEquipments)
		protected.POST("/contracts", controllers.CreateContractDraft)
		protected.POST("/pos/checkout", controllers.CheckoutPOS)
		protected.GET("/pos/history", controllers.GetTransactionHistory)
	}

	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}

// CORSMiddleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
