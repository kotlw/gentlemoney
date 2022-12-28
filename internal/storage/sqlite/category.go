package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
)

// Category is used to acces the persistent storage.
type Category struct {
	executor executor[model.Category]
}

// NewCategory returns new category storage.
func NewCategory(db *sql.DB) (*Category, error) {
	s := &Category{executor[model.Category]{db}}

	if err := s.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("s.CreateTableIfNotExists: %w", err)
	}

	return s, nil
}

// CreateTableIfNotExists creates category table if not exists.
func (s *Category) CreateTableIfNotExists() error {
	q := `CREATE TABLE IF NOT EXISTS category(
            id INTEGER PRIMARY KEY,
            title TEXT NOT NULL UNIQUE);`
	_, err := s.executor.db.Exec(q)
	return err
}

// Insert category into persistent storage.
func (s *Category) Insert(c *model.Category) (int64, error) {
	return s.executor.insert(`INSERT INTO category (title) VALUES (?);`, c.Title)
}

// Update category in persistand storage.
func (s *Category) Update(c *model.Category) error {
	return s.executor.update(`UPDATE category SET title = ? WHERE id = ?;`, c.Title, c.ID)
}

// Delete category from persistent storage.
func (s *Category) Delete(id int64) error {
	return s.executor.update(`DELETE FROM category WHERE id = ?;`, id)
}

// GetAll categories from persistent storage.
func (s *Category) GetAll() ([]*model.Category, error) {
	return s.executor.getAll(`SELECT id, title FROM category;`,
		func() (*model.Category, []any) {
			t := model.NewEmptyCategory()
			return t, []any{&t.ID, &t.Title}
		})
}
