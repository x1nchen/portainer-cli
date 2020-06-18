package cache

import (
	"testing"
	"time"

	perr "github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	bolt "go.etcd.io/bbolt"
)

type BoltStoreTestSuite struct {
	suite.Suite
	bs *bboltStore
}

func TestBoltStore(t *testing.T) {
	suite.Run(t, new(BoltStoreTestSuite))
}

func (s *BoltStoreTestSuite) SetupSuite() {
	bs, err := newTestBoltStore()
	s.NoError(err)
	s.bs = bs
}

func (s *BoltStoreTestSuite) SetupTest() {
	// 在每次测试前执行
}

func (s *BoltStoreTestSuite) TearDownSuite() {
	err := s.bs.Close()
	s.NoError(err)
}

func (s *BoltStoreTestSuite) TearDownTest() {
	// 在每次测试后执行移除数据
	err := s.bs.RemoveAllData()
	s.NoError(err)
}

func newTestBoltStore() (*bboltStore, error) {
	db, err := bolt.Open("test.db", 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, perr.WithMessage(err, "open local db")
	}

	s := &bboltStore{
		host: "test",
		bdb:  db,
	}
	return s, nil
}

func (s *BoltStoreTestSuite) TestSaveToken() {
	err := s.bs.SaveToken("testjwttoken")
	s.NoError(err)
}

func (s *BoltStoreTestSuite) TestGetToken() {
	err := s.bs.SaveToken("testjwttoken")
	s.NoError(err)

	t, err := s.bs.GetToken()
	s.NoError(err)
	s.Equal("testjwttoken", t)
}
