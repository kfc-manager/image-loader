package bucket

import (
	"errors"
	"os"
)

type Bucket interface {
	Get(key string) ([]byte, error)
	List() ([]string, error)
}

type bucket struct {
	path string
}

func New(path string) (*bucket, error) {
	if len(path) < 1 {
		return nil, errors.New("path can't be empty")
	}
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	return &bucket{path: path}, nil
}

func (b *bucket) Get(key string) ([]byte, error) {
	body, err := os.ReadFile(b.path + "/" + key)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (b *bucket) List() ([]string, error) {
	if err := os.MkdirAll(b.path, 0777); err != nil {
		return nil, err
	}

	// Open the directory
	dir, err := os.Open(b.path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read all the directory entries
	files, err := dir.Readdir(-1) // -1 to read all entries
	if err != nil {
		return nil, err
	}

	result := []string{}
	// Iterate through the entries
	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, nil
}
