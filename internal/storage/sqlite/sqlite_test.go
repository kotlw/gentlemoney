package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SqliteStorageTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *SqliteStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	_, err = db.Exec(`CREATE TABLE t(id INTEGER);`)
	require.NoError(s.T(), err, "occurred in SetupSuite")
}

func (s *SqliteStorageTestSuite) TestNewCategoryNegative() {
	_, err := s.db.Exec(`CREATE UNIQUE INDEX category ON t (id);`)
	require.NoError(s.T(), err)

	_, err = sqlite.New(s.db)
	assert.ErrorContains(s.T(), err, "NewCategory: s.CreateTableIfNotExists: there is already an index named category")

	_, err = s.db.Exec(`DROP INDEX category;`)
	require.NoError(s.T(), err)
}

func (s *SqliteStorageTestSuite) TestNewCurrencyNegative() {
	_, err := s.db.Exec(`CREATE UNIQUE INDEX currency ON t (id);`)
	require.NoError(s.T(), err)

	_, err = sqlite.New(s.db)
	assert.ErrorContains(s.T(), err, "NewCurrency: s.CreateTableIfNotExists: there is already an index named currency")

	_, err = s.db.Exec(`DROP INDEX currency;`)
	require.NoError(s.T(), err)
}

func (s *SqliteStorageTestSuite) TestNewAccountNegative() {
	_, err := s.db.Exec(`CREATE UNIQUE INDEX account ON t (id);`)
	require.NoError(s.T(), err)

	_, err = sqlite.New(s.db)
	assert.ErrorContains(s.T(), err, "NewAccount: s.CreateTableIfNotExists: there is already an index named account")

	_, err = s.db.Exec(`DROP INDEX account;`)
	require.NoError(s.T(), err)
}

func (s *SqliteStorageTestSuite) TestNewTransactionNegative() {
	storage, _ := sqlite.New(s.db)
	storage.Category()
	storage.Currency()
	storage.Account()
	storage.Transaction()
}

func (s *SqliteStorageTestSuite) TestStorageGet() {
	_, err := s.db.Exec(`CREATE UNIQUE INDEX "transaction" ON t (id);`)
	require.NoError(s.T(), err)

	_, err = sqlite.New(s.db)
	assert.ErrorContains(s.T(), err, "NewTransaction: s.CreateTableIfNotExists: there is already an index named transaction")

	_, err = s.db.Exec(`DROP INDEX "transaction";`)
	require.NoError(s.T(), err)
}

func (s *SqliteStorageTestSuite) TearDownTest() {
	_, err := s.db.Exec(`DROP TABLE IF EXISTS category;
                         DROP TABLE IF EXISTS currency;
                         DROP TABLE IF EXISTS account;
                         DROP TABLE IF EXISTS "transaction";`)
	require.NoError(s.T(), err, "occurred in TearDownTest")
}

func (s *SqliteStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestSqliteStorageTestSuite(t *testing.T) {
	suite.Run(t, new(SqliteStorageTestSuite))
}
