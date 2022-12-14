package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/kotlw/gentlemoney/internal/model"
)

// Currency is used to acces the persistent storage.
type Currency struct {
	executor executor[model.Currency]
}

// NewCurrency returns new currency storage.
func NewCurrency(db *sql.DB) (*Currency, error) {
	s := &Currency{executor[model.Currency]{db}}

	if err := s.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("s.CreateTableIfNotExists: %w", err)
	}

	return s, nil
}

// CreateTableIfNotExists creates currency table if not exists.
func (s *Currency) CreateTableIfNotExists() error {
	q := `CREATE TABLE IF NOT EXISTS currency(
            id INTEGER PRIMARY KEY,
            abbreviation TEXT NOT NULL UNIQUE);`
	_, err := s.executor.db.Exec(q)
	return err
}

// Insert currency into persistent storage.
func (s *Currency) Insert(c *model.Currency) (int64, error) {
	return s.executor.insert(`INSERT INTO currency(abbreviation) VALUES (?);`, c.Abbreviation)
}

// Update currency in persistand storage.
func (s *Currency) Update(c *model.Currency) error {
	return s.executor.update(
		`UPDATE currency SET abbreviation = ? WHERE id = ?;`, c.Abbreviation, c.ID)
}

// Delete currency from persistent storage.
func (s *Currency) Delete(id int64) error {
	return s.executor.update(`DELETE FROM currency WHERE id = ?;`, id)
}

// GetAll currency from persistent storage.
func (s *Currency) GetAll() ([]*model.Currency, error) {
	return s.executor.getAll(`SELECT id, abbreviation FROM currency;`,
		func() (*model.Currency, []any) {
			t := model.NewEmptyCurrency()
			return t, []any{&t.ID, &t.Abbreviation}
		})
}
