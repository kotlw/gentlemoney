package sqlite

import (
	"database/sql"
	"fmt"

	"gentlemoney/internal/model"
)

// Transaction is used to acces the persistant storage.
type Transaction struct {
	executor executor[model.Transaction]
}

// NewTransaction returns new transaction storage.
func NewTransaction(db *sql.DB) (*Transaction, error) {
	s := &Transaction{executor[model.Transaction]{db}}

	if err := s.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("s.CreateTableIfNotExists: %w", err)
	}

	return s, nil
}

// CreateTableIfNotExists creates transaction table if not exists.
func (s *Transaction) CreateTableIfNotExists() error {
	q := `CREATE TABLE IF NOT EXISTS "transaction"(
            id INTEGER PRIMARY KEY,
            date DATETIME NOT NULL,
            amount INTEGER NOT NULL,
            note TEXT,
            accountId INTEGER NOT NULL,
            categoryId INTEGER NOT NULL,
            FOREIGN KEY(accountId) REFERENCES account(id),
            FOREIGN KEY(categoryId) REFERENCES category(id));`
	_, err := s.executor.db.Exec(q)
	return err
}

// Insert transaction into persistant storage.
func (s *Transaction) Insert(t *model.Transaction) (int64, error) {
	return s.executor.insert(
		`INSERT INTO "transaction" (date, amount, note, accountId, categoryId) VALUES (?, ?, ?, ?, ?);`,
		t.Date, t.Amount, t.Note, t.Account.ID, t.Category.ID)
}

// Update transaction in persistand storage.
func (s *Transaction) Update(t *model.Transaction) error {
	return s.executor.update(
		`UPDATE "transaction" SET date = ?, amount = ?, note = ?, accountId = ?, categoryId = ? WHERE id = ?;`,
		t.Date, t.Amount, t.Note, t.Account.ID, t.Category.ID, t.ID)
}

// Delete transaction from persistant storage.
func (s *Transaction) Delete(id int64) error {
	return s.executor.update(`DELETE FROM "transaction" WHERE id = ?;`, id)
}

// GetAll transaction from persistant storage.
func (s *Transaction) GetAll() ([]*model.Transaction, error) {
	return s.executor.getAll(`SELECT id, date, amount, note, accountId, categoryId FROM "transaction";`,
		func() (*model.Transaction, []any) {
			t := model.NewEmptyTransaction()
			return t, []any{&t.ID, &t.Date, &t.Amount, &t.Note, &t.Account.ID, &t.Category.ID}
		})
}
