// Package main adalah titik masuk utama aplikasi.
package main

import (
	"github.com/gemini-cli/portfolio-chat-ai-go/internal/chat"
	"github.com/gemini-cli/portfolio-chat-ai-go/internal/user"
	"github.com/gemini-cli/portfolio-chat-ai-go/pkg/middleware"
	"github.com/labstack/echo/v4"
)

// Router adalah struct yang tidak memiliki field, hanya digunakan untuk men-grup-kan method SetupRoutes.
// Ini adalah pendekatan untuk menjaga agar logika routing tetap terorganisir.
type Router struct{}

// SetupRoutes mendefinisikan dan mengkonfigurasi semua rute (endpoints) aplikasi.
// Fungsi ini menerima semua handler dan konfigurasi yang dibutuhkan untuk mendaftarkan rute ke instance Echo.
func (h *Router) SetupRoutes(e *echo.Echo, userHandler *user.UserHandler, chatHandler *chat.ChatHandler, jwtSecret, basicUser, basicPass string) {
	// Endpoint publik untuk login, tidak memerlukan autentikasi.
	e.POST("/v1/login", userHandler.Login)

	// Grup rute yang diproteksi menggunakan Basic Auth.
	// Hanya request dengan header Basic Auth yang valid yang bisa mengakses rute di dalam grup ini.
	basicAuthGroup := e.Group("/v1", middleware.BasicAuthMiddleware(basicUser, basicPass))
	basicAuthGroup.POST("/users", userHandler.Create) // Endpoint untuk membuat user baru.

	// Grup rute yang diproteksi menggunakan JWT Auth.
	// Hanya request dengan header `Authorization: Bearer <token>` yang valid yang bisa mengakses rute di grup ini.
	jwtGroup := e.Group("/v1", middleware.JWTMiddleware(jwtSecret))
	jwtGroup.GET("/users/:id", userHandler.GetByID)  // Endpoint untuk mendapatkan data user.
	jwtGroup.GET("/ws", chatHandler.HandleWebSocket) // Endpoint untuk koneksi WebSocket chat.

}
