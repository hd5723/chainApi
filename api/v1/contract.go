package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ContractReq struct {
	g.Meta       `path:"/all" tags:"contract" method:"get" summary:"protocol page""`
	ChainName    string `v:"required" dc:"公链名称"`
	ProtocolCode string `json:"protocolCode"   dc:"协议编码"`
	ContractCode string `json:"contractCode"  dc:"合约编码"`
	IsValid      string `json:"isValid"   dc:"审核状态"` //审核:1 or 未审核:0 , all:
}

type ContractCreateReq struct {
	g.Meta          `path:"/create" tags:"contract" method:"post" summary:"contract create page""`
	ProtocolCode    string `v:"required"  dc:"协议编码"`
	ContractCode    string `v:"required"  dc:"合约编码"`
	ContractAddress string `v:"required"  dc:"合约地址"`
	AbiJson         string `v:"required"  dc:"AbiJson"`
	DeployHeight    int    `v:"required"  dc:"合约部署区块高度"`
	ChainName       string `v:"required"  dc:"公链名称"`
}

type ContractAuditReq struct {
	g.Meta       `path:"/audit" tags:"contract" method:"post" summary:"contract audit page""`
	ProtocolCode string `v:"required"  dc:"协议编码"`
	ContractCode string `v:"required"  dc:"合约编码"`
	OnceHeight   int    `v:"required"  dc:"单次运行爬取的区块"`
	DataType     int    `v:"required"  dc:"爬取区块方式，0:通过扫描区块筛选数据；1:通过查询事件过滤数据"`
	ChainName    string `v:"required"  dc:"公链名称"`
}

type ContractUpdateInfoReq struct {
	g.Meta       `path:"/update" tags:"contract" method:"post" summary:"update contract page""`
	ProtocolCode string `v:"required"  dc:"协议编码"`
	ContractCode string `v:"required"  dc:"合约编码"`
	OnceHeight   string `   dc:"单次运行爬取的区块"`
	DeployHeight string `   dc:"合约部署区块高度"`
	ChainName    string `v:"required"  dc:"公链名称"`
}

//type ContractToAuditListReq struct {
//	g.Meta    `path:"/toAuditList" tags:"contract" method:"get" summary:"get contract toAuditList data""`
//	ChainName string `v:"required" dc:"公链名称"`
//}

type ContractDelReq struct {
	g.Meta       `path:"/del" tags:"contract" method:"delete" summary:"del contract""`
	ContractCode string `v:"required" dc:"合约编码"`
	ProtocolCode string `v:"required"  dc:"协议编码"`
	ChainName    string `v:"required" dc:"公链名称"`
}

type ContractTaskReq struct {
	g.Meta       `path:"/task" tags:"contract" method:"post" summary:"recover contract""`
	ContractCode string `v:"required" dc:"合约编码"`
	ProtocolCode string `v:"required"  dc:"协议编码"`
	ChainName    string `v:"required" dc:"公链名称"`
	Type         int    `v:"required" dc:"操作类型，1:pause, 2:recover"`
}

type ContractQueryByAddressReq struct {
	g.Meta          `path:"/queryByAddress" tags:"contract" method:"get" summary:"query contract by address""`
	ContractAddress string `v:"required" dc:"合约地址"`
	ProtocolCode    string `v:"required"  dc:"协议编码"`
	ChainName       string `v:"required" dc:"公链名称"`
}

type ContractQueryByCodeReq struct {
	g.Meta       `path:"/queryByCode" tags:"contract" method:"get" summary:"query contract by code""`
	ContractCode string `v:"required" dc:"合约编码"`
	ProtocolCode string `v:"required"  dc:"协议编码"`
	ChainName    string `v:"required" dc:"公链名称"`
}
