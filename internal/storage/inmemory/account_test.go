package inmemory_test

import (
	"testing"

	"github.com/kotlw/gentlemoney/internal/model"
	"github.com/kotlw/gentlemoney/internal/storage/inmemory"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccountInmemoryStorageTestSuite struct {
	suite.Suite
	storage        *inmemory.Account
	InitAccounts []*model.Account
}

func (s *AccountInmemoryStorageTestSuite) SetupSuite() {
	s.storage = inmemory.NewAccount()
	s.InitAccounts = []*model.Account{
		{ID: 1, Name: "Card1", Currency: model.NewEmptyCurrency()},
		{ID: 2, Name: "Card2", Currency: model.NewEmptyCurrency()},
	}
}

func (s *AccountInmemoryStorageTestSuite) SetupTest() {
	s.storage.Init(s.InitAccounts)
}

func (s *AccountInmemoryStorageTestSuite) TestInsertPositive() {
	account := &model.Account{ID: 3, Name: "Card3", Currency: model.NewEmptyCurrency()}
	expectedAccounts := append(s.InitAccounts, account)

	s.storage.Insert(account)

	assert.ElementsMatch(s.T(), s.storage.GetAll(), expectedAccounts)
	assert.Equal(s.T(), account, s.storage.GetByID(account.ID))
	assert.Equal(s.T(), account, s.storage.GetByName(account.Name))
}

func (s *AccountInmemoryStorageTestSuite) TestDeletePositive() {
	s.storage.Delete(s.InitAccounts[1])
	assert.ElementsMatch(s.T(), s.storage.GetAll(), s.InitAccounts[:1])
}

func (s *AccountInmemoryStorageTestSuite) TearDownTest() {
    for {
        aa := s.storage.GetAll()
        if len(aa) == 0 {
            break
        }
        s.storage.Delete(aa[0])
    }
}

func TestAccountInmemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(AccountInmemoryStorageTestSuite))
}
