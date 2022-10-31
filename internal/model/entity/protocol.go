package entity

// Protocol is the golang structure for table protocol.
type ProtocolEntity struct {
	ProtocolCode string `json:"protocol_code"           `  // 编码，如pancakeswap
	ProtocolName string `json:"protocol_name"            ` // 名称，如PancakeSwap
	IsValid      bool   `json:"is_valid"             `     // 是否有效, true:有效；false:无效
	UpdateTime   int32  `json:"update_time"          `     // 修改时间(时间戳)
}
