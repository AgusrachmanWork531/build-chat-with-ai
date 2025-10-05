// Package chat berisi semua logika yang terkait dengan domain chat real-time.
package chat

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

// Message adalah struct entitas utama untuk sebuah pesan chat.
type Message struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	status    string    `json:"status"` // "user" atau "system"
}

// ChatRepository mendefinisikan kontrak untuk lapisan persistensi chat.
// Dependensi: lapisan Usecase bergantung pada interface ini.
type ChatRepository interface {
	CreateMessage(ctx context.Context, msg *Message) error
	GetMessagesByRoom(ctx context.Context, roomID string) ([]*Message, error)
}

// ChatUsecase mendefinisikan kontrak untuk lapisan logika bisnis chat.
// Dependensi: lapisan Handler bergantung pada interface ini.
type ChatUsecase interface {
	// HandleStream adalah method utama yang menangani seluruh siklus hidup koneksi WebSocket.
	HandleStream(ctx context.Context, roomID string, c echo.Context) error
}
