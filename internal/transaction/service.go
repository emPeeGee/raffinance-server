package transaction

import (
	"fmt"
	"time"

	"github.com/emPeeGee/raffinance/internal/category"
	"github.com/emPeeGee/raffinance/internal/hub"
	"github.com/emPeeGee/raffinance/pkg/log"
)

type Service interface {
	createTransaction(userId uint, transaction CreateTransactionDTO) (*TransactionResponse, error)
	// TODO: They are not validated, validation is in handler
	CreateAdjustmentTransaction(userId, accountId uint, amount float64, trType TransactionType) (*TransactionResponse, error)
	CreateInitialTransaction(userId, accountId uint, amount float64) (*TransactionResponse, error)
	deleteTransaction(userId, id uint) error
	updateTransaction(usedId, transactionId uint, transaction UpdateTransactionDTO) (*TransactionResponse, error)
	getTransaction(userID, txnId uint) (*TransactionResponse, error)
	GetAccountTransactionsByMonth(accountId uint, year int, month time.Month) ([]TransactionResponse, error)
	GetTransactionsByFilter(filter TransactionFilter) ([]TransactionResponse, error)
	getTransactions(userId uint) ([]TransactionResponse, error)
}

type service struct {
	repo   Repository
	logger log.Logger
	hub    *hub.Hub
}

func NewTransactionService(repo Repository, logger log.Logger, hub *hub.Hub) *service {
	return &service{repo: repo, logger: logger, hub: hub}
}

func (s *service) createTransaction(userId uint, transaction CreateTransactionDTO) (*TransactionResponse, error) {
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, transaction.ToAccountID)
	if err != nil || !ok {
		return nil, fmt.Errorf("toAccountId with id %d doesn't exist or belong to user", transaction.ToAccountID)
	}

	if transaction.TransactionTypeID == byte(TRANSFER) {
		ok, err := s.repo.accountExistsAndBelongsToUser(userId, *transaction.FromAccountID)
		if err != nil || !ok {
			return nil, fmt.Errorf("fromAccountId with id %d doesn't exist or belong to user", *transaction.FromAccountID)
		}
	}

	// Omit validation if it is system category
	if transaction.CategoryID != category.SystemCategoryID {
		exists, err := s.repo.categoryExistsAndBelongsToUser(userId, transaction.CategoryID)
		if err != nil || !exists {
			return nil, fmt.Errorf("categoryId with id %d doesn't exist or belong to user", transaction.CategoryID)
		}
	}

	exist, err := s.repo.tagsExistsAndBelongsToUser(userId, transaction.TagIDs)
	if err != nil || !exist {
		return nil, fmt.Errorf("not all tags belong to user or do not exist %v", transaction.TagIDs)
	}

	return s.repo.createTransaction(userId, transaction)
}

func (s *service) CreateInitialTransaction(userId, accountId uint, amount float64) (*TransactionResponse, error) {
	transaction := CreateTransactionDTO{
		Date:              time.Now(),
		Amount:            amount,
		Description:       "Initial balance",
		Location:          "",
		ToAccountID:       accountId,
		CategoryID:        category.SystemCategoryID,
		TransactionTypeID: byte(INCOME),
	}

	return s.createTransaction(userId, transaction)
}

func (s *service) CreateAdjustmentTransaction(userId, accountId uint, amount float64, trType TransactionType) (*TransactionResponse, error) {
	transaction := CreateTransactionDTO{
		Date:              time.Now(),
		Amount:            amount,
		Description:       "Adjusted balance",
		Location:          "",
		ToAccountID:       accountId,
		CategoryID:        category.SystemCategoryID,
		TransactionTypeID: byte(trType),
	}

	return s.createTransaction(userId, transaction)
}

func (s *service) deleteTransaction(userId, id uint) error {
	ok, err := s.repo.transactionExistsAndBelongsToUser(userId, id)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("transaction with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	return s.repo.deleteTransaction(userId, id)
}

func (s *service) getTransactions(userId uint) ([]TransactionResponse, error) {
	return s.repo.getTransactions(userId)
}

func (s *service) getTransaction(userID, txnId uint) (*TransactionResponse, error) {
	ok, err := s.repo.transactionExistsAndBelongsToUser(userID, txnId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("transaction with ID %d does not exist or belong to user with ID %d", txnId, userID)
	}

	client := s.hub.GetClient(userID)
	if client != nil {
		client.Send([]byte("Here is a string...."))
	}

	return s.repo.getTransaction(txnId)
}

func (s *service) GetAccountTransactionsByMonth(accountId uint, year int, month time.Month) ([]TransactionResponse, error) {
	// TODO: move in validation ??? or it is a query param
	if year < 1900 || year > time.Now().Year()+10 {
		return nil, fmt.Errorf("error getting transactions for account: invalid year provided: %d", year)
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("error getting transactions for account: invalid month provided: %d", month)
	}

	return s.repo.getAccountTransactionsByMonth(accountId, year, month)
}

func (s *service) GetTransactionsByFilter(filter TransactionFilter) ([]TransactionResponse, error) {
	return s.repo.findByFilter(filter)
}

func (s *service) updateTransaction(userId, transactionId uint, transaction UpdateTransactionDTO) (*TransactionResponse, error) {
	exists, err := s.repo.transactionExistsAndBelongsToUser(userId, transactionId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("transaction with ID %d does not exist or belong to user with ID %d", transactionId, userId)
	}

	// TODO:  Duplicate in two places, here and create
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, transaction.ToAccountID)
	if err != nil || !ok {
		return nil, fmt.Errorf("toAccountId with id %d doesn't exist or belong to user", transaction.ToAccountID)
	}

	if transaction.TransactionTypeID == byte(TRANSFER) {
		ok, err := s.repo.accountExistsAndBelongsToUser(userId, *transaction.FromAccountID)
		if err != nil || !ok {
			return nil, fmt.Errorf("fromAccountId with id %d doesn't exist or belong to user", *transaction.FromAccountID)
		}
	}

	exists, err = s.repo.categoryExistsAndBelongsToUser(userId, transaction.CategoryID)
	if err != nil || !exists {
		return nil, fmt.Errorf("categoryId with id %d doesn't exist or belong to user", transaction.CategoryID)
	}

	exist, err := s.repo.tagsExistsAndBelongsToUser(userId, transaction.TagIDs)
	if err != nil || !exist {
		return nil, fmt.Errorf("not all tags belong to user or do not exist %v", transaction.TagIDs)
	}

	return s.repo.updateTransaction(transactionId, transaction)
}
