package main

import (
	"log"

	"github.com/gemini-cli/portfolio-chat-ai-go/internal/chat"
	"github.com/gemini-cli/portfolio-chat-ai-go/internal/user"
	"github.com/gemini-cli/portfolio-chat-ai-go/pkg/bootstrap"
	"github.com/labstack/echo/v4"
)

// main adalah fungsi utama dan titik masuk dari seluruh aplikasi.
func main() {
	// 1. Mem-bootstrap aplikasi: memuat config, konek ke DB, dll.
	app := bootstrap.NewApplication()
	// Menjadwalkan penutupan koneksi database saat fungsi main selesai.
	defer app.Close()

	// 2. Mengambil konfigurasi dan koneksi dari container `app`.
	cfg := app.Env
	db := app.Mongo.Database(cfg.MongoDbName) // TODO: Nama DB seharusnya dari config.

	// 3. Inisialisasi semua lapisan (dependency injection).
	// Menggunakan implementasi repository MongoDB yang sesungguhnya.
	userRepo := user.NewMongoUserRepository(db)
	userUsecase := user.NewUserUsecase(userRepo, cfg.JWTSecret)
	userHandler := user.NewUserHandler(userUsecase)

	// Untuk chat, kita masih menggunakan repository in-memory.
	chatRepo := chat.NewInMemoryChatRepository()
	chatUsecase := chat.NewChatUsecase(chatRepo)
	chatHandler := chat.NewChatHandler(chatUsecase)

	// 4. Membuat instance baru dari web server Echo.
	e := echo.New()

	// 5. Mendaftarkan semua rute (endpoints) ke server Echo.
	router := &Router{}
	router.SetupRoutes(e, userHandler, chatHandler, cfg.JWTSecret, cfg.BasicAuthUser, cfg.BasicAuthPass)

	// 6. Menjalankan server.
	log.Printf("Server berjalan di port %s", cfg.AppPort)
	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
