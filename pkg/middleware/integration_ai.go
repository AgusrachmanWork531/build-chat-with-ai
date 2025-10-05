// Package middleware menyediakan fungsi-fungsi middleware kustom untuk framework Echo.
package middleware

import (
	"github.com/labstack/echo/v4"
)

// GeminiAPIKeyMiddleware memeriksa keberadaan header `X-Gemini-API-Key` pada request.
// Ini digunakan sebagai lapisan otorisasi sederhana untuk endpoint yang memerlukan akses ke Gemini API.
func GeminiAPIKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-Gemini-API-Key")
		if apiKey == "" {
			return echo.NewHTTPError(401, "Missing Gemini API Key")
		}
		// TODO: Tambahkan logika untuk memvalidasi format atau isi API key jika diperlukan.
		return next(c)
	}
}