package category

import (
	"github.com/Gwilides/finance-tracker/pkg/db"
	"gorm.io/gorm/clause"
)

type CategoryRepository struct {
	db *db.Db
}

func NewCategoryRepository(db *db.Db) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) Create(category *Category) error {
	result := repo.db.Create(category)
	return result.Error
}

func (repo *CategoryRepository) GetById(id uint) (*Category, error) {
	var category Category
	result := repo.db.First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (repo *CategoryRepository) GetAll(userID uint) ([]Category, error) {
	var categories []Category
	result := repo.db.Find(&categories, "user_id = ? OR user_id IS NULL", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil

}

func (repo *CategoryRepository) Update(category *Category) (*Category, error) {
	result := repo.db.Clauses(clause.Returning{}).Updates(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (repo *CategoryRepository) Delete(id uint) error {
	result := repo.db.Delete(&Category{}, id)
	return result.Error
}
