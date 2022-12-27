package sqlite_test

import (
	"database/sql"
	"testing"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountSqliteStorageTestSuite struct {
	suite.Suite
	db           *sql.DB
	storage      *sqlite.Account
	InitAccounts []*model.Account
}

func (s *AccountSqliteStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occured in SetupSuite")
	s.db = db

	s.storage, err = sqlite.NewAccount(db)
	require.NoError(s.T(), err, "occured in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitAccounts = []*model.Account{
		{ID: 1, Name: "Card1", Currency: model.NewEmptyCurrency()},
		{ID: 2, Name: "Card2", Currency: model.NewEmptyCurrency()},
	}
}

func (s *AccountSqliteStorageTestSuite) SetupTest() {
	stmt, err := s.db.Prepare(`INSERT INTO account(name, currencyId) VALUES (?, ?);`)
	require.NoError(s.T(), err, "occured in SetupTest")

	for _, account := range s.InitAccounts {
		_, err := stmt.Exec(&account.Name, &account.Currency.ID)
		require.NoError(s.T(), err, "occured in SetupTest")
	}
}

func (s *AccountSqliteStorageTestSuite) TestInsertPositive() {
	account := &model.Account{ID: 3, Name: "Card3", Currency: &model.Currency{ID: 2}}
	expectedAccounts := append(s.InitAccounts, account)

	_, err := s.storage.Insert(account)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedAccounts)
}

func (s *AccountSqliteStorageTestSuite) TestInsertNegative() {
	_, err := s.storage.Insert(s.InitAccounts[1])
    assert.EqualError(s.T(), err, "e.db.Exec: UNIQUE constraint failed: account.name")
}

func (s *AccountSqliteStorageTestSuite) TestUpdatePositive() {
	expectedAccounts := make([]*model.Account, len(s.InitAccounts))
	copy(expectedAccounts, s.InitAccounts)
	expectedAccounts[1].Name = "Card10"

	err := s.storage.Update(expectedAccounts[1])
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedAccounts)
}

func (s *AccountSqliteStorageTestSuite) TestUpdate() {
	err := s.storage.Update(&model.Account{ID: 10, Currency: model.NewEmptyCurrency()})
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *AccountSqliteStorageTestSuite) TestDeletePositive() {
	expectedAccounts := []*model.Account{s.InitAccounts[0]}

	err := s.storage.Delete(2)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedAccounts)
}

func (s *AccountSqliteStorageTestSuite) TestDeleteNegative() {
	err := s.storage.Delete(10)
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *AccountSqliteStorageTestSuite) TestGetAll() {
	allAccounts, err := s.storage.GetAll()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), allAccounts, s.InitAccounts)
}

func (s *AccountSqliteStorageTestSuite) fetchActualData() []*model.Account {
	rows, err := s.db.Query(`SELECT id, name, currencyId FROM account;`)
	require.NoError(s.T(), err)
	defer func() {
		err = rows.Close()
		require.NoError(s.T(), err)
	}()

	res := make([]*model.Account, 0, 3)
	for rows.Next() {
		t := model.NewEmptyAccount()
		err = rows.Scan(&t.ID, &t.Name, &t.Currency.ID)
		require.NoError(s.T(), err)
		res = append(res, t)
	}

	return res
}

func (s *AccountSqliteStorageTestSuite) TearDownTest() {
	stmt, err := s.db.Prepare(`DELETE FROM account;`)
	require.NoError(s.T(), err, "occured in TearDownTest")

	_, err = stmt.Exec()
	require.NoError(s.T(), err, "occured in TearDownTest")
}

func (s *AccountSqliteStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occured in TearDownSuite")
}

func TestAccountSqliteStorageTestSuite(t *testing.T) {
	suite.Run(t, new(AccountSqliteStorageTestSuite))
}
