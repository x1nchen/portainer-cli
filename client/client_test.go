// +build integration

package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PortainerClientTestSuite struct {
	suite.Suite
	token    string
	c        *PortainerClient
	host     string
	user     string
	password string
}

func TestPortainerClientTestSuite(t *testing.T) {
	suite.Run(t, new(PortainerClientTestSuite))
}

func (s *PortainerClientTestSuite) SetupSuite() {
	s.host = ""
	s.user = ""
	s.password = ""
	s.c = NewPortainerClient(s.host, "")
	token, err := s.c.Auth(context.Background(), s.user, s.password)
	s.NoError(err)
	s.c.CarryToken(token)
}

func (s *PortainerClientTestSuite) SetupTest() {
	// 在每次测试前执行
	// Create httpexpect instance

}

func (s *PortainerClientTestSuite) TearDownSuite() {}

func (s *PortainerClientTestSuite) TearDownTest() {
	// 在每次测试后执行
}

func (s *PortainerClientTestSuite) TestListEndpoint() {
	res, err := s.c.ListEndpoint(context.Background())
	s.NoError(err)
	for _, r := range res {
		s.T().Log(r)
	}
}

func (s *PortainerClientTestSuite) TestListContainer() {
	res, err := s.c.ListContainer(context.Background(), 78)
	s.NoError(err)
	for _, r := range res {
		s.T().Logf("%+v", r)
	}
}
