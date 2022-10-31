package entity

// transactionLog is the golang structure for table transaction_log.
type TransactionLogEntity struct {
	TxHash     string `json:"tx_hash"           `   // 当前交易的哈希值
	BlockHash  string `json:"block_hash"          ` // 当前区块的哈希值
	LogIndex   int32  `json:"log_index"           ` // index
	Height     int32  `json:"height"            `   // 区块号
	Removed    bool   `json:"removed"           `   // 是否已移除
	Address    string `json:"address"  `            // 地址
	Data       string `json:"data"             `    // log data
	Type       string `json:"type"         `        // type
	Topics     string `json:"topics"         `      // topics，Array<String>类型
	UpdateTime int32  `json:"update_time"       `   // 修改时间(时间戳)
}
