// Package database menyediakan fungsi helper untuk berinteraksi dengan database.
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectMongo menginisialisasi koneksi baru ke server MongoDB berdasarkan URI yang diberikan.
// Fungsi ini juga melakukan ping ke database untuk memastikan koneksi berhasil dibuat.
// Mengembalikan client MongoDB yang siap digunakan.
func ConnectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	// Menghubungkan ke server MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Melakukan ping ke node primary untuk memverifikasi bahwa koneksi sudah aktif.
	// Memberi batas waktu 2 detik untuk proses ping.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// DisconnectMongo menangani proses pemutusan koneksi dari server MongoDB secara graceful.
func DisconnectMongo(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
