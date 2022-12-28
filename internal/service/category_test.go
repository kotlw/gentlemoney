package service_test

import (
	"database/sql"
	"sort"
	"testing"

	"gentlemoney/internal/model"
	"gentlemoney/internal/service"
	"gentlemoney/internal/storage/inmemory"
	"gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type categoryList []*model.Category

func (cc categoryList) Len() int           { return len(cc) }
func (cc categoryList) Less(i, j int) bool { return cc[i].Title < cc[j].Title }
func (cc categoryList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

type CategoryServiceStorageTestSuite struct {
	suite.Suite
	db                *sql.DB
	persistantStorage *sqlite.Category
	inmemoryStorage   *inmemory.Category
	service           *service.Category
	InitCategories    []*model.Category
}

func (s *CategoryServiceStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occured in SetupSuite")
    s.db = db

	s.persistantStorage, err = sqlite.NewCategory(db)
	require.NoError(s.T(), err, "occured in SetupSuite")

	s.inmemoryStorage = inmemory.NewCategory()

	s.service, err = service.NewCategory(s.persistantStorage, s.inmemoryStorage)
	require.NoError(s.T(), err, "occured in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCategories = []*model.Category{
		{ID: 1, Title: "Health"},
		{ID: 2, Title: "Grocery"},
	}
}

func (s *CategoryServiceStorageTestSuite) SetupTest() {
	for _, c := range s.InitCategories {
		_, err := s.persistantStorage.Insert(c)
		require.NoError(s.T(), err, "occured in SetupTest")
	}

    err := s.service.Init()
    require.NoError(s.T(), err, "occured in SetupTest")
}

func (s *CategoryServiceStorageTestSuite) TestInsertPositive() {
	category := &model.Category{ID: 3, Title: "Sport"}
	expectedCategories := append(s.InitCategories, category)

	err := s.service.Insert(category)
	require.NoError(s.T(), err)

    persistantCategories, err := s.persistantStorage.GetAll()
	require.NoError(s.T(), err)
    inmemoryCategories := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistantCategories, expectedCategories)
	assert.ElementsMatch(s.T(), inmemoryCategories, expectedCategories)
}

func (s *CategoryServiceStorageTestSuite) TestInsertNegative() {
	err := s.service.Insert(s.InitCategories[0])
    assert.ErrorContains(s.T(), err, "s.persistantStorage.Insert: e.db.Exec: UNIQUE constraint failed: category.title")
}

func (s *CategoryServiceStorageTestSuite) TestUpdatePositive() {
    cc := s.service.GetAll()
    cc[0].Title = "Taxi"

	err := s.service.Update(cc[0])
	require.NoError(s.T(), err)

    persistantCategories, err := s.persistantStorage.GetAll()
    inmemoryCategories := s.inmemoryStorage.GetAll()
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), persistantCategories, inmemoryCategories)
}

func (s *CategoryServiceStorageTestSuite) TestUpdateNegative() {
    cc := s.service.GetAll()
    cc[0].Title = "Taxi"
    cc[0].ID = 10

	err := s.service.Update(cc[0])
    assert.ErrorContains(s.T(), err, "s.persistantStorage.Update: total affected rows 0 while expected 1")
    cc[0].ID = 1 // return real id to proper teardown
}

func (s *CategoryServiceStorageTestSuite) TestDeletePositive() {
    cc := s.service.GetAll()
	expectedCategories := []*model.Category{cc[0]}

	err := s.service.Delete(cc[1])
	require.NoError(s.T(), err)

    persistantCategories, err := s.persistantStorage.GetAll()
	require.NoError(s.T(), err)
    inmemoryCategories := s.inmemoryStorage.GetAll()
	assert.ElementsMatch(s.T(), persistantCategories, expectedCategories)
	assert.ElementsMatch(s.T(), inmemoryCategories, expectedCategories)
}

func (s *CategoryServiceStorageTestSuite) TestDeleteNegative() {
    cc := s.service.GetAll()
    cc[0].ID = 10

	err := s.service.Delete(cc[0])
    assert.ErrorContains(s.T(), err, "s.persistantStorage.Delete: total affected rows 0 while expected 1")
    cc[0].ID = 1 // return real id to proper teardown
}

func (s *CategoryServiceStorageTestSuite) TestGetByID() {
    c := s.service.GetByID(2)
    assert.Equal(s.T(), s.InitCategories[1].Title, c.Title)
}

func (s *CategoryServiceStorageTestSuite) TestGetByTitle() {
    c := s.service.GetByTitle(s.InitCategories[1].Title)
    assert.Equal(s.T(), s.InitCategories[1].ID, c.ID)
}

func (s *CategoryServiceStorageTestSuite) TestGetAllSorted() {
    cc := s.service.GetAllSorted()
    expectedCategories := make([]*model.Category, len(s.InitCategories))
    copy(expectedCategories, s.InitCategories)
    sort.Sort(categoryList(expectedCategories))
	assert.ElementsMatch(s.T(), cc, expectedCategories)
}

func (s *CategoryServiceStorageTestSuite) TearDownTest() {
    for {
        cc := s.service.GetAll()
        if len(cc) == 0 {
            break
        }
        
		err := s.persistantStorage.Delete(cc[0].ID)
		require.NoError(s.T(), err, "occured in TearDownTest")
		s.inmemoryStorage.Delete(cc[0])
    }
}

func (s *CategoryServiceStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occured in TearDownSuite")
}

func TestCategoryServiceStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceStorageTestSuite))
}
