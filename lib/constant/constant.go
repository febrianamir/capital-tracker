package constant

type InputState int

const (
	ModeMenu InputState = iota
	ModeListTransaction
	ModeCreateTransaction
)

const (
	TransactionTypeBuy  = "BUY"
	TransactionTypeSell = "SELL"
)
