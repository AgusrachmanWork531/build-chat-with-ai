// Package chat (lapisan handler) bertanggung jawab untuk menangani request terkait chat.
package chat

import (
	"github.com/labstack/echo/v4"
)

// ChatHandler adalah struct yang menangani request HTTP/WebSocket untuk domain Chat.
// Dependensi: bergantung pada ChatUsecase (kontrak lapisan bisnis).
type ChatHandler struct {
	chatUsecase ChatUsecase
}

// NewChatHandler membuat instance baru dari ChatHandler.
func NewChatHandler(chatUsecase ChatUsecase) *ChatHandler {
	return &ChatHandler{chatUsecase: chatUsecase}
}

// HandleWebSocket menangani request untuk upgrade ke koneksi WebSocket (GET /v1/ws/:roomID).
// Tugas utamanya adalah mengekstrak parameter dan meneruskan kontrol ke lapisan use case.
func (h *ChatHandler) HandleWebSocket(c echo.Context) error {
	// Mengambil ID room dari parameter URL.
	roomID := c.Param("roomID")
	// Memanggil use case untuk menangani seluruh logika streaming WebSocket.
	return h.chatUsecase.HandleStream(c.Request().Context(), roomID, c)
}