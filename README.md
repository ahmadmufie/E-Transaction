# 🏗️ E-Transact Backend API

> Sebuah sistem backend tangguh untuk platform penyewaan alat berat B2B (Business-to-Business), dilengkapi dengan integrasi Payment Gateway otomatis.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin-Gonic](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Midtrans](https://img.shields.io/badge/Midtrans-000000?style=for-the-badge&logo=stripe&logoColor=white)

## 📖 Deskripsi Proyek
**E-Transact** adalah layanan RESTful API yang menjembatani transaksi penyewaan alat berat antar perusahaan (*Supplier* dan *Distributor*). Sistem ini dirancang menggunakan standar arsitektur industri dengan fokus pada keamanan data, manajemen stok, kontrak B2B, dan pencatatan pembayaran secara *real-time*.

## ✨ Fitur Utama
* **🔐 Role-Based Authentication:** Sistem *login* aman menggunakan JSON Web Token (JWT) dengan pembagian hak akses (*Role*).
* **🚜 Equipment & Inventory Management:** Operasi CRUD penuh untuk mengelola katalog armada alat berat beserta ketersediaan stok.
* **📄 B2B Contract Generation:** Pembuatan dan pencatatan draf kontrak sewa secara otomatis.
* **💳 Automated Payment Webhook:** Terintegrasi langsung dengan Midtrans (Snap API) untuk memproses pembayaran dan memperbarui status kontrak secara otomatis (*Real-time Webhook*).
* **🛡️ Full-Stack Security:** Dilengkapi dengan konfigurasi CORS yang ketat untuk keamanan lintas domain.

---

## 🛠️ Prasyarat (*Prerequisites*)
Pastikan Anda telah menginstal beberapa perangkat lunak berikut di mesin Anda:
* [Go](https://golang.org/doc/install) (versi 1.18 atau lebih baru)
* [PostgreSQL](https://www.postgresql.org/download/)
* Akun [Midtrans Sandbox](https://dashboard.sandbox.midtrans.com/) (untuk simulasi pembayaran)

---

## 🚀 Cara Menjalankan Aplikasi (*Getting Started*)

### 1. Kloning Repositori
```bash
git clone [https://github.com/username-anda/etransact-backend.git](https://github.com/username-anda/etransact-backend.git)
cd etransact-backend

### 2. Instalasi Dependensi
```bash
go mod tidy

### 3. Konfigurasi Environment
Buat file .env di folder root dengan merujuk pada file .env.example. Isi variabel dengan kredensial Anda:
```bash
# Database Config
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_anda
DB_NAME=etransact

# Midtrans Config
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxxxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxxxxxx
MIDTRANS_ENVIRONMENT=sandbox

### 4. Setup Database
Jalankan skrip SQL yang tersedia di folder database/schema.sql pada pgAdmin atau terminal psql Anda untuk membuat tabel dan menyuntikkan data dummy (seeding).

### 5. Jalankan Server Lokal
```bash
go run main.go
Server akan berjalan di http://localhost:8080.