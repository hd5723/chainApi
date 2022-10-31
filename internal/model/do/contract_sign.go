// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ContractSign is the golang structure of table contract_sign for DAO operations like Where/Data.
type ContractSign struct {
	g.Meta       `orm:"table:contract_sign, do:true"`
	Id           interface{} // 主键
	Name         interface{} // 方法名称
	Sign         interface{} // 签名原始数据
	SignText     interface{} // 签名数据
	SignTextView interface{} // 签名数据，用于页面展示
	Chain        interface{} // 所属链
	Address      interface{} // 合约地址
	Type         interface{} // function、event
	CreateTime   interface{} // 创建时间(时间戳)
	UpdateTime   interface{} // 修改时间(时间戳)
}
