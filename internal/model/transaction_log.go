package model

type TransactionLogCreateInput struct {
	TxHash     string
	BlockHash  string
	Height     int32
	LogIndex   int32
	Removed    bool
	Address    string
	Data       string
	Type       string
	Topics     []string
	UpdateTime int32
}
