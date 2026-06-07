package config

import (
	"log"
	"os"

	"github.com/midtrans/midtrans-go"
)

func InitMidtrans() {
	// Ambil dari .env
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	environment := os.Getenv("MIDTRANS_ENVIRONMENT")

	// Validasi
	if serverKey == "" || clientKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY atau MIDTRANS_CLIENT_KEY tidak ditemukan di .env")
	}

	// Setup Midtrans
	midtrans.ServerKey = serverKey
	midtrans.ClientKey = clientKey

	// Set environment (sandbox atau production)
	if environment == "production" {
		midtrans.Environment = midtrans.Production
	} else {
		midtrans.Environment = midtrans.Sandbox
	}

	log.Println("Midtrans initialized successfully!")
}
