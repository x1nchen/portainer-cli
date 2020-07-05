package token

import (
	"github.com/x1nchen/portainer-cli/cache/internal"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "token"
	KeyToken   = "token"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// GetToken ...
func (service *Service) GetToken() (string, error) {
	var token string
	var err error
	key := internal.StringToBytes(KeyToken)

	token, err = internal.GetString(service.db, BucketName, key)
	if err != nil {
		return token, err
	}

	return token, nil
}

// SaveToken ...
func (service *Service) SaveToken(token string) error {
	key := internal.StringToBytes(KeyToken)
	return internal.UpdateString(service.db, BucketName, key, internal.StringToBytes(token))
}
