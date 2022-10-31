package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ChainTableInstallReq struct {
	g.Meta    `path:"/install" tags:"chain" method:"post" summary:"chainInstall ""`
	ChainName string `v:"required" dc:"公链名称"`
}

type ChainCreateReq struct {
	g.Meta    `path:"/dataCreate" tags:"chain" method:"post" summary:"chainDataCreate""`
	ChainName string `v:"required" dc:"公链名称"`
	ChainId   int    `v:"required" dc:"公链id"`
	BaseCoin  string `v:"required" dc:"coin"`
	RpcUrls   string `v:"required" dc:"Rpc urls"`
}

type ChainUpdateReq struct {
	g.Meta    `path:"/dataUpdate" tags:"chain" method:"post" summary:"ChainUpdateReq""`
	ChainName string `v:"required" dc:"公链名称"`
	ChainId   int    `v:"required" dc:"公链id"`
	BaseCoin  string `  dc:"coin"`
	RpcUrls   string `  dc:"Rpc urls"`
}

type ChainQueryRpcNumReq struct {
	g.Meta    `path:"/queryRpcNum" tags:"chain" method:"get" summary:"queryRpcNum""`
	ChainName string ` dc:"公链名称"`
}

type QueryChainReq struct {
	g.Meta    `path:"/queryById" tags:"chain" method:"get" summary:"queryById""`
	ChainName string `v:"required" dc:"公链名称"`
	ChainId   int    `v:"required" dc:"公链id"`
}

type ChainQueryActiviTaskReq struct {
	g.Meta    `path:"/queryActiviTask" tags:"chain" method:"get" summary:"queryActiviTask""`
	ChainName string ` dc:"公链名称"`
}
