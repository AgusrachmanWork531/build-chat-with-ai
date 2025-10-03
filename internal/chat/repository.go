package chat

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryChatRepository adalah implementasi ChatRepository yang menyimpan pesan di dalam memori.
// PENTING: Implementasi ini hanya untuk tujuan demonstrasi dan data akan hilang saat aplikasi berhenti.
// Seharusnya diganti dengan implementasi yang menggunakan database, seperti MongoChatRepository.
type InMemoryChatRepository struct {
	mu       sync.RWMutex
	messages map[string][]*Message // Kunci adalah roomID, nilai adalah slice dari pesan.
}

// NewInMemoryChatRepository membuat instance baru dari InMemoryChatRepository.
func NewInMemoryChatRepository() *InMemoryChatRepository {
	return &InMemoryChatRepository{
		messages: make(map[string][]*Message),
	}
}

// CreateMessage menyimpan pesan baru ke dalam map di memori.
func (r *InMemoryChatRepository) CreateMessage(ctx context.Context, msg *Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[msg.RoomID] = append(r.messages[msg.RoomID], msg)
	return nil
}

// GetMessagesByRoom mengambil semua pesan untuk sebuah room dari memori.
func (r *InMemoryChatRepository) GetMessagesByRoom(ctx context.Context, roomID string) ([]*Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messages, ok := r.messages[roomID]
	if !ok {
		return nil, fmt.Errorf("tidak ada pesan untuk room %s", roomID)
	}
	return messages, nil
}