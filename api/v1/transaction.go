package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TransactionQueryByAddressAndTxHashReq struct {
	g.Meta    `path:"/queryByTxHash" tags:"transaction" method:"get" summary:"query contract by address"`
	TxHash    string `v:"required" dc:"交易哈希"`
	ChainName string `v:"required" dc:"公链名称"`
}

type TransactionQueryByAddressAndHeightReq struct {
	g.Meta          `path:"/queryByAddressAndHeight" tags:"transaction" method:"get" summary:"query contract by code"`
	Height          int    `v:"required" dc:"区块高度"`
	ContractAddress string `v:"required" dc:"所属合约地址"`
	ChainName       string `v:"required" dc:"公链名称"`
	PageSize        int    ` dc:"每页请求数量"`
	PageNumber      int    ` dc:"页码数"`
}

type TransactionQueryByTimeReq struct {
	g.Meta          `path:"/queryByTime" tags:"transaction" method:"get" summary:"query contract by code"`
	BlockBeginTime  string ` dc:"交易开始时间"`
	BlockEndTime    string ` dc:"交易结束时间"`
	ContractAddress string `v:"required" dc:"所属合约地址"`
	ChainName       string `v:"required" dc:"公链名称"`
	PageSize        int    ` dc:"每页请求数量"`
	PageNumber      int    ` dc:"页码数"`
	Address         string ` dc:"用户地址"`
}

type TransactionModifyPriceTypeAndDataReq struct {
	g.Meta    `path:"/modifyPriceType" tags:"transaction" method:"post" summary:"TransactionModifyPriceTypeAndData""`
	ChainName string `v:"required" dc:"公链名称"`
}
