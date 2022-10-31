package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type QueryAbiNameListByContractAddressReq struct {
	g.Meta          `path:"/queryAbiListByContract" tags:"abi" method:"get" summary:"QueryAbiNameListByContractAddressReq""`
	ContractAddress string `v:"required" dc:"合约地址"`
	ChainName       string `v:"required" dc:"公链名称"`
	Type            string `v:"required" dc:"类型"`
	PageSize        int    ` dc:"每页请求数量"`
	PageNumber      int    ` dc:"页码数"`
}
