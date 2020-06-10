package cache

import (
	perr "github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"path"
	"time"
)

var (
	Token    = "token"
	KeyToken = []byte(Token)
)

// Store is the data cache for portainer
type Store interface {
	// 保存用户 token
	SaveToken(token string) error
	// 获取用户 token
	GetToken() (string, error)
	RemoveAllData() error
	Close() error
}

var _ Store = &bboltStore{}

type bboltStore struct {
	host string // host is for bbolt bucket
	bdb  *bolt.DB
}

func NewBoltStore(datadir string, host string) (Store, error) {
	db, err := bolt.Open(path.Join(datadir, "portainer.db"), 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, perr.WithMessage(err, "open local db")
	}

	s := &bboltStore{
		host: host,
		bdb:  db,
	}
	return s, nil
}

func (b *bboltStore) SaveToken(token string) error {
	return b.bdb.Update(func(tx *bolt.Tx) error {
		// 如果 bucket 不存在则，创建一个 bucket
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.host))
		if err != nil {
			return perr.WithMessage(err, "create bucket")
		}

		// 将 key-value 写入到 bucket 中
		err = bucket.Put(KeyToken, []byte(token))
		if err != nil {
			return perr.WithMessage(err, "bucket put")
		}
		return nil
	})
}

func (b *bboltStore) GetToken() (string, error) {
	val, err := b.get(Token)
	if err != nil {
		return "", err
	}
	return BytesToString(val), nil
}

func (b *bboltStore) get(key string) ([]byte, error) {
	var result []byte
	var err error
	err = b.bdb.View(func(tx *bolt.Tx) error {
		bucket, ierr := tx.CreateBucketIfNotExists([]byte(b.host))
		if ierr != nil {
			return perr.WithMessage(err, "create bucket")
		}

		// 将 key-value 写入到 bucket 中
		result = bucket.Get(StringToBytes(key))
		return nil
	})

	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

// 删除全部当前 host 下的数据
func (b *bboltStore) RemoveAllData() error {
	return b.bdb.Update(func(tx *bolt.Tx) error {
		var ierr error
		if ierr = tx.DeleteBucket([]byte(b.host)); ierr != nil {
			return ierr
		}

		return nil
	})
}

func (b *bboltStore) Close() error {
	return b.bdb.Close()
}
