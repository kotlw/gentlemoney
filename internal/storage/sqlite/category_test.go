package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CategorySqliteStorageTestSuite struct {
	suite.Suite
	db             *sql.DB
	storage        *sqlite.Category
	InitCategories []*model.Category
}

func (s *CategorySqliteStorageTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err, "occurred in SetupSuite")
	s.db = db

	s.storage, err = sqlite.NewCategory(db)
	require.NoError(s.T(), err, "occurred in SetupSuite")

	// id's settled by sqlite on insert incrementally starting from 1,
	// so here they are initialized for match purpose
	s.InitCategories = []*model.Category{
		{ID: 1, Title: "Grocery"},
		{ID: 2, Title: "Health"},
	}
}

func (s *CategorySqliteStorageTestSuite) SetupTest() {
	stmt, err := s.db.Prepare(`INSERT INTO category (title) VALUES (?);`)
	require.NoError(s.T(), err, "occurred in SetupTest")

	for _, category := range s.InitCategories {
		_, err := stmt.Exec(category.Title)
		require.NoError(s.T(), err, "occurred in SetupTest")
	}
}

func (s *CategorySqliteStorageTestSuite) TestInsertPositive() {
	category := &model.Category{ID: 3, Title: "Sport"}
	expectedCategories := append(s.InitCategories, category)

	_, err := s.storage.Insert(category)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCategories)
}

func (s *CategorySqliteStorageTestSuite) TestInsertNegative() {
	_, err := s.storage.Insert(s.InitCategories[1])
	assert.ErrorContains(s.T(), err, "e.db.Exec: UNIQUE constraint failed: category.title")
}

func (s *CategorySqliteStorageTestSuite) TestUpdatePositive() {
	expectedCategories := make([]*model.Category, len(s.InitCategories))
	copy(expectedCategories, s.InitCategories)
	expectedCategories[1].Title = "Taxi"

	err := s.storage.Update(expectedCategories[1])
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCategories)
}

func (s *CategorySqliteStorageTestSuite) TestUpdateNegative() {
	err := s.storage.Update(&model.Category{ID: 10})
	assert.ErrorContains(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *CategorySqliteStorageTestSuite) TestDeletePositive() {
	expectedCategories := []*model.Category{s.InitCategories[0]}

	err := s.storage.Delete(2)
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), s.fetchActualData(), expectedCategories)
}

func (s *CategorySqliteStorageTestSuite) TestDeleteNegative() {
	err := s.storage.Delete(10)
	assert.EqualError(s.T(), err, "total affected rows 0 while expected 1")
}

func (s *CategorySqliteStorageTestSuite) TestGetAll() {
	actualCategories, err := s.storage.GetAll()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), actualCategories, s.InitCategories)
}

func (s *CategorySqliteStorageTestSuite) fetchActualData() []*model.Category {
	rows, err := s.db.Query(`SELECT id, title FROM category;`)
	require.NoError(s.T(), err)
	defer func() {
		err = rows.Close()
		require.NoError(s.T(), err)
	}()

	res := make([]*model.Category, 0, 3)
	for rows.Next() {
		t := model.NewEmptyCategory()
		err = rows.Scan(&t.ID, &t.Title)
		require.NoError(s.T(), err)
		res = append(res, t)
	}

	return res
}

func (s *CategorySqliteStorageTestSuite) TearDownTest() {
	_, err := s.db.Exec(`DELETE FROM category;`)
	require.NoError(s.T(), err, "occurred in TearDownTest")
}

func (s *CategorySqliteStorageTestSuite) TearDownSuite() {
	err := s.db.Close()
	require.NoError(s.T(), err, "occurred in TearDownSuite")
}

func TestCategorySqliteStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CategorySqliteStorageTestSuite))
}
