package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type FunctionInfoQueryByTxHashReq struct {
	g.Meta `path:"/queryByTxHash" tags:"functionInfo" method:"get" summary:"QueryByTxHash""`
	TxHash string `v:"required" dc:"交易哈希"`
	//Height       int    ` dc:"区块高"`
	ProtocolCode string `v:"required" dc:"协议编码"`
	ContractCode string `v:"required" dc:"合约编码"`
	Function     string `v:"required" dc:"ABI方法"`
	ChainName    string `v:"required" dc:"公链名称"`
}

type FunctionInfoQueryByHeightReq struct {
	g.Meta `path:"/queryByHeight" tags:"functionInfo" method:"get" summary:"QueryByHeight""`
	//TxHash string ` dc:"交易哈希"`
	BeginHeight  int    ` dc:"开始区块高"`
	EndHeight    int    ` dc:"结束区块高"`
	Address      string ` dc:"用户地址"`
	ProtocolCode string `v:"required" dc:"协议编码"`
	ContractCode string `v:"required" dc:"合约编码"`
	Function     string `v:"required" dc:"ABI方法"`
	ChainName    string `v:"required" dc:"公链名称"`
	PageSize     int    ` dc:"每页请求数量"`
	PageNumber   int    ` dc:"页码数"`
}

type FunctionInfoQueryByBlockTimeReq struct {
	g.Meta         `path:"/queryByBlockTime" tags:"functionInfo" method:"get" summary:"QueryByBlockTime""`
	TxHash         string ` dc:"交易哈希"`
	BlockBeginTime string ` dc:"交易开始时间"`
	BlockEndTime   string ` dc:"交易结束时间"`
	Address        string ` dc:"用户地址"`
	ProtocolCode   string `v:"required" dc:"协议编码"`
	ContractCode   string `v:"required" dc:"合约编码"`
	Function       string `v:"required" dc:"ABI方法"`
	ChainName      string `v:"required" dc:"公链名称"`
	PageSize       int    ` dc:"每页请求数量"`
	PageNumber     int    ` dc:"页码数"`
}

type FunctionInfoQueryAllReq struct {
	g.Meta    `path:"/all" tags:"functionInfo" method:"get" summary:"query all""`
	ChainName string `v:"required" dc:"公链名称"`
}

type FunctionInfoQueryAllTableReq struct {
	g.Meta       `path:"/allTable" tags:"functionInfo" method:"get" summary:"QueryAllTable ""`
	ProtocolCode string `  dc:"合约编码"`
	ContractCode string `  dc:"合约编码"`
	ChainName    string `v:"required" dc:"公链名称"`
	PageSize     int    ` dc:"每页请求数量"`
	PageNumber   int    ` dc:"页码数"`
}
