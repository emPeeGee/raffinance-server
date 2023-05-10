package category

import (
	"fmt"

	"github.com/emPeeGee/raffinance/pkg/log"
)

type Service interface {
	createCategory(userId uint, category createCategoryDTO) (*categoryResponse, error)
	deleteCategory(userId, id uint) error
	getCategories(userId uint) ([]categoryResponse, error)
	updateCategory(usedId, categoryId uint, category updateCategoryDTO) (*categoryResponse, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewCategoryService(repo Repository, logger log.Logger) *service {
	return &service{repo: repo, logger: logger}
}

func (s *service) createCategory(userID uint, category createCategoryDTO) (*categoryResponse, error) {
	if err := checkForBlacklist(category.Name); err != nil {
		return nil, err
	}

	// name should be unique per user
	exists, err := s.repo.categoryExistsAndBelongsToUser(userID, 0, category.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("category with name %s exists", category.Name)
	}

	return s.repo.createCategory(userID, category)
}

func (s *service) deleteCategory(userId, id uint) error {
	ok, err := s.repo.categoryExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("category with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	if err := s.repo.categoryIsUsed(id); err != nil {
		return err
	}

	return s.repo.deleteCategory(userId, id)
}

func (s *service) getCategories(userId uint) ([]categoryResponse, error) {
	return s.repo.getCategories(userId)
}

func (s *service) updateCategory(userID, categoryId uint, category updateCategoryDTO) (*categoryResponse, error) {
	if err := checkForBlacklist(category.Name); err != nil {
		return nil, err
	}

	// TODO: bug, when user updates only icon or color. the category won't be updated
	exists, err := s.repo.categoryExistsAndBelongsToUser(userID, categoryId, "")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("category with ID %d does not exist or belong to user with ID %d", categoryId, userID)
	}

	// check if such name exists, name should be unique per user
	exists, err = s.repo.categoryExistsAndBelongsToUser(userID, 0, category.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("category with name %s already exists", category.Name)
	}

	return s.repo.updateCategory(userID, categoryId, category)
}
