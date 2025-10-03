// Package bootstrap bertanggung jawab untuk menginisialisasi semua komponen inti aplikasi.
// Ini termasuk memuat konfigurasi, membuat koneksi database, dan menyatukan semuanya
// dalam satu container aplikasi yang siap digunakan.
package bootstrap

import (
	"context"
	"log"

	"github.com/gemini-cli/portfolio-chat-ai-go/pkg/config"
	"github.com/gemini-cli/portfolio-chat-ai-go/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Application adalah container utama yang menampung semua dependensi inti aplikasi,
// seperti konfigurasi (Env) dan koneksi database (Mongo).
// Jika ada koneksi lain seperti Redis, ia akan ditambahkan di sini.
type Application struct {
	Env   *config.Config // Menyimpan semua konfigurasi dari environment.
	Mongo *mongo.Client  // Client untuk koneksi ke MongoDB.
}

// NewApplication memuat konfigurasi, menginisialisasi koneksi, dan mengembalikan container aplikasi.
// Ini adalah titik awal untuk membangun seluruh dependensi aplikasi.
func NewApplication() *Application {
	app := &Application{}

	// 1. Memuat semua konfigurasi dari environment.
	app.Env = config.NewConfig()

	// 2. Menggunakan konfigurasi untuk membuat koneksi ke MongoDB.
	mongoClient, err := database.ConnectMongo(context.Background(), app.Env.MongoURI)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}
	app.Mongo = mongoClient
	log.Println("Berhasil terhubung ke MongoDB.")

	// Inisialisasi koneksi lain (misal: Redis) bisa ditambahkan di sini.

	return app
}

// Close secara graceful menutup semua koneksi yang ada di dalam container aplikasi.
// Fungsi ini dipanggil menggunakan `defer` di `main.go` untuk memastikan semua koneksi ditutup
// saat aplikasi berhenti.
func (app *Application) Close() {
	if err := database.DisconnectMongo(context.Background(), app.Mongo); err != nil {
		log.Fatalf("Gagal menutup koneksi MongoDB: %v", err)
	}
	log.Println("Berhasil menutup koneksi MongoDB.")
}
