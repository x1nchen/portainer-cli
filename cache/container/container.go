package container

import (
	"bytes"
	"context"
	"fmt"

	"github.com/x1nchen/portainer-cli/cache/internal"
	"github.com/x1nchen/portainer-cli/err"
	climodel "github.com/x1nchen/portainer-cli/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketContainerName represents the name of the bucket where this service stores data.
	BucketContainerName = "container-name"
	// BucketContainerID only stores keys for container id
	BucketContainerID = "container-id"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketContainerName)
	if err != nil {
		return nil, err
	}

	err = internal.CreateBucket(db, BucketContainerID)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// GetContainByID find container by ID.
func (service *Service) GetContainByID(id string) (*climodel.ContainerExtend, error) {
	var mc climodel.ContainerExtend
	identifier := internal.StringToBytes(id)

	err := internal.GetObject(service.db, BucketContainerID, identifier, &mc)
	if err != nil {
		return nil, err
	}

	return &mc, nil
}

// UpdateContainer UpdateEndpoint updates an endpoint.
// TODO need to update data in BucketContainerName
func (service *Service) UpdateContainer(ID string, container *climodel.ContainerExtend) error {
	identifier := internal.StringToBytes(ID)

	return internal.UpdateObject(service.db, BucketContainerID, identifier, container)
}

// DeleteContainer deletes an endpoint.
func (service *Service) DeleteContainer(ID string) error {
	identifier := internal.StringToBytes(ID)
	err := service.db.Update(func(tx *bolt.Tx) error {
		bucketCI := tx.Bucket([]byte(BucketContainerID))
		value := bucketCI.Get(identifier)
		if value == nil {
			return err.ErrObjectNotFound
		}
		var container climodel.ContainerExtend

		if err := internal.UnmarshalObjectWithJsoniter(value, &container); err != nil {
			return err
		}

		// delete key in bucket-container-name
		if err := bucketCI.DeleteBucket(identifier); err != nil {
			return err
		}

		// delete key in bucket-container-id
		bucketCN := tx.Bucket([]byte(BucketContainerName))
		keyCN := container.KeyWithEndpoint()
		if err := bucketCN.Delete(internal.StringToBytes(keyCN)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// FindAllContainers return an array containing all containers
func (service *Service) FindAllContainers() ([]climodel.ContainerExtend, error) {
	var containers = make([]climodel.ContainerExtend, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketContainerName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var container climodel.ContainerExtend
			err := internal.UnmarshalObjectWithJsoniter(v, &container)
			if err != nil {
				return err
			}
			containers = append(containers, container)
		}

		return nil
	})

	return containers, err
}

// FuzzyFindContainerByName Endpoints return an array containing all the endpoints.
func (service *Service) FuzzyFindContainerByName(name string) ([]climodel.ContainerExtend, error) {
	var containers = make([]climodel.ContainerExtend, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketContainerName))
		cursor := bucket.Cursor()

		match := internal.StringToBytes(name)

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			// find by simple byte slice match
			if bytes.Contains(k, match) {
				var container climodel.ContainerExtend
				err := internal.UnmarshalObjectWithJsoniter(v, &container)
				if err != nil {
					return err
				}
				containers = append(containers, container)
			}
		}

		return nil
	})

	return containers, err
}

// SyncEndpointContainer sync containers upon endpoint
func (service *Service) SyncEndpointContainer(
	ctx context.Context,
	endpointID int,
	containers ...climodel.ContainerExtend,
) error {
	var containerIDList []string

	err := service.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketContainerName)).Cursor()
		prefix := internal.StringToBytes(fmt.Sprintf("%d:", endpointID))

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var cont climodel.ContainerExtend
			err := internal.UnmarshalObjectWithJsoniter(v, &cont)
			if err != nil {
				return err
			}
			containerIDList = append(containerIDList, cont.ID)
			if err = c.Delete(); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err = service.DeleteContainerInIDBucket(ctx, containerIDList...);err != nil {
		return err
	}

	if err = service.BatchUpdateContainers(ctx, containers...); err != nil {
		return err
	}

	return  nil
}

func (service *Service) DeleteContainerInIDBucket(ctx context.Context, idList ...string) error {
	return service.db.Batch(func(tx *bolt.Tx) error {
		bucketID := tx.Bucket([]byte(BucketContainerID))
		for _, id := range idList {
			err := bucketID.Delete(internal.StringToBytes(id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdateContainers 批量更新 container
// 1. bucket key 格式是容器的 name#container_id
// 2. bucket key 格式是容器的 id
func (service *Service) BatchUpdateContainers(ctx context.Context, containers ...climodel.ContainerExtend) error {
	return service.db.Batch(func(tx *bolt.Tx) error {
		bucketCN := tx.Bucket([]byte(BucketContainerName))
		for _, container := range containers {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}

			key := container.KeyWithEndpoint()
			err = bucketCN.Put(internal.StringToBytes(key), data)
			if err != nil {
				return err
			}
		}

		bucketCI := tx.Bucket([]byte(BucketContainerID))
		for _, container := range containers {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}

			key := container.KeyWithContainerID()
			err = bucketCI.Put(internal.StringToBytes(key), data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Synchronize creates, updates and deletes endpoints inside a single transaction.
func (service *Service) Synchronize(toCreate, toUpdate, toDelete []*climodel.ContainerExtend) error {
	err := service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketContainerName))
		for _, container := range toCreate {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}
			err = bucket.Put(internal.StringToBytes(container.KeyWithEndpoint()), data)
			if err != nil {
				return err
			}
		}
		for _, container := range toUpdate {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}
			err = bucket.Put(internal.StringToBytes(container.KeyWithEndpoint()), data)
			if err != nil {
				return err
			}
		}
		for _, container := range toDelete {
			err := bucket.Delete(internal.StringToBytes(container.KeyWithEndpoint()))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketContainerID))
		for _, container := range toCreate {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}
			err = bucket.Put(internal.StringToBytes(container.ID), data)
			if err != nil {
				return err
			}
		}
		for _, container := range toUpdate {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}
			err = bucket.Put(internal.StringToBytes(container.ID), data)
			if err != nil {
				return err
			}
		}
		for _, container := range toDelete {
			err := bucket.Delete(internal.StringToBytes(container.ID))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// TruncateDatabase delete bucket
func (service *Service) TruncateDatabase() error {
	return service.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket([]byte(BucketContainerName)); err != nil {
			return err
		}
		if err := tx.DeleteBucket([]byte(BucketContainerID)); err != nil {
			return err
		}
		return nil
	})
}

// CreateDatabase create bucket
func (service *Service) CreateDatabase() (
	bucketCN *bolt.Bucket,
	bucketCI *bolt.Bucket,
	err error,
) {

	err = service.db.Update(func(tx *bolt.Tx) error {
		if bucketCN, err = tx.CreateBucketIfNotExists([]byte(BucketContainerName)); err != nil {
			return err
		}

		if bucketCN, err = tx.CreateBucketIfNotExists([]byte(BucketContainerID)); err != nil {
			return err
		}
		return nil
	})

	return
}

// DB return db instance
func (service *Service) DB() *bolt.DB {
	return service.db
}
