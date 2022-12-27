package inmemory

import (
	"gentlemoney/internal/model"
)

// Category is used to acces inmemory storage.
type Category struct {
	categories        []*model.Category
	categoryByID      map[int64]*model.Category
	categoryByTitle   map[string]*model.Category
}

// NewCategory returns new category inmemory storage.
func NewCategory() *Category {
	return &Category{
		categories: make([]*model.Category, 0, 20),
		categoryByID: make(map[int64]*model.Category),
		categoryByTitle: make(map[string]*model.Category),
	}
}

// Init initialize inmemory storage with given slice of data.
func (s *Category) Init(cc []*model.Category) {
	for _, c := range cc {
		s.categoryByID[c.ID] = c
		s.categoryByTitle[c.Title] = c
	}
	s.categories = cc
}

// Insert appends category to inmemory storage.
func (s *Category) Insert(c *model.Category) {
	s.categoryByID[c.ID] = c
	s.categoryByTitle[c.Title] = c
	s.categories = append(s.categories, c)
}

// Delete removes category from current inmemory storage.
func (s *Category) Delete(c *model.Category) {
	delete(s.categoryByID, c.ID)
	delete(s.categoryByTitle, c.Title)

	if len(s.categories) == 1 {
		s.categories = make([]*model.Category, 0, 20)
	} else {
		i := findIndex(c, s.categories)
		s.categories = append(s.categories[:i], s.categories[i+1:]...)
	}
}

// GetAll returns slice of categories.
func (s *Category) GetAll() []*model.Category {
	return s.categories
}

// GetByID returns category by its id.
func (s *Category) GetByID(id int64) *model.Category {
	return s.categoryByID[id]
}

// GetByID returns category by its title.
func (s *Category) GetByTitle(title string) *model.Category {
	return s.categoryByTitle[title]
}
