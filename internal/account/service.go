package account

import (
	"fmt"
	"math"
	"time"

	"github.com/emPeeGee/raffinance/internal/transaction"
	"github.com/emPeeGee/raffinance/pkg/log"
)

type Service interface {
	createAccount(userId uint, account createAccountDTO) (*accountResponse, error)
	deleteAccount(userId, id uint) error
	updateAccount(usedId, accountId uint, account updateAccountDTO) (*accountResponse, error)

	getAccounts(userId uint) ([]accountResponse, error)
	getAccount(userId, accountId uint) (*accountDetailsResponse, error)
	getAccountWithTransactions(userId, id uint) (*accountDetailsResponse, error)
	getAccountTransactionsByMonth(accountId uint, year int, month time.Month) ([]transaction.TransactionResponse, error)
	getAccountBalance(userId, id uint) (float64, error)
	getUserBalance(userId uint) (float64, error)
}

type service struct {
	repo Repository
	// NOTE: I need this service to make transaction on account
	transactionService transaction.Service
	logger             log.Logger
}

func NewAccountService(transactionService transaction.Service, repo Repository, logger log.Logger) *service {
	return &service{
		transactionService: transactionService,
		repo:               repo,
		logger:             logger,
	}
}

func (s *service) createAccount(userId uint, account createAccountDTO) (*accountResponse, error) {
	// check if such name or email already exists, email and name should be unique per user
	exists, err := s.repo.accountExistsAndBelongsToUser(userId, 0, account.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("account with name %s exists", account.Name)
	}

	// First, create account and then if needed the first transaction
	createdAccount, err := s.repo.createAccount(userId, account)
	if err != nil {
		return nil, err
	}

	// In case the account is created with some default balance, create a transaction
	if account.Balance > 0 {
		_, err := s.transactionService.CreateInitialTransaction(userId, createdAccount.ID, account.Balance)
		if err != nil {
			return nil, fmt.Errorf("could not create an initial balance of %f for account %d", account.Balance, createdAccount.ID)
		}
	}

	return createdAccount, nil
}

func (s *service) deleteAccount(userId, id uint) error {
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("account with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	// account with references is prohibited to delete
	if err := s.repo.accountIsUsed(id); err != nil {
		return err
	}

	return s.repo.deleteAccount(userId, id)
}

func (s *service) getAccounts(userId uint) ([]accountResponse, error) {
	return s.repo.getAccounts(userId)
}

func (s *service) updateAccount(userId, accountId uint, account updateAccountDTO) (*accountResponse, error) {
	exists, err := s.repo.accountExistsAndBelongsToUser(userId, accountId, "")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("account with ID %d does not exist or belong to user with ID %d", accountId, userId)
	}

	// check if such name exists, name should be unique per user
	exists, err = s.repo.accountExistsAndBelongsToUser(userId, 0, account.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("account with name %s already exists", account.Name)
	}

	// If the balance is modified, get the current one
	currentAccountBalance, err := s.getAccountBalance(userId, accountId)
	if err != nil {
		return nil, err
	}

	// calculate the difference
	// current 100, modified 200 => 100 - 200 = -100. If adjustment is negative. It means we should make an income with adjusted amount to adjust balance
	// current 200, modified 100 => 200 - 100 = 100. If adjustment is positive. It means we should make an expense with adjusted amount to adjust balance
	adjustedAmount := currentAccountBalance - account.Balance

	if adjustedAmount > 0 {
		_, err := s.transactionService.CreateAdjustmentTransaction(userId, accountId, math.Abs(adjustedAmount), transaction.EXPENSE)
		if err != nil {
			return nil, fmt.Errorf("failed to create an adjustment balance for account %d with err: %s", accountId, err.Error())
		}
	} else if adjustedAmount < 0 {
		_, err := s.transactionService.CreateAdjustmentTransaction(userId, accountId, math.Abs(adjustedAmount), transaction.INCOME)
		if err != nil {
			return nil, fmt.Errorf("failed to create an adjustment balance for account %d with err: %s", accountId, err.Error())
		}
	}

	return s.repo.updateAccount(userId, accountId, account)
}

func (s *service) getAccount(userId, id uint) (*accountDetailsResponse, error) {
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return nil, fmt.Errorf("error checking account ownership: %v", err)
	}

	if !ok {
		return nil, fmt.Errorf("account with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	account, err := s.repo.getAccount(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// Returns account details and transaction from current month
func (s *service) getAccountWithTransactions(userId, id uint) (*accountDetailsResponse, error) {
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return nil, fmt.Errorf("error checking account ownership: %v", err)
	}

	if !ok {
		return nil, fmt.Errorf("account with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	account, err := s.repo.getAccount(id)
	if err != nil {
		return nil, err
	}

	today := time.Now()
	year := today.Year()
	month := today.Month()

	accountTransactionsThisMonth, err := s.transactionService.GetAccountTransactionsByMonth(id, year, month)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions for account: %v", err)
	}

	account.Transactions = accountTransactionsThisMonth
	return account, nil
}

func (s *service) getAccountTransactionsByMonth(accountId uint, year int, month time.Month) ([]transaction.TransactionResponse, error) {
	transactions, err := s.transactionService.GetAccountTransactionsByMonth(accountId, year, month)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *service) getAccountBalance(userId, id uint) (float64, error) {
	ok, err := s.repo.accountExistsAndBelongsToUser(userId, id, "")
	if err != nil {
		return -1, err
	}

	if !ok {
		return -1, fmt.Errorf("account with ID %d does not exist or belong to user with ID %d", id, userId)
	}

	return s.repo.getAccountBalance(id, nil)
}

// TODO: to be moved in user
func (s *service) getUserBalance(userId uint) (float64, error) {
	return s.repo.getUserBalance(userId)
}
