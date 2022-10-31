package entity

// BlockChain is the golang structure for table block_chain.
type BlockChainEntity struct {
	ChainName     string `json:"chain_name"           ` // 区块链名
	BaseCoin      string `json:"base_coin"            ` // e.g BNB
	ChainId       int32  `json:"chain_id"             ` // 链编号
	RpcUrls       string `json:"rpc_urls"             ` // RPC配置
	UpdateTime    int32  `json:"update_time"          ` // 修改时间(时间戳)
	CurrentHeight int32  `json:"current_height"       ` // 当前解析高度
}
