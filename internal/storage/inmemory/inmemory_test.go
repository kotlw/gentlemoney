package inmemory_test

import (
	"testing"

	"gentlemoney/internal/storage/inmemory"

	"github.com/stretchr/testify/suite"
)

type InmemoryStorageTestSuite struct {
	suite.Suite
}

func (s *InmemoryStorageTestSuite) TestStorageGet() {
	storage := inmemory.New()
	storage.Category()
	storage.Currency()
	storage.Account()
	storage.Transaction()
}

func TestInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(InmemoryStorageTestSuite))
}
