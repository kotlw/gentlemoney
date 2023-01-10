package service_test

import (
	"database/sql"
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/service"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	db                *sql.DB
	persistentStorage *sqlite.Category
	inmemoryStorage   *inmemory.Category
	service           *service.Category
	InitCategories    []*model.Category
}

func (s *CategoryServiceTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	s.persistentStorage, err = sqlite.NewCategory(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	s.inmemoryStorage = inmemory.NewCategory()

	s.service, err = service.NewCategory(s.persistentStorage, s.inmemoryStorage)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCategories = []*model.Category{
		{ID: 1, Title: "Health"},
		{ID: 2, Title: "Grocery"},
	}
}

func (s *CategoryServiceTestSuite) SetupTest() {
	for _, c := range s.InitCategories {
		_, err := s.persistentStorage.Insert(c)
		require.NoError(s.T(), err, "occurred in SetupTest")
	}

	err := s.service.Init()
	require.NoError(s.T(), err, "occurred in SetupTest")
}

func (s *CategoryServiceTestSuite) TestInsertPositive() {
	category := &model.Category{ID: 3, Title: "Sport"}
	expectedCategories := append(s.InitCategories, category)

	err := s.service.Insert(category)
	require.NoError(s.T(), err)

	persistentCategories, err := s.persistentStorage.GetAll()
	require.NoError(s.T(), err)
	inmemoryCategories := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistentCategories, expectedCategories)
	assert.ElementsMatch(s.T(), inmemoryCategories, expectedCategories)
}

func (s *CategoryServiceTestSuite) TestInsertNegative() {
	err := s.service.Insert(s.InitCategories[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Insert: e.db.Exec: UNIQUE constraint failed: category.title")
}

func (s *CategoryServiceTestSuite) TestUpdatePositive() {
	expectedCategories := make([]*model.Category, 2)
	copy(expectedCategories, s.InitCategories)
	expectedCategories[0].Title = "Taxi"

	err := s.service.Update(expectedCategories[0])
	require.NoError(s.T(), err)

	persistentCategories, err := s.persistentStorage.GetAll()
	inmemoryCategories := s.inmemoryStorage.GetAll()
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), expectedCategories, inmemoryCategories)
	assert.ElementsMatch(s.T(), persistentCategories, inmemoryCategories)
}

func (s *CategoryServiceTestSuite) TestUpdateNegative() {
	cc := s.service.GetAll()
	cc[0].Title = "Taxi"
	cc[0].ID = 10

	err := s.service.Update(cc[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Update: total affected rows 0 while expected 1")
	cc[0].ID = 1 // return real id to proper teardown
}

func (s *CategoryServiceTestSuite) TestDeletePositive() {
	cc := s.service.GetAll()
	expectedCategories := []*model.Category{cc[0]}

	err := s.service.Delete(cc[1])
	require.NoError(s.T(), err)

	persistentCategories, err := s.persistentStorage.GetAll()
	require.NoError(s.T(), err)
	inmemoryCategories := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistentCategories, expectedCategories)
	assert.ElementsMatch(s.T(), inmemoryCategories, expectedCategories)
}

func (s *CategoryServiceTestSuite) TestDeleteNegative() {
	cc := s.service.GetAll()
	cc[0].ID = 10

	err := s.service.Delete(cc[0])
	assert.ErrorContains(s.T(), err, "s.persistentStorage.Delete: total affected rows 0 while expected 1")
	cc[0].ID = 1 // return real id to proper teardown
}

func (s *CategoryServiceTestSuite) TestGetByID() {
	c := s.service.GetByID(2)
	assert.Equal(s.T(), s.InitCategories[1].Title, c.Title)
}

func (s *CategoryServiceTestSuite) TestGetByTitle() {
	c := s.service.GetByTitle(s.InitCategories[1].Title)
	assert.Equal(s.T(), s.InitCategories[1].ID, c.ID)
}

func (s *CategoryServiceTestSuite) TearDownTest() {
	for {
		cc := s.service.GetAll()
		if len(cc) == 0 {
			break
		}

		err := s.persistentStorage.Delete(cc[0].ID)
		require.NoError(s.T(), err, "occurred in TearDownTest")
		s.inmemoryStorage.Delete(cc[0])
	}
}

func (s *CategoryServiceTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestCategoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}
