package presenter_test

import (
	"database/sql"
	"testing"

	"github.com/kotlw/gentlemoney/internal/presenter"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PresenterTestSuite struct {
	suite.Suite
}

func (s *PresenterTestSuite) TestPresenterGet() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")

	persistentStorage, err := sqlite.New(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	inmemoryStorage := inmemory.New()

	service, err := service.New(persistentStorage, inmemoryStorage)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	presenter := presenter.New(service)

	presenter.Category()
	presenter.Currency()
	presenter.Account()
	presenter.Transaction()
}

func TestInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(PresenterTestSuite))
}
