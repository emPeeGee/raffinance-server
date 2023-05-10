package category

import (
	"errors"
	"fmt"

	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/emPeeGee/raffinance/pkg/util"

	"gorm.io/gorm"
)

type Repository interface {
	getCategories(userId uint) ([]categoryResponse, error)
	createCategory(userId uint, category createCategoryDTO) (*categoryResponse, error)
	updateCategory(userId, categoryId uint, category updateCategoryDTO) (*categoryResponse, error)
	deleteCategory(userId, id uint) error
	categoryExistsAndBelongsToUser(userID, id uint, name string) (bool, error)
	categoryIsUsed(categoryId uint) error
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewCategoryRepository(db *gorm.DB, logger log.Logger) *repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) createCategory(userId uint, category createCategoryDTO) (*categoryResponse, error) {
	newCategory := entity.Category{
		Name:   category.Name,
		Color:  category.Color,
		Icon:   category.Icon,
		UserID: &userId,
	}

	if err := r.db.Create(&newCategory).Error; err != nil {
		return nil, err
	}

	r.logger.Info("new category, ", util.StringifyAny(newCategory))
	createdCategory := &categoryResponse{
		ID:        newCategory.ID,
		Name:      newCategory.Name,
		Color:     newCategory.Color,
		Icon:      newCategory.Icon,
		CreatedAt: newCategory.CreatedAt,
		UpdatedAt: newCategory.UpdatedAt,
	}

	return createdCategory, nil
}

func (r *repository) updateCategory(userId, categoryId uint, category updateCategoryDTO) (*categoryResponse, error) {
	// NOTE: When update with struct, GORM will only update non-zero fields, you might want to use
	// map to update attributes or use Select to specify fields to update
	if err := r.db.Model(&entity.Category{}).Where("id = ?", categoryId).Updates(map[string]interface{}{
		"name":  category.Name,
		"color": category.Color,
		"icon":  category.Icon,
	}).Error; err != nil {
		return nil, err
	}

	var updatedCategory categoryResponse
	if err := r.db.Model(&entity.Category{}).First(&updatedCategory, categoryId).Error; err != nil {
		return nil, err
	}

	return &updatedCategory, nil
}

func (r *repository) deleteCategory(userId, id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

func (r *repository) getCategories(userId uint) ([]categoryResponse, error) {
	var categories []categoryResponse = make([]categoryResponse, 0)
	var user entity.User

	if err := r.db.Preload("Categories").Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	for _, category := range user.Categories {
		categories = append(categories, categoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Color:     category.Color,
			Icon:      category.Icon,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}

	return categories, nil
}

func (r *repository) categoryExistsAndBelongsToUser(userID, id uint, name string) (bool, error) {
	var count int64
	var whereClause string
	var values []interface{}

	whereClause = "user_id = ?"
	values = append(values, userID)

	if id > 0 {
		whereClause += " AND id = ?"
		values = append(values, id)
	} else if name != "" {
		whereClause += " AND name = ?"
		values = append(values, name)
	} else {
		return false, errors.New("id or name parameter is required")
	}

	if err := r.db.Model(&entity.Category{}).
		Where(whereClause, values...).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// TODO: Can this be achieved using constraints ???
func (r *repository) categoryIsUsed(categoryId uint) error {
	var count int64
	if err := r.db.Model(&entity.Category{}).
		Joins("JOIN transactions ON transactions.category_id = categories.id").
		Where("categories.id = ?", categoryId).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category %d that is used in %d transactions", categoryId, count)
	}

	return nil
}
