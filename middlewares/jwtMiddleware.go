package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key (Harus sama persis dengan yang ada di authController)
var jwtSecret = []byte("super_rahasia_fintech_2026")

// Fungsi pelindung jalur API
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil token dari Header HTTP (Authorization: Bearer <token>)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Token tidak ditemukan."})
			return
		}

		// 2. Pisahkan kata "Bearer" dan ambil tokennya saja
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma yang digunakan sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode enkripsi tidak valid")
			}
			return jwtSecret, nil
		})

		// 4. Cek apakah ada error saat parse atau token tidak valid
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa!"})
			return
		}

		// 5. (Opsional) Ambil data di dalam token dan simpan di context untuk digunakan controller
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userEmail", claims["email"])
			c.Set("userRole", claims["role"])
		}

		// 6. Loloskan request ke Controller selanjutnya
		c.Next()
	}
}
