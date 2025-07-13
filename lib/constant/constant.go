package constant

type InputState int

const (
	ModeMenu InputState = iota
	ModeListWatchlist
	ModeListTransaction
	ModeCreateTransaction
)

const (
	TransactionTypeBuy  = "BUY"
	TransactionTypeSell = "SELL"
)
