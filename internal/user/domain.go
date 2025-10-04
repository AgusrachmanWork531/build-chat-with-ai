// Package user berisi semua logika yang terkait dengan domain pengguna (user).
package user

import (
	"context"
	"time"
)

// User adalah struct entitas utama untuk domain pengguna.
// Struct ini merepresentasikan data seorang pengguna di dalam sistem.
// Tag `json` digunakan untuk serialisasi/deserialisasi JSON saat berkomunikasi via API.
// Tag `bson` digunakan oleh driver MongoDB untuk memetakan struct ke dokumen BSON.
type User struct {
	ID           string    `json:"id" bson:"_id"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"password_hash"` // `json:"-"` berarti field ini tidak akan pernah dikirim dalam response JSON.
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type Filter struct {
	Email string
}

// UserRepository mendefinisikan kontrak (interface) untuk lapisan persistensi (database).
// Setiap struct repository (misal: MongoUserRepository) harus mengimplementasikan semua method ini.
// Ini memungkinkan kita untuk menukar implementasi database tanpa mengubah logika bisnis.
// Dependensi: lapisan Usecase bergantung pada interface ini.
type UserRepository interface {
	Create(ctx context.Context, filter Filter, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// UserUsecase mendefinisikan kontrak (interface) untuk lapisan logika bisnis (use case).
// Setiap struct use case (misal: UserUsecaseImpl) harus mengimplementasikan semua method ini.
// Dependensi: lapisan Handler bergantung pada interface ini.
type UserUsecase interface {
	Create(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, actorID, userID string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error) // Mengembalikan token JWT.
}
