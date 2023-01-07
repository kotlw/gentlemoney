package inmemory

import (
	"github.com/kotlw/gentlemoney/internal/model"
)

// Transaction is used to acces inmemory storage.
type Transaction struct {
	transactions    []*model.Transaction
	transactionByID map[int64]*model.Transaction
}

// NewTransaction returns new transaction inmemory storage.
func NewTransaction() *Transaction {
	return &Transaction{
		transactions:    make([]*model.Transaction, 0, 20),
		transactionByID: make(map[int64]*model.Transaction),
	}
}

// Init initialize inmemory storage with given slice of data.
func (s *Transaction) Init(tt []*model.Transaction) {
	for _, t := range tt {
		s.transactionByID[t.ID] = t
	}
	s.transactions = tt
}

// Insert appends transaction to inmemory storage.
func (s *Transaction) Insert(t *model.Transaction) {
	s.transactionByID[t.ID] = t
	s.transactions = append(s.transactions, t)
}

// Delete removes transaction from current inmemory storage.
func (s *Transaction) Delete(t *model.Transaction) {
	delete(s.transactionByID, t.ID)

	for i, tt := range s.transactions {
		if tt.ID == t.ID {
			last := len(s.transactions) - 1
			s.transactions[i] = s.transactions[last]
			s.transactions = s.transactions[:last]
		}
	}
}

// GetAll returns slice of transactions.
func (s *Transaction) GetAll() []*model.Transaction {
	return s.transactions
}

// GetByID returns transaction by its id.
func (s *Transaction) GetByID(id int64) *model.Transaction {
	return s.transactionByID[id]
}
