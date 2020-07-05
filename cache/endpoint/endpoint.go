package endpoint

import (
	"github.com/x1nchen/portainer-cli/cache/internal"
	"github.com/x1nchen/portainer/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "endpoint"
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

// Endpoint returns an endpoint by ID.
func (service *Service) Endpoint(ID int) (*model.Endpoint, error) {
	var endpoint model.Endpoint
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

// UpdateEndpoint updates an endpoint.
func (service *Service) UpdateEndpoint(ID int, endpoint *model.Endpoint) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, endpoint)
}

// DeleteEndpoint deletes an endpoint.
func (service *Service) DeleteEndpoint(ID int) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// Endpoints return an array containing all the endpoints.
func (service *Service) Endpoints() ([]model.Endpoint, error) {
	var endpoints = make([]model.Endpoint, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var endpoint model.Endpoint
			err := internal.UnmarshalObjectWithJsoniter(v, &endpoint)
			if err != nil {
				return err
			}
			endpoints = append(endpoints, endpoint)
		}

		return nil
	})

	return endpoints, err
}

// CreateEndpoint assign an ID to a new endpoint and saves it.
func (service *Service) CreateEndpoint(endpoint *model.Endpoint) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		// We manually manage sequences for endpoints
		err := bucket.SetSequence(uint64(endpoint.Id))
		if err != nil {
			return err
		}

		data, err := internal.MarshalObject(endpoint)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(endpoint.Id)), data)
	})
}

// GetNextIdentifier returns the next identifier for an endpoint.
func (service *Service) GetNextIdentifier() int {
	return internal.GetNextIdentifier(service.db, BucketName)
}

// Synchronize creates, updates and deletes endpoints inside a single transaction.
func (service *Service) Synchronize(toCreate, toUpdate, toDelete []*model.Endpoint) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		for _, endpoint := range toCreate {
			id, _ := bucket.NextSequence()
			endpoint.Id = int32(id)

			data, err := internal.MarshalObject(endpoint)
			if err != nil {
				return err
			}

			err = bucket.Put(internal.StringToBytes(string(endpoint.Id)), data)
			if err != nil {
				return err
			}
		}

		for _, endpoint := range toUpdate {
			data, err := internal.MarshalObject(endpoint)
			if err != nil {
				return err
			}

			err = bucket.Put(internal.StringToBytes(string(endpoint.Id)), data)
			if err != nil {
				return err
			}
		}

		for _, endpoint := range toDelete {
			err := bucket.Delete(internal.StringToBytes(string(endpoint.Id)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// BatchUpdateEndpoints ...
func (service *Service) BatchUpdateEndpoints(endpoints ...model.Endpoint) error {
	return service.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		for _, endpoint := range endpoints {
			data, err := internal.MarshalObject(endpoint)
			if err != nil {
				return err
			}

			err = bucket.Put(internal.StringToBytes(string(endpoint.Id)), data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
