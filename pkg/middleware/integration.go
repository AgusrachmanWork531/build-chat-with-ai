// Package middleware menyediakan fungsi-fungsi middleware untuk aplikasi.
package middleware

import (
	"github.com/gemini-cli/portfolio-chat-ai-go/pkg/config"
	"github.com/labstack/echo/v4"
)

// Middleware adalah struct yang berfungsi sebagai container untuk semua middleware yang digunakan dalam aplikasi.
// Ini memusatkan logika pembuatan middleware dan membuatnya mudah untuk di-inject ke router.
type Middleware struct {
	BasicAuth    echo.MiddlewareFunc
	JWT          echo.MiddlewareFunc
	GeminiAPIKey echo.MiddlewareFunc
}

// NewMiddleware membuat instance baru dari struct Middleware.
// Ia mengambil semua dependensi yang diperlukan (seperti config) untuk menginisialisasi
// semua middleware yang dibutuhkan oleh aplikasi.
func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		BasicAuth:    BasicAuthMiddleware(cfg.BasicAuthUser, cfg.BasicAuthPass),
		JWT:          JWTMiddleware(cfg.JWTSecret),
		GeminiAPIKey: GeminiAPIKeyMiddleware,
	}
}
