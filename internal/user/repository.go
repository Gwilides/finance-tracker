package user

import "github.com/Gwilides/finance-tracker/pkg/db"

type UserRepository struct {
	db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Create(user *User) error {
	result := repo.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
