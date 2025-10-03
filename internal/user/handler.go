// Package user (lapisan handler) bertanggung jawab untuk menangani request HTTP yang masuk.
package user

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// UserHandler adalah struct yang menangani request HTTP untuk domain User.
// Dependensi: bergantung pada UserUsecase (kontrak lapisan bisnis).
type UserHandler struct {
	userUsecase UserUsecase
}

// NewUserHandler membuat instance baru dari UserHandler.
func NewUserHandler(userUsecase UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// Create menangani request untuk membuat pengguna baru (POST /v1/users).
// Endpoint ini diproteksi oleh Basic Auth.
func (h *UserHandler) Create(c echo.Context) error {
	// Struct untuk menampung data dari request body.
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Mengikat (bind) data JSON dari body request ke dalam struct `req`.
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "body request tidak valid"})
	}

	// Memanggil lapisan use case untuk menjalankan logika bisnis pembuatan user.
	user, err := h.userUsecase.Create(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Mengembalikan data user yang baru dibuat dengan status 201 Created.
	return c.JSON(http.StatusCreated, user)
}

// GetByID menangani request untuk mendapatkan data pengguna berdasarkan ID (GET /v1/users/:id).
// Endpoint ini diproteksi oleh JWT Auth.
func (h *UserHandler) GetByID(c echo.Context) error {
	// Mengambil data pengguna dari token JWT yang sudah divalidasi oleh middleware.
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*jwt.RegisteredClaims)
	actorID := claims.Subject // ID pengguna yang melakukan request (aktor).

	// Mengambil ID pengguna yang ingin dilihat datanya dari parameter URL.
	userID := c.Param("id")

	// Memanggil lapisan use case, yang akan berisi logika otorisasi.
	user, err := h.userUsecase.GetByID(c.Request().Context(), actorID, userID)
	if err != nil {
		// Jika use case mengembalikan error (misal: akses ditolak), kirim status 403 Forbidden.
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// Login menangani request autentikasi pengguna (POST /v1/login).
// Endpoint ini publik.
func (h *UserHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "body request tidak valid"})
	}

	// Memanggil lapisan use case untuk memvalidasi kredensial dan mendapatkan token.
	token, err := h.userUsecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// Jika kredensial tidak valid, kirim status 401 Unauthorized.
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Mengembalikan token JWT ke client.
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}