package tag

import (
	"errors"
	"fmt"

	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/emPeeGee/raffinance/pkg/util"

	"gorm.io/gorm"
)

type Repository interface {
	getTags(userId uint) ([]tagResponse, error)
	createTag(userId uint, Tag createTagDTO) (*tagResponse, error)
	updateTag(userId, tagId uint, tag updateTagDTO) (*tagResponse, error)
	deleteTag(userId, id uint) error
	tagExistsAndBelongsToUser(userID, id uint, name string) (bool, error)
	tagIsUsed(tagId uint) error
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewTagRepository(db *gorm.DB, logger log.Logger) *repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) createTag(userId uint, tag createTagDTO) (*tagResponse, error) {
	newTag := entity.Tag{
		Name:   tag.Name,
		Icon:   tag.Icon,
		Color:  tag.Color,
		UserID: &userId,
	}

	if err := r.db.Create(&newTag).Error; err != nil {
		return nil, err
	}

	r.logger.Info("new tag, ", util.StringifyAny(newTag))
	createdTag := &tagResponse{
		ID:        newTag.ID,
		Name:      newTag.Name,
		Color:     newTag.Color,
		Icon:      newTag.Icon,
		CreatedAt: newTag.CreatedAt,
		UpdatedAt: newTag.UpdatedAt,
	}

	return createdTag, nil
}

func (r *repository) updateTag(userId, tagId uint, tag updateTagDTO) (*tagResponse, error) {
	// NOTE: When update with struct, GORM will only update non-zero fields, you might want to use
	// map to update attributes or use Select to specify fields to update
	if err := r.db.Model(&entity.Tag{}).Where("id = ?", tagId).Updates(map[string]interface{}{
		"name":  tag.Name,
		"icon":  tag.Icon,
		"color": tag.Color,
	}).Error; err != nil {
		return nil, err
	}

	var updatedTag tagResponse
	if err := r.db.Model(&entity.Tag{}).First(&updatedTag, tagId).Error; err != nil {
		return nil, err
	}

	return &updatedTag, nil
}

func (r *repository) deleteTag(userId, id uint) error {
	return r.db.Delete(&entity.Tag{}, id).Error
}

func (r *repository) getTags(userId uint) ([]tagResponse, error) {
	var tags []tagResponse = make([]tagResponse, 0)
	var user entity.User

	if err := r.db.Preload("Tags").Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	r.logger.Info(user)

	for _, tag := range user.Tags {
		tags = append(tags, tagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Color:     tag.Color,
			Icon:      tag.Icon,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return tags, nil
}

// if id < 0, will search by name
func (r *repository) tagExistsAndBelongsToUser(userID, id uint, name string) (bool, error) {
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

	if err := r.db.Model(&entity.Tag{}).
		Where(whereClause, values...).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// TODO: Can this be achieved using constraints ???
func (r *repository) tagIsUsed(tagId uint) error {
	var count int64
	if err := r.db.Model(&entity.Tag{}).
		Joins("JOIN transaction_tags ON transaction_tags.tag_id = tags.id").
		Where("transaction_tags.tag_id = ?", tagId).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete tag %d that is used in %d transactions", tagId, count)
	}

	return nil
}
