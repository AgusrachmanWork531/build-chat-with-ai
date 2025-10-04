package chat

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoChatRepository struct {
	// Tambahkan field yang diperlukan untuk koneksi MongoDB, misalnya:
	db         *mongo.Database
	collection string // Nama koleksi untuk menyimpan pesan, misalnya "messages"
}

// NewMongoChatRepository membuat instance baru dari MongoChatRepository.
func NewMongoChatRepository(db *mongo.Database) *MongoChatRepository {
	return &MongoChatRepository{
		db:         db,
		collection: "messages",
	}
}

// CreateMessage menyimpan pesan baru ke dalam koleksi `messages` di MongoDB.
func (r *MongoChatRepository) CreateMessage(ctx context.Context, msg *Message) error {
	// Implementasi penyimpanan pesan ke MongoDB.
	// Contoh (pseudo-code):
	fmt.Printf("Simulasi menyimpan pesan ke MongoDB: %+v\n", msg)
	collection := r.db.Collection(r.collection)
	_, err := collection.InsertOne(ctx, msg)
	return err
}

func (r *MongoChatRepository) GetMessagesByRoom(ctx context.Context, roomID string) ([]*Message, error) {
	// Implementasi pengambilan pesan berdasarkan roomID dari MongoDB.
	// Contoh (pseudo-code):
	/*
		collection := r.db.Collection("messages")
		cursor, err := collection.Find(ctx, bson.M{"room_id": roomID})
		if err != nil {
			return nil, err
		}
		var messages []*Message
		if err := cursor.All(ctx, &messages); err != nil {
			return nil, err
		}
		return messages, nil
	*/
	fmt.Println("Simulasi mengambil pesan dari MongoDB untuk roomID:", roomID)
	return []*Message{}, nil
}
