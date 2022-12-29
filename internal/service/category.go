package service

import (
	"fmt"
	"sort"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"
	"github.com/kotlw/gentlemoney/internal/storage/sqlite"
)

// categoryList is a wrapper to perform sort by model.Category.Title.
type categoryList []*model.Category

func (cc categoryList) Len() int           { return len(cc) }
func (cc categoryList) Less(i, j int) bool { return cc[i].Title < cc[j].Title }
func (cc categoryList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

// Category service contains business logic related to model.Category.
type Category struct {
	persistentStorage *sqlite.Category
	inmemoryStorage   *inmemory.Category
}

// NewCategory returns Category service.
func NewCategory(persistentStorage *sqlite.Category, inmemoryStorage *inmemory.Category) (*Category, error) {
	c := &Category{
		persistentStorage: persistentStorage,
		inmemoryStorage:   inmemoryStorage,
	}

	if err := c.Init(); err != nil {
		return nil, fmt.Errorf("c.Init: %w", err)
	}

	return c, nil
}

// Init initialize inmemory storage with data from persistent storage.
func (s *Category) Init() error {
	cc, err := s.persistentStorage.GetAll()
	if err != nil {
		return fmt.Errorf("s.persistentStorage.GetAll: %w", err)
	}

	s.inmemoryStorage.Init(cc)

	return nil
}

// Insert appends category to both persistent and inmemory storages.
func (s *Category) Insert(c *model.Category) error {
	id, err := s.persistentStorage.Insert(c)
	if err != nil {
		return fmt.Errorf("s.persistentStorage.Insert: %w", err)
	}

	c.ID = id
	s.inmemoryStorage.Insert(c)

	return nil
}

// Update updates category in persistent storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistent storage as well.
func (s *Category) Update(c *model.Category) error {
	if err := s.persistentStorage.Update(c); err != nil {
		return fmt.Errorf("s.persistentStorage.Update: %w", err)
	}
	return nil
}

// Delete deletes category from inmemory and persistent storages.
func (s *Category) Delete(c *model.Category) error {
	if err := s.persistentStorage.Delete(c.ID); err != nil {
		return fmt.Errorf("s.persistentStorage.Delete: %w", err)
	}
	s.inmemoryStorage.Delete(c)
	return nil
}

// GetAll returns all categories.
func (s *Category) GetAll() []*model.Category {
	return s.inmemoryStorage.GetAll()
}

// GetByID returns category by given model.Category.ID.
func (s *Category) GetByID(id int64) *model.Category {
	return s.inmemoryStorage.GetByID(id)
}

// GetByTitle returns category by given model.Category.Title.
func (s *Category) GetByTitle(title string) *model.Category {
	return s.inmemoryStorage.GetByTitle(title)
}

// GetAllSorted returns all categories sorted by model.Category.Title.
func (s *Category) GetAllSorted() []*model.Category {
	cc := s.inmemoryStorage.GetAll()
	sort.Sort(categoryList(cc))
	return cc
}