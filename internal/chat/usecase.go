package chat

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// upgrader adalah instance dari gorilla/websocket yang menangani proses upgrade koneksi HTTP ke WebSocket.
var (
	upgrader = websocket.Upgrader{
		// CheckOrigin mengizinkan koneksi dari origin manapun. Untuk production, ini harus dibatasi.
		// TODO: Ganti dengan daftar origin yang diizinkan di lingkungan produksi.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// ChatUsecaseImpl adalah implementasi dari ChatUsecase yang menangani logika real-time chat.
// Dependensi: bergantung pada ChatRepository untuk menyimpan pesan.
type ChatUsecaseImpl struct {
	chatRepo  ChatRepository
	chatMongo *MongoChatRepository
	mu        sync.RWMutex
	// rooms adalah map untuk menampung koneksi WebSocket yang aktif untuk setiap room.
	// Kunci luar adalah roomID, kunci dalam adalah pointer ke koneksi WebSocket.
	rooms map[string]map[*websocket.Conn]bool
}

// NewChatUsecase membuat instance baru dari ChatUsecaseImpl.
func NewChatUsecase(chatRepo ChatRepository, chatMongo *MongoChatRepository) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		chatRepo:  chatRepo,
		chatMongo: chatMongo,
		rooms:     make(map[string]map[*websocket.Conn]bool),
	}
}

// HandleStream adalah method utama yang menangani siklus hidup koneksi WebSocket.
func (uc *ChatUsecaseImpl) HandleStream(ctx context.Context, roomID string, c echo.Context) error {

	// 1. Upgrade koneksi HTTP ke koneksi WebSocket.
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// 2. Tambahkan koneksi baru ini ke dalam daftar koneksi aktif untuk room ini.
	uc.addConnection(roomID, ws)
	// Pastikan koneksi dihapus saat fungsi ini berakhir (koneksi terputus).
	defer uc.removeConnection(roomID, ws)

	// 3. Masuk ke loop tak terbatas untuk membaca pesan dari client.
	for {
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			// Jika ada error saat membaca (misal: client disconnect), hentikan loop.
			log.Println("read error:", err)
			break
		}

		// Mengambil token JWT dari context, yang seharusnya sudah divalidasi oleh middleware.
		userToken, ok := c.Get("user").(*jwt.Token)
		if !ok {
			// Jika token tidak ditemukan atau tipe-nya salah, ini adalah error internal.
			log.Println("failed to get user token from context")
			// Sebaiknya kirim pesan error ke client dan tutup koneksi.
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error"))
			break
		}
		// Mengekstrak claims dari token.
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("failed to cast claims from token")
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error"))
			break
		}

		// Mengambil userID dari claims.
		userID, ok := claims["sub"].(string)
		if !ok {
			log.Println("failed to get userID from claims")
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error"))
			break
		}

		// 4. Buat entitas Message baru.
		newMessage := &Message{
			ID:        uuid.NewString(),
			RoomID:    roomID,
			UserID:    userID,
			Content:   string(msgBytes),
			CreatedAt: time.Now(),
		}
		// 5. Simpan pesan ke database melalui repository.
		if err := uc.chatMongo.CreateMessage(ctx, newMessage); err != nil {
			log.Println("write error:", err)
			continue
		}

		// 6. Siarkan (broadcast) pesan ke semua client lain di room yang sama.
		uc.broadcast(roomID, newMessage)
	}

	return nil
}

// addConnection secara aman (thread-safe) menambahkan koneksi baru ke map `rooms`.
func (uc *ChatUsecaseImpl) addConnection(roomID string, ws *websocket.Conn) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.rooms[roomID]; !ok {
		uc.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	uc.rooms[roomID][ws] = true
}

// removeConnection secara aman (thread-safe) menghapus koneksi dari map `rooms`.
func (uc *ChatUsecaseImpl) removeConnection(roomID string, ws *websocket.Conn) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.rooms[roomID]; ok {
		delete(uc.rooms[roomID], ws)
		if len(uc.rooms[roomID]) == 0 {
			delete(uc.rooms, roomID)
		}
	}
}

// broadcast mengirimkan pesan ke semua koneksi yang aktif di sebuah room.
func (uc *ChatUsecaseImpl) broadcast(roomID string, msg *Message) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	for conn := range uc.rooms[roomID] {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("write error:", err)
			uc.removeConnection(roomID, conn)
		}
	}
}
