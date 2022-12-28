package inmemory_test

import (
	"testing"
	"time"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage/inmemory"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionInmemoryStorageTestSuite struct {
	suite.Suite
	storage          *inmemory.Transaction
	InitTransactions []*model.Transaction
}

func (s *TransactionInmemoryStorageTestSuite) SetupSuite() {
	s.storage = inmemory.NewTransaction()
	s.InitTransactions = []*model.Transaction{
		{
			ID:       1,
			Date:     time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			Account:  model.NewEmptyAccount(),
			Category: model.NewEmptyCategory(),
			Amount:   12345,
			Note:     "note1",
		},
		{
			ID:       2,
			Date:     time.Date(2022, time.Month(2), 22, 1, 10, 30, 0, time.UTC),
			Account:  model.NewEmptyAccount(),
			Category: model.NewEmptyCategory(),
			Amount:   67890,
			Note:     "note2",
		},
	}
}

func (s *TransactionInmemoryStorageTestSuite) SetupTest() {
	s.storage.Init(s.InitTransactions)
}

func (s *TransactionInmemoryStorageTestSuite) TestInsertPositive() {
	transaction := &model.Transaction{
		ID:       3,
		Date:     time.Date(2022, time.Month(2), 23, 1, 10, 30, 0, time.UTC),
		Account:  model.NewEmptyAccount(),
		Category: model.NewEmptyCategory(),
		Amount:   4321,
		Note:     "note3",
	}
	expectedTransactions := append(s.InitTransactions, transaction)

	s.storage.Insert(transaction)

	assert.ElementsMatch(s.T(), s.storage.GetAll(), expectedTransactions)
	assert.Equal(s.T(), transaction, s.storage.GetByID(transaction.ID))
}

func (s *TransactionInmemoryStorageTestSuite) TestDeletePositive() {
	s.storage.Delete(s.InitTransactions[1])

	assert.ElementsMatch(s.T(), s.storage.GetAll(), s.InitTransactions[:1])
}

func (s *TransactionInmemoryStorageTestSuite) TearDownTest() {
    for {
        tt := s.storage.GetAll()
        if len(tt) == 0 {
            break
        }
        s.storage.Delete(tt[0])
    }
}

func TestTransactionInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionInmemoryStorageTestSuite))
}
