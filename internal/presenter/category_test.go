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
	expected := map[string]string{"ID": "0", "Title": "Health"}
	actual := s.presenter.ToMap(category)
	assert.Equal(s.T(), expected, actual)
}

func (s *CategoryPresenterTestSuite) TestFromMapPositive() {
	for _, tc := range []struct {
		name     string
		give     map[string]string
		expected *model.Category
	}{
		{
			name:     "ExistingID",
			give:     map[string]string{"ID": "91", "Title": "Health"},
			expected: &model.Category{ID: 91, Title: "Health"},
		},
		{
			name:     "NotExistingID",
			give:     map[string]string{"Title": "Health"},
			expected: &model.Category{ID: 0, Title: "Health"},
		},
	} {
		s.Run(tc.name, func() {
			actual, err := s.presenter.FromMap(tc.give)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, actual)
		})
	}
}

func (s *CategoryPresenterTestSuite) TestFromMapNegative() {
	_, err := s.presenter.FromMap(map[string]string{"Name": "Health"})
	assert.EqualError(s.T(), err, `checkKeys: key "Title" is missing`)
}

func TestCategoryPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryPresenterTestSuite))
}
