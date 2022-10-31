package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ProtocolReq struct {
	g.Meta       `path:"/all" tags:"protocol" method:"get" summary:"get protocol all data"`
	ChainName    string `v:"required" dc:"公链名称"`
	ProtocolCode string `json:"protocolCode"   dc:"协议编码"`
	IsValid      string `json:"isValid"   dc:"审核状态"` //审核:1 or 未审核:0 , all:
}

type ProtocolCreateReq struct {
	g.Meta       `path:"/create" tags:"protocol" method:"post" summary:"protocol create page"`
	ProtocolCode string `json:"protocolCode" v:"required"  dc:"协议编码"`
	ProtocolName string `json:"protocolName" v:"required"  dc:"协议名称"`
	ChainName    string `json:"chainName" v:"required"  dc:"公链名称"`
}

type ProtocolAuditReq struct {
	g.Meta       `path:"/audit" tags:"protocol" method:"post" summary:"protocol audit page"`
	ProtocolCode string `json:"protocolCode" v:"required"  dc:"协议编码"`
	ChainName    string `json:"chainName" v:"required"  dc:"公链名称"`
}

//type ProtocolToAuditListReq struct {
//	g.Meta    `path:"/toAuditList" tags:"protocol" method:"get" summary:"get protocol toAuditList data"`
//	ChainName string `v:"required" dc:"公链名称"`
//}

type ProtocolDelReq struct {
	g.Meta       `path:"/del" tags:"protocol" method:"delete" summary:"del protocol"`
	ProtocolCode string `v:"required" dc:"协议编码"`
	ChainName    string `v:"required" dc:"公链名称"`
}
