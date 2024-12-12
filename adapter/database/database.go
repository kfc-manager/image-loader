package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-manager/image-loader/domain/image"
)

type Database interface {
	Close()
	InsertImage(img *image.Image) error
}

type database struct {
	conn *pgxpool.Pool
}

func (db *database) createTables() error {
	_, err := db.conn.Exec(
		context.Background(),
		`CREATE EXTENSION IF NOT EXISTS vector;`,
	)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS "image" (
		hash VARCHAR(64) PRIMARY KEY,
		size INTEGER NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		entropy DOUBLE PRECISION NOT NULL,
		format VARCHAR(10) NOT NULL,
		embedding VECTOR(768) DEFAULT NULL
	);`)
	return err
}

func New(host, port, name, user, pass string) (*database, error) {
	conn, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name),
	)
	if err != nil {
		return nil, err
	}
	db := &database{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *database) Close() {
	db.conn.Close()
}

func (db *database) InsertImage(img *image.Image) error {
	_, err := db.conn.Exec(context.Background(), `INSERT INTO "image" 
	(hash, size, width, height, entropy, format) 
	VALUES ($1, $2, $3, $4, $5, $6);`,
		img.Hash,
		img.Size,
		img.Width,
		img.Height,
		img.Entropy,
		img.Format,
	)
	return err
}
