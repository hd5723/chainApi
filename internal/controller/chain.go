package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/model"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"context"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
)

var (
	Chain = cChain{}
)

type cChain struct{}

func (c *cChain) ChainInstall(ctx context.Context, req *v1.ChainTableInstallReq) (res *utils.ResponseRes, err error) {
	chainSqlPath := gfile.Pwd() + "/sql/chain.sql"
	if !gfile.IsReadable(chainSqlPath) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("sql file can't read!"))
		return
	}

	dataBaseSql := " CREATE database  IF NOT EXISTS  " + req.ChainName
	_, err = g.DB().Exec(context.Background(), dataBaseSql)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	chainSqlContent := strings.ReplaceAll(gfile.GetContents(chainSqlPath), "\n", "")
	v := g.View()
	v.SetDelimiters("${", "}")
	chainSql, err := v.ParseContent(context.Background(), chainSqlContent,
		map[string]interface{}{
			"database": req.ChainName,
		})
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return nil, err
	}

	sqlArray := strings.Split(chainSql, ";")

	for i := 0; i < len(sqlArray); i++ {
		sql := sqlArray[i]
		if g.IsEmpty(sql) {
			continue
		}

		_, err := g.DB().Exec(ctx, sql)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return nil, err
		}
	}
	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	return
}

func (c *cChain) ChainCreate(ctx context.Context, req *v1.ChainCreateReq) (res *utils.ResponseRes, err error) {
	var in model.BlockChainCreateInput
	in.UpdateTime = int32(gtime.Timestamp())
	in.ChainId = int32(req.ChainId)
	in.ChainName = req.ChainName
	in.BaseCoin = req.BaseCoin
	in.RpcUrls = req.RpcUrls
	err = service.Chain().ExecInsert(ctx, req.ChainName, in)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}
	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	return
}

func (c *cChain) ChainUpdate(ctx context.Context, req *v1.ChainUpdateReq) (res *utils.ResponseRes, err error) {
	var in model.BlockChainCreateInput
	in.UpdateTime = int32(gtime.Timestamp())
	in.ChainId = int32(req.ChainId)
	in.ChainName = req.ChainName
	in.BaseCoin = req.BaseCoin
	in.RpcUrls = req.RpcUrls
	err = service.Chain().Update(ctx, req.ChainName, in)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}
	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	return
}

/*
func (c *cChain) ChainQueryRpcNum(ctx context.Context, req *v1.ChainQueryRpcNumReq) (res *utils.ResponseRes, err error) {
	rpcConnectErrNumKey := req.ChainName + consts.RPC_ERR_NUM
	rpcConnectSuccessNumKey := req.ChainName + consts.RPC_SUCCESS_NUM

	errNum, err := g.Redis().Do(ctx, "GET", rpcConnectErrNumKey)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	successNum, err := g.Redis().Do(ctx, "GET", rpcConnectSuccessNumKey)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	var gMap gmap.Map
	gMap.SetIfNotExist("errNum", errNum.Int())
	gMap.SetIfNotExist("successNum", successNum.Int())

	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(gMap))
	return
}*/

func (c *cChain) ChainQueryActiviTask(ctx context.Context, req *v1.ChainQueryActiviTaskReq) (res *utils.ResponseRes, err error) {
	key := req.ChainName + consts.ACTIVITY_TASK_KEY
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}
	if g.IsEmpty(v) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	} else {
		values := strings.Split(v.String(), ",")
		taskSet := gset.NewSet(true)
		for i := 0; i < len(values); i++ {
			value := values[i]
			if !g.IsEmpty(value) && !taskSet.Contains(value) {
				taskSet.Add(value)
			}
		}
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(taskSet, taskSet.Size()))
	}
	return
}

func (c *cChain) ChainQueryById(ctx context.Context, req *v1.QueryChainReq) (res *utils.ResponseRes, err error) {
	rcd := service.Chain().QueryOneByChainId(ctx, int32(req.ChainId), req.ChainName)
	if g.IsEmpty(rcd.ChainName) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		return
	}
	chainHeightKey := req.ChainName + "_curr_height_key"
	v, err := g.Redis().Do(ctx, "GET", chainHeightKey)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}
	rcd.CurrentHeight = v.Int32()
	g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(rcd))
	return
}
