package tag

import (
	"fmt"

	"github.com/emPeeGee/raffinance/pkg/log"
)

type Service interface {
	createTag(userId uint, tag createTagDTO) (*tagResponse, error)
	deleteTag(userId, id uint) error
	getTags(userId uint) ([]tagResponse, error)
	updateTag(usedId, tagId uint, tag updateTagDTO) (*tagResponse, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewTagService(repo Repository, logger log.Logger) *service {
	return &service{repo: repo, logger: logger}
}

func (s *service) createTag(userID uint, tag createTagDTO) (*tagResponse, error) {
	// check if such name or email already exists, email and name should be unique per user
	exists, err := s.repo.tagExistsAndBelongsToUser(userID, 0, tag.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("tag with name %s exists", tag.Name)
	}

	return s.repo.createTag(userID, tag)
}

func (s *service) deleteTag(userId, id uint) error {
	// another user hasn't to be able to access tags
	ok, err := s.repo.tagExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("tag with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	if err := s.repo.tagIsUsed(id); err != nil {
		return err
	}

	return s.repo.deleteTag(userId, id)
}

func (s *service) getTags(userId uint) ([]tagResponse, error) {
	return s.repo.getTags(userId)
}

func (s *service) updateTag(userId, tagId uint, tag updateTagDTO) (*tagResponse, error) {
	exists, err := s.repo.tagExistsAndBelongsToUser(userId, tagId, "")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("tag with ID %d does not exist or belong to user with ID %d", tagId, userId)
	}

	// check if such name exists, name should be unique per user
	exists, err = s.repo.tagExistsAndBelongsToUser(userId, 0, tag.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("tag with name %s already exists", tag.Name)
	}

	return s.repo.updateTag(userId, tagId, tag)
}
