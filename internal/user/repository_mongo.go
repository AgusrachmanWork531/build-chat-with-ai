package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserRepository adalah implementasi dari UserRepository yang menggunakan MongoDB sebagai penyimpanannya.
// Dependensi: bergantung pada koneksi database MongoDB (*mongo.Database).
type MongoUserRepository struct {
	db         *mongo.Database // Koneksi ke database spesifik di MongoDB.
	collection string          // Nama koleksi (tabel) yang digunakan, yaitu "users".
}

// NewMongoUserRepository membuat instance baru dari MongoUserRepository.
// Menerima koneksi database sebagai dependensi.
func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		db:         db,
		collection: "users",
	}
}

// Create menyimpan sebuah entitas User baru ke dalam koleksi `users` di MongoDB.
func (r *MongoUserRepository) Create(ctx context.Context, filter Filter, user *User) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"updated_at":    user.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        user.ID,
			"created_at": user.CreatedAt,
		},
	}

	result := r.db.Collection(r.collection).FindOneAndUpdate(ctx, filter, update, opts)
	return result.Err()
}

// GetByID mencari dan mengembalikan seorang pengguna berdasarkan ID-nya dari database.
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	// Mencari satu dokumen di koleksi `users` dimana field `_id` sama dengan id yang diberikan.
	err := r.db.Collection(r.collection).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

// GetByEmail mencari dan mengembalikan seorang pengguna berdasarkan email-nya dari database.
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	// Mencari satu dokumen di koleksi `users` dimana field `email` sama dengan email yang diberikan.
	err := r.db.Collection(r.collection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}
