// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ContractSign is the golang structure for table contract_sign.
type ContractSignEntity struct {
	Id           int64  `json:"id"             ` // 主键
	Address      string `json:"address"        ` // 合约地址
	Name         string `json:"name"           ` // 方法名称
	Sign         string `json:"sign"           ` // 签名原始数据
	SignText     string `json:"sign_text"      ` // 签名数据
	SignTextView string `json:"sign_text_view" ` // 签名数据，用于页面展示
	Chain        string `json:"chain"          ` // 所属链
	Type         string `json:"type"           ` // function、event
	UpdateTime   int32  `json:"update_time"    ` // 修改时间(时间戳)
}
