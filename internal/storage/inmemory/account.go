package inmemory

import (
	"github.com/kotlw/gentlemoney/internal/model"
)

// Account is used to acces inmemory storage.
type Account struct {
	accounts      []*model.Account
	accountByID   map[int64]*model.Account
	accountByName map[string]*model.Account
}

// NewAccount returns new account inmemory storage.
func NewAccount() *Account {
	return &Account{
		accounts:      make([]*model.Account, 0, 20),
		accountByID:   make(map[int64]*model.Account),
		accountByName: make(map[string]*model.Account),
	}
}

// Init initialize inmemory storage with given slice of data.
func (s *Account) Init(aa []*model.Account) {
	for _, a := range aa {
		s.accountByID[a.ID] = a
		s.accountByName[a.Name] = a
	}
	s.accounts = aa
}

// Insert appends account to inmemory storage.
func (s *Account) Insert(a *model.Account) {
	s.accountByID[a.ID] = a
	s.accountByName[a.Name] = a
	s.accounts = append(s.accounts, a)
}

// Update updates currency of inmemory storage.
func (s *Account) Update(a *model.Account) {
	delete(s.accountByName, s.accountByID[a.ID].Name)

	s.accountByID[a.ID] = a
	s.accountByName[a.Name] = a

	for i, aa := range s.accounts {
		if aa.ID == a.ID {
			s.accounts[i] = a
			return
		}
	}
}

// Delete removes account from current inmemory storage.
func (s *Account) Delete(a *model.Account) {
	delete(s.accountByID, a.ID)
	delete(s.accountByName, a.Name)

	for i, aa := range s.accounts {
		if aa.ID == a.ID {
			last := len(s.accounts) - 1
			s.accounts[i] = s.accounts[last]
			s.accounts = s.accounts[:last]
		}
	}
}

// GetAll returns slice of accounts.
func (s *Account) GetAll() []*model.Account {
	return s.accounts
}

// GetByID returns account by its id.
func (s *Account) GetByID(id int64) *model.Account {
	return s.accountByID[id]
}

// GetByName returns account by its name.
func (s *Account) GetByName(name string) *model.Account {
	return s.accountByName[name]
}
