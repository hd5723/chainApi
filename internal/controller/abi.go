package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	Abi = cAbi{}
)

type cAbi struct{}

func (p *cAbi) QueryAbiNameListByContractAddress(ctx context.Context, req *v1.QueryAbiNameListByContractAddressReq) (res *utils.ResponseRes, err error) {

	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNmber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNmber = req.PageNumber
	}
	limitBegin := pageSize * (pageNmber - 1)
	address := common.HexToAddress(req.ContractAddress).Hex()
	contractSign, err := service.ContractSign().QueryListByNameAndType(ctx, address, limitBegin, pageSize, req.Type, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(contractSign, contractSign.Size()))
	}
	return
}
