package controllers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"etransact-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan Password wajib diisi!"})
		return
	}

	// 1. CARI USER DI DATABASE
	var id, fullName, role, passwordHash string
	// Menggunakan Raw SQL PostgreSQL ($1) - Pastikan config.DB adalah *sql.DB
	err := config.DB.QueryRow("SELECT id, full_name, role, password_hash FROM users WHERE email = $1", input.Email).
		Scan(&id, &fullName, &role, &passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak terdaftar!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error database"})
		}
		return
	}

	// 2. VERIFIKASI PASSWORD (BCRYPT)
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
		return
	}

	// 3. GENERATE TOKEN DENGAN CLAIM NYATA (RBAC)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   input.Email,
		"role":    role, // SUPERADMIN/CASHIER/RESELLER
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// [SECURITY PATCH] Ambil secret dari .env, jangan di-hardcode!
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "super_rahasia_fintech_2026" // Fallback sementara untuk lokal
	}

	tokenString, _ := token.SignedString([]byte(secretKey))

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
		"user": gin.H{
			"id":   id,
			"name": fullName,
			"role": role,
		},
	})
}
