package category

import (
	"github.com/Gwilides/finance-tracker/internal/user"
	"gorm.io/gorm"
)

type UserProvider interface {
	FindByEmail(email string) (*user.User, error)
}

type CategoryServiceDeps struct {
	UserRepository     UserProvider
	CategoryRepository *CategoryRepository
}

type CategoryService struct {
	userRepository     UserProvider
	categoryRepository *CategoryRepository
}

func NewCategoryService(deps *CategoryServiceDeps) *CategoryService {
	return &CategoryService{
		userRepository:     deps.UserRepository,
		categoryRepository: deps.CategoryRepository,
	}
}

func (service *CategoryService) Create(email string, body *CategoryCreateRequest) (*Category, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	category := &Category{
		UserID: &user.ID,
		Title:  body.Title,
	}
	err = service.categoryRepository.Create(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (service *CategoryService) GetAll(email string) ([]Category, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	categories, err := service.categoryRepository.GetAll(user.ID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (service *CategoryService) GetById(email string, id uint) (*Category, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	category, err := service.categoryRepository.GetById(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	if category.UserID == nil {
		return category, nil
	}
	if *category.UserID != user.ID {
		return nil, ErrForbidden
	}
	return category, nil
}

func (service *CategoryService) Update(email string, id uint, body *CategoryUpdateRequest) (*Category, error) {
	_, err := service.GetById(email, id)
	if err != nil {
		return nil, err
	}
	return service.categoryRepository.Update(&Category{
		Model: gorm.Model{ID: id},
		Title: body.Title,
	})
}

func (service *CategoryService) Delete(email string, id uint) error {
	_, err := service.GetById(email, id)
	if err != nil {
		return err
	}
	return service.categoryRepository.Delete(id)
}
