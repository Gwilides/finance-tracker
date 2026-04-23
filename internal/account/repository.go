package account

import (
	"github.com/Gwilides/finance-tracker/pkg/db"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	db *db.Db
}

func NewAccountRepository(db *db.Db) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (repo *AccountRepository) Create(account *Account) (*Account, error) {
	result := repo.db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (repo *AccountRepository) GetById(id uint) (*Account, error) {
	var account Account
	result := repo.db.First(&account, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (repo *AccountRepository) Update(account *Account) (*Account, error) {
	result := repo.db.Clauses(clause.Returning{}).Updates(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (repo *AccountRepository) Delete(id uint) error {
	result := repo.db.Delete(&Account{}, id)
	return result.Error
}
