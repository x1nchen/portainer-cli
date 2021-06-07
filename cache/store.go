package cache

import (
	"fmt"
	"path"
	"time"

	"github.com/x1nchen/portainer-cli/cache/container"
	"github.com/x1nchen/portainer-cli/cache/registry"

	"github.com/x1nchen/portainer-cli/cache/internal"

	"github.com/x1nchen/portainer-cli/cache/endpoint"
	"github.com/x1nchen/portainer-cli/cache/token"

	perr "github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

var (
	Token    = "token"
	KeyToken = []byte(Token)
)

// Store is the data cache for portainer
type Store struct {
	host             string // host is for bbolt bucket
	DBName           string
	bdb              *bolt.DB
	TokenService     *token.Service
	EndpointService  *endpoint.Service
	ContainerService *container.Service
	RegistryService  *registry.Service
}

// var _ Store = &bboltStore{}

func NewBoltStore(datadir string, host string) (*Store, error) {
	dbNamePrefix := internal.MD5Hash(host)
	dbName := fmt.Sprintf("%s.db", dbNamePrefix)
	db, err := bolt.Open(path.Join(datadir, dbName), 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, perr.WithMessage(err, "open local db")
	}

	s := &Store{
		host:   host,
		DBName: dbName,
		bdb:    db,
	}

	err = s.initServices()
	if err != nil {
		return nil, perr.WithMessage(err, "init store services")
	}
	return s, nil
}

func (b *Store) initServices() error {
	tokenService, err := token.NewService(b.bdb)
	if err != nil {
		return err
	}
	b.TokenService = tokenService

	endpointService, err := endpoint.NewService(b.bdb)
	if err != nil {
		return err
	}
	b.EndpointService = endpointService

	containerService, err := container.NewService(b.bdb)
	if err != nil {
		return err
	}
	b.ContainerService = containerService
	registryService, err := registry.NewService(b.bdb)
	if err != nil {
		return err
	}
	b.RegistryService = registryService
	return nil
}

// func (b *bboltStore) SaveToken(token string) error {
// 	return b.bdb.Update(func(tx *bolt.Tx) error {
// 		// 如果 bucket 不存在则，创建一个 bucket
// 		bucket, err := tx.CreateBucketIfNotExists(StringToBytesUnsafe(b.host))
// 		if err != nil {
// 			return perr.WithMessage(err, "create bucket")
// 		}
//
// 		// 将 key-value 写入到 bucket 中
// 		err = bucket.Put(KeyToken, []byte(token))
// 		if err != nil {
// 			return perr.WithMessage(err, "bucket put")
// 		}
// 		return nil
// 	})
// }
//
// // Endpoint gives access to the Endpoint data management layer
// func (b *bboltStore) Endpoint() *endpoint.Service {
// 	return b.EndpointService
// }
//
// func (b *bboltStore) GetToken() (string, error) {
// 	val, err := b.get(Token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return BytesToString(val), nil
// }
//
// func (b *bboltStore) get(key string) ([]byte, error) {
// 	var result []byte
// 	var err error
// 	err = b.bdb.View(func(tx *bolt.Tx) error {
// 		buc := tx.Bucket(StringToBytes(b.host))
//
// 		// 将 key-value 写入到 bucket 中
// 		result = buc.Get(StringToBytes(key))
// 		return nil
// 	})
//
// 	if err != nil {
// 		return []byte{}, err
// 	}
//
// 	return result, nil
// }
//
// 删除全部当前 host 下的数据
func (b *Store) RemoveAllData() error {
	return b.bdb.Update(func(tx *bolt.Tx) error {
		var ierr error
		if ierr = tx.DeleteBucket([]byte(b.host)); ierr != nil {
			return ierr
		}

		return nil
	})
}

func (b *Store) Close() error {
	return b.bdb.Close()
}
