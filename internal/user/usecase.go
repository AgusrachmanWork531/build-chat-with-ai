package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecaseImpl adalah implementasi dari UserUsecase yang berisi logika bisnis utama.
// Dependensi: bergantung pada UserRepository (untuk akses data) dan JWT Secret.
type UserUsecaseImpl struct {
	userRepo  UserRepository // Kontrak ke lapisan repository.
	jwtSecret []byte         // Kunci rahasia untuk membuat token.
}

// NewUserUsecase membuat instance baru dari UserUsecaseImpl.
func NewUserUsecase(userRepo UserRepository, jwtSecret string) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Create adalah logika bisnis untuk membuat pengguna baru.
// Termasuk hashing password dan memanggil repository untuk menyimpan data.
func (uc *UserUsecaseImpl) Create(ctx context.Context, email, password string) (*User, error) {
	// Melakukan hash pada password menggunakan bcrypt untuk keamanan.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa hash password: %w", err)
	}

	// Membuat entitas User baru.
	now := time.Now()
	newUser := &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	filter := Filter{
		Email: email,
	}

	// Memanggil lapisan repository untuk menyimpan user baru.
	if err := uc.userRepo.Create(ctx, filter, newUser); err != nil {
		return nil, fmt.Errorf("tidak bisa membuat user: %w", err)
	}

	return newUser, nil
}

// GetByID adalah logika bisnis untuk mendapatkan pengguna berdasarkan ID.
// Termasuk logika otorisasi untuk memeriksa apakah pengguna yang meminta berhak melihat data ini.
func (uc *UserUsecaseImpl) GetByID(ctx context.Context, actorID, userID string) (*User, error) {
	// Logika Otorisasi: hanya pengguna itu sendiri yang bisa melihat datanya.
	if actorID != userID {
		return nil, fmt.Errorf("akses ditolak: aktor %s tidak bisa mengakses data user %s", actorID, userID)
	}

	// Memanggil repository untuk mengambil data dari database.
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa mendapatkan user: %w", err)
	}

	return user, nil
}

// Login adalah logika bisnis untuk autentikasi pengguna.
// Memvalidasi kredensial dan membuat token JWT jika berhasil.
func (uc *UserUsecaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Cari pengguna berdasarkan email.
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("kredensial tidak valid: %w", err)
	}

	// 2. Bandingkan password yang diberikan dengan hash yang ada di database.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("kredensial tidak valid: %w", err)
	}

	// 3. Jika kredensial valid, buat token JWT.
	claims := &jwt.RegisteredClaims{
		Issuer:    "portfolio-chat-ai-go",
		Subject:   user.ID, // ID pengguna disimpan di dalam token.
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token berlaku 24 jam.
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("tidak bisa menandatangani token: %w", err)
	}

	return signedToken, nil
}
