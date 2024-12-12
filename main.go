package main

import (
	"fmt"
	"os"

	"github.com/kfc-manager/image-loader/adapter/bucket"
	"github.com/kfc-manager/image-loader/adapter/database"
	"github.com/kfc-manager/image-loader/adapter/queue"
	"github.com/kfc-manager/image-loader/service/load"
)

func main() {
	db, err := database.New(
		envOrPanic("DB_HOST"),
		envOrPanic("DB_PORT"),
		envOrPanic("DB_NAME"),
		envOrPanic("DB_USER"),
		envOrPanic("DB_PASS"),
	)
	if err != nil {
		panic(err)
	}

	b, err := bucket.New(envOrPanic("BUCKET_PATH"))
	if err != nil {
		panic(err)
	}

	q, err := queue.New(
		envOrPanic("QUEUE_HOST"),
		envOrPanic("QUEUE_PORT"),
		envOrPanic("QUEUE_NAME"),
	)
	if err != nil {
		panic(err)
	}

	if err := load.New(db, b, q).Load(); err != nil {
		panic(err)
	}
}

func envOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env variables: '%s'", key))
	}
	return val
}
