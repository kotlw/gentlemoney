package storage_test

import (
	"database/sql"
	"testing"
	"time"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionStorageTestSuite struct {
	suite.Suite
	db               *sql.DB
	storage          *storage.Transaction
	InitTransactions []*model.Transaction
}

func (s *TransactionStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occured in SetupSuite")
	s.db = db

	s.storage, err = storage.NewTransaction(db)
	require.NoError(s.T(), err, "occured in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
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

func (s *TransactionStorageTestSuite) SetupTest() {
	stmt, err := s.db.Prepare(`INSERT INTO "transaction" (date, accountId, categoryId, amount, note) VALUES (?, ?, ?, ?, ?);`)
	require.NoError(s.T(), err, "occured in SetupTest")

	for _, tr := range s.InitTransactions {
		_, err := stmt.Exec(&tr.Date, &tr.Account.ID, &tr.Category.ID, &tr.Amount, &tr.Note)
		require.NoError(s.T(), err, "occured in SetupTest")
	}
}

func (s *TransactionStorageTestSuite) TestInsertPositive() {
	transaction := &model.Transaction{
		ID:       3,
		Date:     time.Date(2022, time.Month(2), 23, 1, 10, 30, 0, time.UTC),
		Account:  model.NewEmptyAccount(),
		Category: model.NewEmptyCategory(),
		Amount:   4321,
		Note:     "note3",
	}
	expectedTransactions := append(s.InitTransactions, transaction)

	_, err := s.storage.Insert(transaction)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedTransactions)
}

func (s *TransactionStorageTestSuite) TestUpdatePositive() {
	expectedTransactions := make([]*model.Transaction, len(s.InitTransactions))
	copy(expectedTransactions, s.InitTransactions)
	expectedTransactions[1].Amount = 10987
	expectedTransactions[1].Note = "changed"

	err := s.storage.Update(expectedTransactions[1])
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedTransactions)
}

func (s *TransactionStorageTestSuite) TestUpdateNegative() {
    tr := model.NewEmptyTransaction()
    tr.ID = 10
	err := s.storage.Update(tr)
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *TransactionStorageTestSuite) TestDeletePositive() {
	expectedTransactions := []*model.Transaction{s.InitTransactions[0]}

	err := s.storage.Delete(2)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedTransactions)
}

func (s *TransactionStorageTestSuite) TestDeleteNegative() {
	err := s.storage.Delete(10)
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *TransactionStorageTestSuite) TestGetAll() {
	allTransactions, err := s.storage.GetAll()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), s.InitTransactions, allTransactions)
}

func (s *TransactionStorageTestSuite) fetchActualData() []*model.Transaction {
	rows, err := s.db.Query(`SELECT id, date, amount, note, accountId, categoryId FROM "transaction";`)
	require.NoError(s.T(), err)
	defer func() {
		err = rows.Close()
		require.NoError(s.T(), err)
	}()

	res := make([]*model.Transaction, 0, 3)
	for rows.Next() {
		t := model.NewEmptyTransaction() 
		err = rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Note, &t.Account.ID, &t.Category.ID)
		require.NoError(s.T(), err)
		res = append(res, t)
	}

	return res
}

func (s *TransactionStorageTestSuite) TearDownTest() {
	stmt, err := s.db.Prepare(`DELETE FROM "transaction";`)
	require.NoError(s.T(), err, "occured in TearDownTest")

	_, err = stmt.Exec()
	require.NoError(s.T(), err, "occured in TearDownTest")
}

func (s *TransactionStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occured in TearDownSuite")
}

func TestTransactionStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionStorageTestSuite))
}
