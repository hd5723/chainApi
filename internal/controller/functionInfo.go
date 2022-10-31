package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	_ "math/big"
	"time"
)

var (
	FunctionInfo = cFunctionInfo{}
)

type cFunctionInfo struct{}

func (p *cFunctionInfo) QueryAllTable(ctx context.Context, req *v1.FunctionInfoQueryAllTableReq) (res *utils.ResponseRes, err error) {

	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNmber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNmber = req.PageNumber
	}
	limitBegin := pageSize * (pageNmber - 1)

	tableWithDatabase := "information_schema.tables"

	var temp gdb.Result
	if !g.IsEmpty(req.ProtocolCode) && !g.IsEmpty(req.ContractCode) {
		contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		if g.IsEmpty(contract.ContractCode) {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
			return nil, err
		}
		table_name := contract.ProtocolCode + "_" + contract.ContractCode + "_call_%"
		temp, err = g.DB().Model(tableWithDatabase).Safe().Ctx(ctx).Fields("table_schema", "table_name").Where("table_schema=?", contract.ProtocolCode).WhereLike("table_name", table_name).Limit(limitBegin, pageSize).All()
	} else if !g.IsEmpty(req.ProtocolCode) && g.IsEmpty(req.ContractCode) {
		temp, err = g.DB().Model(tableWithDatabase).Safe().Ctx(ctx).Fields("table_schema", "table_name").Where("table_schema=?", req.ProtocolCode).
			WhereLike("table_name", req.ProtocolCode+"_%").
			WhereLike("table_name", "%_call_%").
			Limit(limitBegin, pageSize).All()
	} else if g.IsEmpty(req.ProtocolCode) && !g.IsEmpty(req.ContractCode) {
		contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		if g.IsEmpty(contract.ContractCode) {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
			return nil, err
		}
		table_name := contract.ProtocolCode + "_" + contract.ContractCode + "_call_%"
		temp, err = g.DB().Model(tableWithDatabase).Safe().Ctx(ctx).Fields("table_schema", "table_name").WhereLike("table_name", table_name).Limit(limitBegin, pageSize).All()
	} else if g.IsEmpty(req.ProtocolCode) && g.IsEmpty(req.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("The protocolCode and contractCode both empty!"))
		return res, err
	}

	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		resultList := utils.AdaptBigIntResult(temp)
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(resultList, len(resultList)))
	}
	return
}

func (p *cFunctionInfo) QueryByTxHash(ctx context.Context, req *v1.FunctionInfoQueryByTxHashReq) (res *utils.ResponseRes, err error) {
	contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	if g.IsEmpty(contract.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
		return nil, err
	}

	contractSign, err := service.ContractSign().QueryOneByName(ctx, contract.ContractAddress, req.Function, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	if g.IsEmpty(contractSign.Address) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function is not exists"))
		return nil, err
	}

	if contractSign.Type != "function" {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function name is not exists"))
		return
	}

	tableName := contract.ProtocolCode + "_" + contract.ContractCode + "_" + service.ContractSign().AdaptFunctionType(contractSign.Type) + "_" + req.Function
	tableWithDatabase := contract.ProtocolCode + "." + tableName
	temp, err := g.DB().Model(tableWithDatabase).Safe().Ctx(ctx).Where("call_tx_hash =?", req.TxHash).One()
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		result := utils.AdaptBigIntRecord(temp)
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(result))
	}
	return
}

func (p *cFunctionInfo) QueryByHeight(ctx context.Context, req *v1.FunctionInfoQueryByHeightReq) (res *utils.ResponseRes, err error) {
	contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return nil, err
	}

	if g.IsEmpty(contract.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
		return nil, err
	}

	contractSign, err := service.ContractSign().QueryOneByName(ctx, contract.ContractAddress, req.Function, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	if g.IsEmpty(contractSign.Address) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function is not exists"))
		return nil, err
	}

	if contractSign.Type != "function" {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function name is not exists"))
		return
	}

	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNumber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNumber = req.PageNumber
	}
	limitBegin := pageSize * (pageNumber - 1)

	tableName := contract.ProtocolCode + "_" + contract.ContractCode + "_" + service.ContractSign().AdaptFunctionType(contractSign.Type) + "_" + req.Function
	tableWithDatabase := contract.ProtocolCode + "." + tableName
	var dataModel *gdb.Model
	dataModel = g.DB().Model(tableWithDatabase).Safe().Ctx(ctx)
	wBuild := dataModel.Builder()
	if !g.IsEmpty(req.Address) {
		address := common.HexToAddress(req.Address).Hex()
		wBuild = wBuild.Where("call_from =?", address)
	}
	if !g.IsEmpty(req.BeginHeight) {
		wBuild = wBuild.WhereGTE("call_height", req.BeginHeight)
	}
	if !g.IsEmpty(req.EndHeight) {
		wBuild = wBuild.WhereLTE("call_height", req.EndHeight)
	}
	wBuild = wBuild.Where("1=1") // 上述条件都不满足，where 需要一个默认条件
	temp, err := dataModel.Where(wBuild).OrderDesc("call_block_time").Limit(limitBegin, pageSize).All()

	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		resultList := utils.AdaptBigIntResult(temp)
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(resultList, len(resultList)))
	}
	return
}

func (p *cFunctionInfo) QueryByBlockTime(ctx context.Context, req *v1.FunctionInfoQueryByBlockTimeReq) (res *utils.ResponseRes, err error) {
	curTime := time.Now()
	contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	g.Log().Info(ctx, "QueryByBlockTime.QueryByContractCode doing time:", time.Now().Sub(curTime))
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return nil, err
	}

	if g.IsEmpty(contract.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
		return nil, err
	}

	contractSign, err := service.ContractSign().QueryOneByName(ctx, contract.ContractAddress, req.Function, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return nil, err
	}

	if g.IsEmpty(contractSign.Address) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function is not exists"))
		return nil, err
	}

	if contractSign.Type != "function" {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the function name is not exists"))
		return
	}

	tableName := contract.ProtocolCode + "_" + contract.ContractCode + "_" + service.ContractSign().AdaptFunctionType(contractSign.Type) + "_" + req.Function
	tableWithDatabase := contract.ProtocolCode + "." + tableName

	var pageSize = 50
	if !g.IsEmpty(req.PageSize) {
		pageSize = req.PageSize
	}

	var pageNumber = 1
	if !g.IsEmpty(req.PageNumber) {
		pageNumber = req.PageNumber
	}
	limitBegin := pageSize * (pageNumber - 1)

	var dataModel *gdb.Model
	dataModel = g.DB().Model(tableWithDatabase).Safe().Ctx(ctx)
	wBuild := dataModel.Builder()
	if !g.IsEmpty(req.Address) {
		address := common.HexToAddress(req.Address).Hex()
		wBuild = wBuild.Where("call_from =?", address)
	}
	if !g.IsEmpty(req.BlockBeginTime) {
		wBuild = wBuild.WhereGTE("call_block_time", req.BlockBeginTime)
	}
	if !g.IsEmpty(req.BlockEndTime) {
		wBuild = wBuild.WhereLTE("call_block_time", req.BlockEndTime)
	}
	wBuild = wBuild.Where("1=1") // 上述条件都不满足，where 需要一个默认条件
	temp, err := dataModel.Where(wBuild).OrderDesc("call_block_time").Limit(limitBegin, pageSize).All()
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		resultList := utils.AdaptBigIntResult(temp)
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(resultList, len(resultList)))
	}
	defer g.Log().Info(ctx, "QueryByBlockTime doing time:", time.Now().Sub(curTime), " , fun:",
		req.Function, " , protocol:", req.ProtocolCode, " , contract:", req.ContractCode)
	return
}
