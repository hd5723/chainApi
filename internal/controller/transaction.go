package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/dao"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"OnchainParser/internal/web3/web3Client"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

var (
	Transaction = cTransaction{}
)

type cTransaction struct{}

func (p *cTransaction) QueryOneByTxHash(ctx context.Context, req *v1.TransactionQueryByAddressAndTxHashReq) (res *utils.ResponseRes, err error) {
	result, err := service.Transaction().QueryOneByTxHash(ctx, req.TxHash, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		if g.IsEmpty(result.TxHash) {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		} else {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(result))
		}
	}
	return
}

func (p *cTransaction) QueryByAddressAndBlockNumber(ctx context.Context, req *v1.TransactionQueryByAddressAndHeightReq) (res *utils.ResponseRes, err error) {
	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNmber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNmber = req.PageNumber
	}
	address := common.HexToAddress(req.ContractAddress).Hex()
	result, err := service.Transaction().QueryByAddressAndBlockNumber(ctx, address, req.Height, pageSize, pageNmber, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(result, len(result)))
	}
	return
}

func (p *cTransaction) QueryByTime(ctx context.Context, req *v1.TransactionQueryByTimeReq) (res *utils.ResponseRes, err error) {
	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNmber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNmber = req.PageNumber
	}
	address := common.HexToAddress(req.ContractAddress).Hex()
	result, err := service.Transaction().QueryListByTime(ctx, address, req.BlockBeginTime, req.BlockEndTime, req.Address, pageSize, pageNmber, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(result, len(result)))
	}
	return
}

func (p *cTransaction) TransactionModifyPriceTypeAndData(ctx context.Context, req *v1.TransactionModifyPriceTypeAndDataReq) (res *utils.ResponseRes, err error) {
	//从参数获取操作的 database
	tEnts, err := service.Transaction().QueryLteListByGasPrice(ctx, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	//修复交易表里的gas数据
	modifyTypeSql := "ALTER TABLE " + req.ChainName + ".`transaction`  MODIFY COLUMN gas_used Int64,MODIFY COLUMN gas_limit Int64,MODIFY COLUMN gas_price Int64;"
	_, err = dao.Transaction.DB().Exec(ctx, modifyTypeSql)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	rcd := service.Chain().QueryOneByChainName(ctx, req.ChainName)
	links := strings.Split(rcd.RpcUrls, ",")
	ethClient, _, err := web3Client.GetClientByLinks(ctx, links, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	for i := 0; i < len(tEnts); i++ {
		ent := tEnts[i]
		hexTxHash := common.HexToHash(ent.TxHash)
		//通过RPC，获取 GasPrice
		tx, _, err := ethClient.TransactionByHash(ctx, hexTxHash)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		//修复因类型转换丢失精度的GasPrice数据
		err = dao.Transaction.UpdateGasprice(ctx, ent.TxHash, tx.GasPrice().Int64(), req.ChainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
	}
	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	return
}
