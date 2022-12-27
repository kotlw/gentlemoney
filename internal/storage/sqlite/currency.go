package sqlite

import (
	"database/sql"
	"fmt"

	"gentlemoney/internal/model"
)

// Currency is used to acces the persistant storage.
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
            abbreviation TEXT NOT NULL UNIQUE,
            exchangeRate INTEGER,
            isMain INTEGER NOT NULL);`
	_, err := s.executor.db.Exec(q)
	return err
}

// Insert currency into persistant storage.
func (s *Currency) Insert(c *model.Currency) (int64, error) {
	return s.executor.insert(
		`INSERT INTO currency(abbreviation, exchangeRate, isMain) VALUES (?, ?, ?);`,
		c.Abbreviation, c.ExchangeRate, c.IsMain)
}

// Update currency in persistand storage.
func (s *Currency) Update(c *model.Currency) error {
	return s.executor.update(
		`UPDATE currency SET abbreviation = ?, exchangeRate = ?, isMain = ? WHERE id = ?;`,
		c.Abbreviation, c.ExchangeRate, c.IsMain, c.ID)
}

// Delete currency from persistant storage.
func (s *Currency) Delete(id int64) error {
	return s.executor.update(`DELETE FROM currency WHERE id = ?;`, id)
}

// GetAll currency from persistant storage.
func (s *Currency) GetAll() ([]*model.Currency, error) {
	return s.executor.getAll(`SELECT id, abbreviation, exchangeRate, isMain FROM currency;`,
		func() (*model.Currency, []any) {
			t := model.NewEmptyCurrency()
			return t, []any{&t.ID, &t.Abbreviation, &t.ExchangeRate, &t.IsMain}
		})
}
