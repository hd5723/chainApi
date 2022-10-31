package model

type TransactionCreateInput struct {
	TxHash          string
	TxIndex         int32
	Height          int32
	ContractAddress string
	Status          int32
	From            string
	To              string
	Value           int64
	GasUsed         int64
	GasLimit        int64
	GasPrice        int64
	TransactionFee  int64
	CreateTime      int32
	UpdateTime      int32
	Timestamp       int32
}
