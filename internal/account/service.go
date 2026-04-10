package account

import (
	"github.com/Gwilides/finance-tracker/internal/user"
)

type UserProvider interface {
	FindByEmail(email string) (*user.User, error)
}

type AccountServiceDeps struct {
	UserRepository    UserProvider
	AccountRepository *AccountRepository
}

type AccountService struct {
	userRepository    UserProvider
	accountRepository *AccountRepository
}

func NewAccountService(deps *AccountServiceDeps) *AccountService {
	return &AccountService{
		userRepository:    deps.UserRepository,
		accountRepository: deps.AccountRepository,
	}
}

func (service *AccountService) Create(email string, body *AccountCreateRequest) (*Account, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	account := &Account{
		UserID: user.ID,
		Type:   body.Type,
		Title:  body.Title,
	}
	_, err = service.accountRepository.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *AccountService) GetById(email string, id uint) (*Account, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	account, err := service.accountRepository.GetById(id)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	if account.UserID != user.ID {
		return nil, ErrForbidden
	}
	return account, nil
}

func (service *AccountService) Update(email string, account *Account) (*Account, error) {
	_, err := service.GetById(email, account.ID)
	if err != nil {
		return nil, err
	}
	return service.accountRepository.Update(account)
}

func (service *AccountService) Delete(email string, id uint) error {
	_, err := service.GetById(email, id)
	if err != nil {
		return err
	}
	return service.accountRepository.Delete(id)
}
