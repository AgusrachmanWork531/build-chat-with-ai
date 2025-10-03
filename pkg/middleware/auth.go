// Package middleware menyediakan fungsi-fungsi middleware kustom untuk framework Echo.
package middleware

import (
	"crypto/subtle"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTMiddleware membuat dan mengembalikan sebuah middleware Echo untuk validasi token JWT.
// Middleware ini menggunakan secret yang diberikan untuk memverifikasi tanda tangan token.
// Dependensi: github.com/labstack/echo-jwt/v4
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey: []byte(secret),
	}
	return echojwt.WithConfig(config)
}

// BasicAuthMiddleware membuat dan mengembalikan sebuah middleware Echo untuk validasi Basic Auth.
// Middleware ini membandingkan username dan password dari request dengan yang seharusnya
// menggunakan `subtle.ConstantTimeCompare` untuk mencegah timing attacks.
// Dependensi: github.com/labstack/echo/v4/middleware
func BasicAuthMiddleware(username, password string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(user, pass string, c echo.Context) (bool, error) {
		// Menggunakan ConstantTimeCompare untuk keamanan.
		if subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1 {
			return true, nil
		}
		return false, nil
	})
}