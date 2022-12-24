package storage

import (
	"database/sql"
	"fmt"

	"gentlemoney/internal/model"
)

// Account is used to acces the persistant storage.
type Account struct {
	executor executor[model.Account]
}

// NewAccount returns new account storage.
func NewAccount(db *sql.DB) (*Account, error) {
	s := &Account{executor[model.Account]{db}}

	if err := s.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("s.CreateTableIfNotExists: %w", err)
	}

	return s, nil
}

// CreateTableIfNotExists creates account table if not exists.
func (s *Account) CreateTableIfNotExists() error {
	q := `CREATE TABLE IF NOT EXISTS account(
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL UNIQUE,
            currencyId INTEGER NOT NULL,
            FOREIGN KEY(currencyId) REFERENCES currency(id));`
	_, err := s.executor.db.Exec(q)
	return err
}

// Insert account into persistant storage.
func (s *Account) Insert(a *model.Account) (int64, error) {
	return s.executor.insert(
		`INSERT INTO account(name, currencyId) VALUES (?, ?);`,
		a.Name, a.Currency.ID)
}

// Update account in persistand storage.
func (s *Account) Update(a *model.Account) error {
	return s.executor.update(
		`UPDATE account SET name = ?, currencyId = ? WHERE id = ?;`,
		a.Name, a.Currency.ID, a.ID)
}

// Delete account from persistant storage.
func (s *Account) Delete(id int64) error {
	return s.executor.update(`DELETE FROM account WHERE id = ?;`, id)
}

// GetAll accounts from persistant storage.
func (s *Account) GetAll() ([]*model.Account, error) {
	return s.executor.getAll(`SELECT id, name, currencyId FROM account;`,
		func() (*model.Account, []any) {
			t := model.NewEmptyAccount()
			return t, []any{&t.ID, &t.Name, &t.Currency.ID}
		})
}
