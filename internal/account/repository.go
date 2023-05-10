package account

import (
	"errors"
	"fmt"
	"time"

	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/internal/transaction"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/emPeeGee/raffinance/pkg/util"

	"gorm.io/gorm"
)

var ErrAccountBalanceNotFound = errors.New("account balance not found")

type Repository interface {
	getAccounts(userId uint) ([]accountResponse, error)
	getAccount(accountId uint) (*accountDetailsResponse, error)
	createAccount(userId uint, Account createAccountDTO) (*accountResponse, error)
	updateAccount(userId, accountId uint, account updateAccountDTO) (*accountResponse, error)
	deleteAccount(userId, id uint) error
	accountExistsAndBelongsToUser(userID, id uint, name string) (bool, error)
	accountIsUsed(accountId uint) error
	getAccountBalance(id uint, month *time.Time) (float64, error)
	getUserBalance(userID uint) (float64, error)
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewAccountRepository(db *gorm.DB, logger log.Logger) *repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) createAccount(userId uint, account createAccountDTO) (*accountResponse, error) {
	newAccount := entity.Account{
		Name:     account.Name,
		Currency: account.Currency,
		Color:    account.Color,
		Icon:     account.Icon,
		UserID:   &userId,
	}

	if err := r.db.Create(&newAccount).Error; err != nil {
		return nil, err
	}

	r.logger.Info("new account, ", util.StringifyAny(newAccount))
	createdAccount := &accountResponse{
		ID:        newAccount.ID,
		Name:      newAccount.Name,
		Currency:  newAccount.Currency,
		Color:     newAccount.Color,
		Icon:      newAccount.Icon,
		Balance:   account.Balance,
		CreatedAt: newAccount.CreatedAt,
		UpdatedAt: newAccount.UpdatedAt,
	}

	return createdAccount, nil
}

func (r *repository) updateAccount(userId, accountId uint, account updateAccountDTO) (*accountResponse, error) {
	// NOTE: When update with struct, GORM will only update non-zero fields, you might want to use
	// map to update attributes or use Select to specify fields to update
	if err := r.db.Model(&entity.Account{}).Where("id = ?", accountId).Updates(map[string]interface{}{
		"name":     account.Name,
		"currency": account.Currency,
		"icon":     account.Icon,
		"color":    account.Color,
	}).Error; err != nil {
		return nil, err
	}

	var updatedAccount accountResponse
	if err := r.db.Model(&entity.Account{}).First(&updatedAccount, accountId).Error; err != nil {
		return nil, err
	}

	// Calculating the balance dynamically
	accountBalance, err := r.getAccountBalance(accountId, nil)
	if err != nil {
		return nil, err
	}

	updatedAccount.Balance = accountBalance

	return &updatedAccount, nil
}

func (r *repository) deleteAccount(userId, id uint) error {
	return r.db.Delete(&entity.Account{}, id).Error
}

func (r *repository) getAccounts(userId uint) ([]accountResponse, error) {
	var accounts []accountResponse
	var accountsR []accountResponse

	query := `
		SELECT ac.id, ac.created_at, ac.updated_at, ac.name, ac.color, ac.currency, ac.icon,
      (SELECT COUNT(DISTINCT id)
				FROM transactions AS t
				WHERE t.deleted_at IS NULL AND (t.from_account_id = ac.id OR t.to_account_id = ac.id)) AS transaction_count
    FROM accounts as ac
    WHERE ac.user_id = ? AND ac.deleted_at IS NULL;
	`

	if err := r.db.Raw(query, userId).Scan(&accounts).Error; err != nil {
		return nil, err
	}

	for _, account := range accounts {
		// TODO: get the account balance dynamically. Is there a better way for it?
		accountBalance, err := r.getAccountBalance(account.ID, nil)
		if err != nil {
			return nil, err
		}

		now := new(time.Time)
		*now = time.Now()
		lastMonth := now.AddDate(0, -1, 0)

		accountBalanceThisMonth, err := r.getAccountBalance(account.ID, now)
		if err != nil {
			return nil, err
		}

		accountBalanceLastMonth, err := r.getAccountBalance(account.ID, &lastMonth)
		if err != nil {
			return nil, err
		}

		var rate float64
		if accountBalanceLastMonth == 0 {
			rate = 0 // avoid division by zero
		} else {
			rate = ((accountBalanceThisMonth - accountBalanceLastMonth) / accountBalanceLastMonth) * 100
		}

		diff := accountBalanceLastMonth - accountBalanceThisMonth
		r.logger.Debugf("%f %f and DIFF %f", accountBalanceThisMonth, accountBalanceLastMonth, diff)

		accountsR = append(accountsR, accountResponse{
			ID:                account.ID,
			Name:              account.Name,
			Currency:          account.Currency,
			Balance:           accountBalance,
			Color:             account.Color,
			Icon:              account.Icon,
			CreatedAt:         account.CreatedAt,
			UpdatedAt:         account.UpdatedAt,
			TransactionCount:  account.TransactionCount,
			RateWithPrevMonth: &rate,
		})
	}

	return accountsR, nil
}

func (r *repository) accountExistsAndBelongsToUser(userID, id uint, name string) (bool, error) {
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

	if err := r.db.Model(&entity.Account{}).
		Where(whereClause, values...).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// TODO: Can this be achieved using constraints ???
func (r *repository) accountIsUsed(accountId uint) error {
	var count int64

	r.db.Model(&entity.Transaction{}).
		Where("from_account_id = ? OR to_account_id = ?", accountId, accountId).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("cannot delete account with ID %d because it is used by %d transactions", accountId, count)
	}

	return nil
}

func (r *repository) getAccountBalance(id uint, month *time.Time) (float64, error) {
	// Calculate the total balance of this account for the given month
	var nonTransferBalance, transferBalance float64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Calculate the non-transfer total balance
		txQuery := tx.Table("transactions").
			Where("deleted_at is null and to_account_id = ? and transaction_type_id <> ?", id, transaction.TRANSFER).
			Select("COALESCE(SUM(CASE WHEN transaction_type_id = ? THEN amount ELSE -amount END), 0)", transaction.INCOME)

		if month != nil {
			txQuery = txQuery.Where("to_char(date, 'YYYY-MM') = ?", month.Format("2006-01"))
		}

		err := txQuery.Row().Scan(&nonTransferBalance)
		if err != nil {
			return err
		}

		// Calculate the transfer total balance
		txQuery = tx.Table("transactions").
			Where("deleted_at is null and (to_account_id = ? or from_account_id = ?) and transaction_type_id = ?", id, id, transaction.TRANSFER).
			Select("COALESCE(SUM(CASE WHEN to_account_id = ? THEN amount ELSE -amount END), 0)", id)

		if month != nil {
			txQuery = txQuery.Where("to_char(date, 'YYYY-MM') = ?", month.Format("2006-01"))
		}

		err = txQuery.Row().Scan(&transferBalance)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrAccountBalanceNotFound
		}

		return 0, err
	}

	r.logger.Infof("Account %d balance: balance excluding transfers %f, transfer total %f", id, nonTransferBalance, transferBalance)

	return nonTransferBalance + transferBalance, nil
}

func (r *repository) getUserBalance(userID uint) (float64, error) {
	var totalBalance float64

	// Retrieve all accounts associated with the user
	var accounts []entity.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return 0, err
	}

	// Iterate over each account to calculate the total balance
	for _, account := range accounts {
		accountBalance, err := r.getAccountBalance(account.ID, nil)
		if err != nil {
			return 0, err
		}

		totalBalance += accountBalance
	}

	r.logger.Infof("User %d balance: %f", userID, totalBalance)

	return totalBalance, nil
}

func (r *repository) getAccount(accountId uint) (*accountDetailsResponse, error) {
	var account *accountDetailsResponse

	query := `
		SELECT ac.id, ac.created_at, ac.updated_at, ac.name, ac.color, ac.currency, ac.icon,
      (SELECT COUNT(DISTINCT t.id)
				FROM transactions AS t
				WHERE t.deleted_at IS NULL AND (t.from_account_id = ac.id OR t.to_account_id = ac.id)) AS transaction_count
    FROM accounts as ac
    WHERE ac.id = ? AND ac.deleted_at IS NULL;
		`

	if err := r.db.Raw(query, accountId).Scan(&account).Error; err != nil {
		return nil, err
	}

	accountBalance, err := r.getAccountBalance(accountId, nil)
	if err != nil {
		return nil, err
	}

	account.Balance = accountBalance

	return account, nil
}
