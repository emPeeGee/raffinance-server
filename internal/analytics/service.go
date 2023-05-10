package analytics

import (
	"time"

	"github.com/emPeeGee/raffinance/internal/transaction"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/emPeeGee/raffinance/pkg/util"
)

type Service interface {
	GetCashFlowReport(userID uint, params *RangeDateParams) (*Report, error)
	GetBalanceEvolution(userID uint, params *BalanceEvolutionParams) (*Report, error)
	GetTopTransactions(userID uint, params *TopTransactionsParams) (*Report, error)
	GetCategoriesSpending(userID uint, params *RangeDateParams) (*Report, error)
	GetCategoriesIncome(userID uint, params *RangeDateParams) (*Report, error)
	GetTransactionsCountByDay(userID uint, params *YearlyTransactionsParams) (*Report, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewAnalyticsService(repo Repository, logger log.Logger) *service {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetCashFlowReport(userID uint, params *RangeDateParams) (*Report, error) {
	if params.EndDate != nil && params.StartDate != nil {
		params.EndDate = util.EndOfTheDay(*params.EndDate)
	}

	data, err := s.repo.GetCashFlowReport(userID, params)
	if err != nil {
		return nil, err
	}

	return &Report{
		Title: "Cash flow",
		Data:  data,
	}, nil
}

func (s *service) GetBalanceEvolution(userID uint, params *BalanceEvolutionParams) (*Report, error) {
	if params.EndDate != nil && params.StartDate != nil {
		params.EndDate = util.EndOfTheDay(*params.EndDate)
	}

	data, err := s.repo.GetBalanceEvolutionReport(userID, params)
	if err != nil {
		return nil, err
	}

	return &Report{
		Title: "Balance evolution",
		Data:  data,
	}, nil
}

func (s *service) GetTopTransactions(userID uint, params *TopTransactionsParams) (*Report, error) {
	if params.EndDate != nil && params.StartDate != nil {
		params.EndDate = util.EndOfTheDay(*params.EndDate)
	}

	transactions, err := s.repo.GetTopTransactions(userID, params)
	if err != nil {
		return nil, err
	}

	var topTxns []transaction.TransactionResponse
	for _, t := range transactions {
		topTxns = append(topTxns, transaction.EntityToResponse(&t))
	}

	return &Report{
		Title: "Top transactions",
		Data:  topTxns,
	}, nil
}

func (s *service) GetCategoriesSpending(userID uint, params *RangeDateParams) (*Report, error) {
	if params.EndDate != nil && params.StartDate != nil {
		params.EndDate = util.EndOfTheDay(*params.EndDate)
	}

	data, err := s.repo.GetCategoriesReport(userID, transaction.EXPENSE, params)
	if err != nil {
		return nil, err
	}

	return &Report{
		Title: "Categories Spending",
		Data:  data,
	}, nil
}

func (s *service) GetCategoriesIncome(userID uint, params *RangeDateParams) (*Report, error) {
	if params.EndDate != nil && params.StartDate != nil {
		params.EndDate = util.EndOfTheDay(*params.EndDate)
	}

	data, err := s.repo.GetCategoriesReport(userID, transaction.INCOME, params)
	if err != nil {
		return nil, err
	}

	return &Report{
		Title: "Categories Income",
		Data:  data,
	}, nil
}

func (s *service) GetTransactionsCountByDay(userID uint, params *YearlyTransactionsParams) (*Report, error) {

	// If year is not provided, use the current year
	if params.Year == 0 {
		params.Year = time.Now().Year()
	}

	data, err := s.repo.GetTransactionCountByDay(userID, params)
	if err != nil {
		return nil, err
	}

	return &Report{
		Title: "Count transactions by day",
		Data:  data,
	}, nil
}
