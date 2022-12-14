package inmemory

// InmemoryStorage is a facade structure which aggregates all inmemory storages. It is used for convenience.
type InmemoryStorage struct {
	category    *Category
	currency    *Currency
	account     *Account
	transaction *Transaction
}

// New returns new InmemoryStorage.
func New() *InmemoryStorage {
	return &InmemoryStorage{
		category:    NewCategory(),
		currency:    NewCurrency(),
		account:     NewAccount(),
		transaction: NewTransaction(),
	}
}

// Category returns category inmemory storage.
func (s *InmemoryStorage) Category() *Category {
	return s.category
}

// Currency returns currency inmemory storage.
func (s *InmemoryStorage) Currency() *Currency {
	return s.currency
}

// Account returns account inmemory storage.
func (s *InmemoryStorage) Account() *Account {
	return s.account
}

// Transaction returns transaction inmemory storage.
func (s *InmemoryStorage) Transaction() *Transaction {
	return s.transaction
}
