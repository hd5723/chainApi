package consts

type TxExecuteInfo struct {
	TxHash string
	Err    error
}
type BlockExecuteInfo struct {
	BlockNumber int64
	Err         error
}

type BlockInfo struct {
	BlockNumber int64
}
type TxInfo struct {
	TxHash string
}
