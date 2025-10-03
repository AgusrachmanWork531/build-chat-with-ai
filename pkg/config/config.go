// Package config bertugas untuk membaca dan mengelola semua konfigurasi aplikasi.
// Konfigurasi diambil dari environment variables, dengan dukungan file .env untuk development.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua variabel konfigurasi aplikasi yang diambil dari environment.
type Config struct {
	AppPort       string // Port yang digunakan oleh server aplikasi.
	JWTSecret     string // Kunci rahasia untuk menandatangani token JWT.
	BasicAuthUser string // Username untuk Basic Authentication.
	BasicAuthPass string // Password untuk Basic Authentication.
	MongoURI      string // URI untuk koneksi ke database MongoDB.
	MongoDbName   string // Nama database MongoDB yang digunakan.
}

// NewConfig membuat instance Config baru dengan membaca environment variables.
// Fungsi ini akan menghentikan aplikasi jika variabel penting tidak ditemukan.
func NewConfig() *Config {
	// Memuat file .env jika ada. Perintah ini akan diabaikan jika file tidak ditemukan,
	// memungkinkan variabel di-set langsung di environment (umum untuk production).
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading variables from environment.")
	}

	return &Config{
		// Untuk AppPort, nilai default diberikan jika tidak ada di environment.
		AppPort: getEnvWithFallback("APP_PORT", "8080"),

		// Untuk variabel krusial, aplikasi akan berhenti jika tidak di-set.
		JWTSecret:     getEnvOrFatal("JWT_SECRET"),
		BasicAuthUser: getEnvOrFatal("BASIC_AUTH_USER"),
		BasicAuthPass: getEnvOrFatal("BASIC_AUTH_PASS"),
		MongoURI:      getEnvOrFatal("MONGO_URI"),
		MongoDbName:   getEnvOrFatal("MONGO_DB"),
	}
}

// getEnvWithFallback membaca environment variable berdasarkan key, atau mengembalikan nilai fallback jika tidak ada.
func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvOrFatal membaca environment variable berdasarkan key, atau menghentikan aplikasi jika tidak ada.
func getEnvOrFatal(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// Menghentikan eksekusi dan mencatat error fatal jika variabel yang dibutuhkan tidak ada.
	log.Fatalf("FATAL ERROR: Required environment variable not set: %s", key)
	return ""
}
