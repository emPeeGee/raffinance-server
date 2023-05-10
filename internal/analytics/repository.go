package analytics

import (
	"time"

	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/internal/transaction"
	"github.com/emPeeGee/raffinance/pkg/log"
	"gorm.io/gorm"
)

type Repository interface {
	GetCashFlowReport(userID uint, params *RangeDateParams) ([]CashFlowReport, error)
	GetBalanceEvolutionReport(userID uint, params *BalanceEvolutionParams) ([]DateValue, error)

	GetTopTransactions(userID uint, params *TopTransactionsParams) ([]entity.Transaction, error)
	GetCategoriesReport(userID uint, txnType transaction.TransactionType, params *RangeDateParams) ([]LabelValue, error)
	GetTransactionCountByDay(userID uint, params *YearlyTransactionsParams) ([]DateValue, error)
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewAnalyticsRepository(db *gorm.DB, logger log.Logger) *repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) GetCashFlowReport(userID uint, params *RangeDateParams) ([]CashFlowReport, error) {
	var trends []CashFlowReport

	query := r.db.Table("transactions").
		// Select("date::date AS date, SUM(CASE WHEN transactions.transaction_type_id = 1 THEN transactions.amount WHEN transactions.transaction_type_id = 2 THEN -transactions.amount ELSE 0  END) AS value").
		Select(`date::date AS date, 
		 SUM(CASE WHEN transactions.transaction_type_id = 1 THEN transactions.amount ELSE 0 END) AS income, 
		 SUM(CASE WHEN transactions.transaction_type_id = 2 THEN transactions.amount ELSE 0 END) AS expense,
		 SUM(CASE WHEN transactions.transaction_type_id = 1 THEN transactions.amount WHEN transactions.transaction_type_id = 2 THEN -transactions.amount ELSE 0 END) AS cash_flow`).
		Joins("JOIN accounts ON transactions.to_account_id = accounts.id").
		Where("transactions.deleted_at IS NULL AND accounts.user_id = ?", userID).
		Group("date::date").
		Order("date::date ASC")

	if params.StartDate != nil && params.EndDate != nil {
		query.Where("date BETWEEN ? AND ?", params.StartDate, params.EndDate)
	}

	if err := query.Scan(&trends).Error; err != nil {
		return nil, err
	}

	return trends, nil
}

func (r *repository) GetBalanceEvolutionReport(userID uint, params *BalanceEvolutionParams) ([]DateValue, error) {
	var reports []DateValue

	query := r.db.Table("transactions").
		Select("date::date AS date, SUM(CASE WHEN transactions.transaction_type_id = 1 THEN transactions.amount WHEN transactions.transaction_type_id = 2 THEN -transactions.amount ELSE 0  END) AS value").
		Joins("JOIN accounts ON transactions.to_account_id = accounts.id").
		Where("transactions.deleted_at IS NULL AND accounts.user_id = ?", userID).
		Group("date::date").
		Order("date::date ASC")

	if params.StartDate != nil {
		query = query.Where("date >= ?", params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("date <= ?", params.EndDate)
	}

	if params.AccountID != nil {
		// If there are account id, when it should be calculated different:
		// transfers should be taken into account
		query = query.Where("transactions.to_account_id = ?", params.AccountID)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balance float64 = 0
	for rows.Next() {
		var report DateValue
		var date time.Time
		var amount float64
		if err := rows.Scan(&date, &amount); err != nil {
			return nil, err
		}
		balance += amount
		report.Date = date
		report.Value = balance
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *repository) GetTopTransactions(userID uint, params *TopTransactionsParams) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	query := r.db.Model(&entity.Transaction{}).
		Preload("Tags").
		Preload("Category").
		Joins("JOIN accounts ON transactions.to_account_id = accounts.id").
		Where("transactions.deleted_at IS NULL AND accounts.user_id = ?", userID).
		Order("amount DESC").
		Limit(int(params.Limit))

	if params.StartDate != nil && params.EndDate != nil {
		query.Where("date BETWEEN ? AND ?", params.StartDate, params.EndDate)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *repository) GetCategoriesReport(userID uint, txnType transaction.TransactionType, params *RangeDateParams) ([]LabelValue, error) {
	var byCategory []LabelValue

	// Query the category-wise spending data for the user within the specified date range
	query := r.db.Table("transactions").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Joins("JOIN accounts ON transactions.from_account_id = accounts.id OR transactions.to_account_id = accounts.id").
		Select("categories.name AS label, SUM(transactions.amount) AS value").
		Where("accounts.user_id = ? AND transactions.transaction_type_id = ?", userID, txnType).
		Where("transactions.deleted_at IS NULL AND categories.deleted_at IS NULL").
		Group("categories.name")

	if params.StartDate != nil && params.EndDate != nil {
		query.Where("date BETWEEN ? AND ?", params.StartDate, params.EndDate)
	}

	err := query.Scan(&byCategory).Error
	if err != nil {
		return nil, err
	}

	return byCategory, nil
}

func (r *repository) GetTransactionCountByDay(userID uint, params *YearlyTransactionsParams) ([]DateValue, error) {
	startDate := time.Date(params.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(1, 0, 0).Add(-time.Nanosecond)

	var counts []DateValue

	query := r.db.Table("transactions").
		Select("date::date AS date, COUNT(*) AS value").
		Joins("JOIN accounts ON transactions.to_account_id = accounts.id").
		Where("transactions.deleted_at IS NULL AND accounts.user_id = ?", userID).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("date::date").
		Order("date::date ASC")

	if err := query.Scan(&counts).Error; err != nil {
		return nil, err
	}

	return counts, nil
}
