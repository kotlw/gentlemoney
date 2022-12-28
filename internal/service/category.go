package service

import (
	"fmt"
	"sort"

	"gentlemoney/internal/model"
	"gentlemoney/internal/storage/inmemory"
	"gentlemoney/internal/storage/sqlite"
)

// categoryList is a wrapper to perform sort by model.Category.Title.
type categoryList []*model.Category

func (cc categoryList) Len() int           { return len(cc) }
func (cc categoryList) Less(i, j int) bool { return cc[i].Title < cc[j].Title }
func (cc categoryList) Swap(i, j int)      { cc[i], cc[j] = cc[j], cc[i] }

// Category service contains business logic related to model.Category.
type Category struct {
	persistantStorage *sqlite.Category
	inmemoryStorage   *inmemory.Category
}

// NewCategory returns Category service.
func NewCategory(persistantStorage *sqlite.Category, inmemoryStorage *inmemory.Category) (*Category, error) {
    c := &Category{
        persistantStorage: persistantStorage,
        inmemoryStorage: inmemoryStorage,
    }

    if err := c.Init(); err != nil {
        return nil, fmt.Errorf("c.Init: %w", err)
    }

    return c, nil
}

// Init initialize inmemory storage with data from persistant storage.
func (s *Category) Init() error {
    cc, err := s.persistantStorage.GetAll()
    if err != nil {
        return fmt.Errorf("s.persistantStorage.GetAll: %w", err)
    }

    s.inmemoryStorage.Init(cc)

    return nil
}

// Insert appends category to both persistant and inmemory storages.
func (s *Category) Insert(c *model.Category) error {
    id, err := s.persistantStorage.Insert(c)
    if err != nil {
        return fmt.Errorf("s.persistantStorage.Insert: %w", err)
    }
    
    c.ID = id
    s.inmemoryStorage.Insert(c)

    return nil
}

// Update updates category in persistant storage. Since GetAll returns pointers to inmemory data
// after update the category we need to update it in persistant storage as well.
func (s *Category) Update(c *model.Category) error {
    if err := s.persistantStorage.Update(c); err != nil {
        return fmt.Errorf("s.persistantStorage.Update: %w", err)
    }
    return nil
}

// Delete deletes category from inmemory and persistant storages.
func (s *Category) Delete(c *model.Category) error {
    if err := s.persistantStorage.Delete(c.ID); err != nil {
        return fmt.Errorf("s.persistantStorage.Delete: %w", err)
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
