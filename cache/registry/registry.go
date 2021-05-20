package registry

import (
	"github.com/x1nchen/portainer-cli/cache/internal"
	"github.com/x1nchen/portainer-cli/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "registry"
	Identifier = "user"
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

// GetUser get registry user
func (service *Service) GetUser() (*model.RegistryUser, error) {
	var user model.RegistryUser

	err := internal.GetObject(
		service.db,
		BucketName,
		internal.StringToBytes(Identifier),
		&user,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser update registry user
func (service *Service) UpdateUser(user *model.RegistryUser) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		data, err := internal.MarshalObject(user)
		if err != nil {
			return err
		}

		return bucket.Put(internal.StringToBytes(Identifier), data)
	})
}


// TruncateDatabase delete bucket
func (service *Service) TruncateDatabase() error {
	return service.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(BucketName))
	})
}

// CreateDatabase create bucket
func (service *Service) CreateDatabase() (*bolt.Bucket, error) {
	var buc *bolt.Bucket
	var err error
	err = service.db.Update(func(tx *bolt.Tx) error {
		buc, err = tx.CreateBucketIfNotExists([]byte(BucketName))
		return err
	})

	return buc, err
}