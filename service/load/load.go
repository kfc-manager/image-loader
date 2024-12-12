package load

import (
	"errors"

	"github.com/kfc-manager/image-loader/adapter/bucket"
	"github.com/kfc-manager/image-loader/adapter/database"
	"github.com/kfc-manager/image-loader/adapter/queue"
	"github.com/kfc-manager/image-loader/domain/image"
)

type Service interface {
	Load() error
}

type service struct {
	db     database.Database
	bucket bucket.Bucket
	queue  queue.Queue
}

func New(db database.Database, b bucket.Bucket, q queue.Queue) *service {
	return &service{db: db, bucket: b, queue: q}
}

func (s *service) Load() error {
	list, err := s.bucket.List()
	if err != nil {
		return err
	}

	for _, key := range list {
		b, err := s.bucket.Get(key)
		if err != nil {
			return err
		}

		img, err := image.LoadImage(b)
		if err != nil {
			return err
		}
		if img.Hash != key {
			return errors.New("hash mismatch")
		}

		err = s.db.InsertImage(img)
		if err != nil {
			return err
		}

		err = s.queue.Push(img.Hash)
		if err != nil {
			return err
		}
	}

	return nil
}
