package cache

import (
	"fmt"
	"testing"
	"time"

	perr "github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/x1nchen/portainer-cli/model"
	bolt "go.etcd.io/bbolt"
)

type BoltStoreTestSuite struct {
	suite.Suite
	bs  *Store
	bdb *bolt.DB
}

func TestBoltStore(t *testing.T) {
	suite.Run(t, new(BoltStoreTestSuite))
}

func (s *BoltStoreTestSuite) SetupSuite() {
	bs, bdb, err := newTestBoltStore()
	s.NoError(err)
	s.bs = bs
	s.bdb = bdb
}

func (s *BoltStoreTestSuite) SetupTest() {
	// 在每次测试前执行
}

func (s *BoltStoreTestSuite) TearDownSuite() {
	err := s.bs.Close()
	s.NoError(err)
}

func (s *BoltStoreTestSuite) TearDownTest() {
}

func newTestBoltStore() (*Store, *bolt.DB, error) {
	db, err := bolt.Open("test.db", 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, nil, perr.WithMessage(err, "open local db")
	}

	s := &Store{
		host: "test",
		bdb:  db,
	}
	err = s.initServices()
	if err != nil {
		return s, db, perr.WithMessage(err, "init store services")
	}
	return s, db, nil
}

func (s *BoltStoreTestSuite) TestBoltDBSeek() {
	err := s.bdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b, err := tx.CreateBucketIfNotExists([]byte("test"))

		if err != nil {
			return err
		}

		err = b.Put([]byte("zzztest1"), []byte("test1"))

		if err != nil {
			return err
		}
		err = b.Put([]byte("xtes"), []byte("test2"))
		if err != nil {
			return err
		}

		return nil
	})

	s.NoError(err)

	err = s.bdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte("test")).Cursor()

		prefix := []byte("test")
		// for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		for k, v := c.Seek(prefix); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})

	s.NoError(err)

}

func (s *BoltStoreTestSuite) TestSaveToken() {
	err := s.bs.TokenService.SaveToken("testjwttoken")
	s.Require().NoError(err)
}

func (s *BoltStoreTestSuite) TestGetToken() {
	err := s.bs.TokenService.SaveToken("testjwttoken")
	s.NoError(err)

	t, err := s.bs.TokenService.GetToken()
	s.Require().NoError(err)
	s.Equal("testjwttoken", t)
}

func (s *BoltStoreTestSuite) TestRegistyGetUser() {
	user := &model.RegistryUser{
		Username:      "testuser",
		Password:      "testpwd",
		Email:         "testemail",
		ServerAddress: "testaddr",
	}
	err := s.bs.RegistryService.UpdateUser(user)
	s.Require().NoError(err)

	userGet, err := s.bs.RegistryService.GetUser()
	s.Require().NoError(err)
	s.Equal(user, userGet)
}
