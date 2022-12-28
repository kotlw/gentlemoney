package inmemory_test

import (
	"testing"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage/inmemory"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoryInmemoryStorageTestSuite struct {
	suite.Suite
	storage        *inmemory.Category
	InitCategories []*model.Category
}

func (s *CategoryInmemoryStorageTestSuite) SetupSuite() {
	s.storage = inmemory.NewCategory()
	s.InitCategories = []*model.Category{
		{ID: 1, Title: "Grocery"},
		{ID: 2, Title: "Health"},
	}
}

func (s *CategoryInmemoryStorageTestSuite) SetupTest() {
	s.storage.Init(s.InitCategories)
}

func (s *CategoryInmemoryStorageTestSuite) TestInsertPositive() {
	category := &model.Category{ID: 3, Title: "Sport"}
	expectedCategories := append(s.InitCategories, category)

	s.storage.Insert(category)

	assert.ElementsMatch(s.T(), s.storage.GetAll(), expectedCategories)
	assert.Equal(s.T(), category, s.storage.GetByID(category.ID))
	assert.Equal(s.T(), category, s.storage.GetByTitle(category.Title))
}

func (s *CategoryInmemoryStorageTestSuite) TestDeletePositive() {
	s.storage.Delete(s.InitCategories[1])

	assert.ElementsMatch(s.T(), s.storage.GetAll(), s.InitCategories[:1])
}

func (s *CategoryInmemoryStorageTestSuite) TearDownTest() {
    for {
        cc := s.storage.GetAll()
        if len(cc) == 0 {
            break
        }
        s.storage.Delete(cc[0])
    }
}

func TestCategoryInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryInmemoryStorageTestSuite))
}
