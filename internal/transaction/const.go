package transaction

type TransactionType byte

const (
	INCOME   TransactionType = 1
	EXPENSE  TransactionType = 2
	TRANSFER TransactionType = 3
)
