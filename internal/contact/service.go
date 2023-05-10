package contact

import (
	"fmt"

	"github.com/emPeeGee/raffinance/pkg/log"
)

type Service interface {
	createContact(userId uint, contact createContactDTO) (*contactResponse, error)
	deleteContact(userId, id uint) error
	getContacts(userId uint) ([]contactResponse, error)
	updateContact(userId, contactId uint, contact updateContactDTO) (*contactResponse, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewContactService(repo Repository, logger log.Logger) *service {
	return &service{repo: repo, logger: logger}
}

func (s *service) createContact(userId uint, contact createContactDTO) (*contactResponse, error) {
	// check if such name or email already exists, email and name should be unique per user
	err := s.repo.contactExists(contact.Name, contact.Email)
	if err != nil {
		return nil, err
	}

	return s.repo.createContact(userId, contact)
}

func (s *service) deleteContact(userId, id uint) error {
	ok, err := s.repo.contactExistsAndBelongsToUser(userId, id)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("contact with ID %d exists or does not belong to user with ID %d", id, userId)
	}

	return s.repo.deleteContact(userId, id)
}

func (s *service) getContacts(userId uint) ([]contactResponse, error) {
	return s.repo.getContacts(userId)
}

func (s *service) updateContact(userId, contactId uint, contact updateContactDTO) (*contactResponse, error) {
	// TODO: check if the contact belongs to the current user

	// check if id exists
	err := s.repo.contactExistsByID(contactId)
	if err != nil {
		return nil, err
	}

	// check if such name or email already exists, email and name should be unique per user
	err = s.repo.contactExists(contact.Name, contact.Email)
	if err != nil {
		return nil, err
	}

	return s.repo.updateContact(userId, contactId, contact)
}
