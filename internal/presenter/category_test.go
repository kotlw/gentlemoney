package presenter_test

import (
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/presenter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CategoryPresenterTestSuite struct {
	suite.Suite
	presenter *presenter.Category
}

func (s *CategoryPresenterTestSuite) SetupSuite() {
	s.presenter = presenter.NewCategory()
}

func (s *CategoryPresenterTestSuite) TestToMap() {
	category := &model.Category{Title: "Health"}
	expected := map[string]string{"Title": "Health"}
	actual := s.presenter.ToMap(category)
	assert.Equal(s.T(), expected, actual)
}

func (s *CategoryPresenterTestSuite) TestFromMapPositive() {
	m := map[string]string{"Title": "Health"}
	expected := &model.Category{Title: "Health"}
	actual, err := s.presenter.FromMap(m)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expected, actual)
}

func (s *CategoryPresenterTestSuite) TestFromMapNegative() {
	_, err := s.presenter.FromMap(map[string]string{"Name": "Health"})
	assert.EqualError(s.T(), err, `checkKeys: key "Title" is missing`)
}

func TestCategoryPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryPresenterTestSuite))
}
