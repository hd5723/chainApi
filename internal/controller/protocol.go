package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/model"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"bytes"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	Protocol = cProtocol{}
)

type cProtocol struct{}

func (p *cProtocol) Create(ctx context.Context, req *v1.ProtocolCreateReq) (res *utils.ResponseRes, err error) {
	g.Log().Info(ctx, "ProtocolName:", req.ProtocolName, " ,ProtocolCode", req.ProtocolCode)

	if utils.ExistsSpecialLetters(req.ProtocolCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("Cannot contain special characters"))
		return res, err
	}

	var in model.ProtocolCreateInput
	in.ProtocolCode = req.ProtocolCode
	in.ProtocolName = req.ProtocolName
	in.IsValid = false //初始化设置，需要管理员审核是否通过验证
	in.UpdateTime = int32(gtime.Timestamp())

	chainName := req.ChainName
	rcd, err := service.Protocol().QueryOneByProtocolCode(context.Background(), in.ProtocolCode, chainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	if g.IsEmpty(rcd.ProtocolCode) {
		//新增Protocol数据
		service.Protocol().ExecInsert(context.Background(), chainName, in)
		// 创建dataBase
		var dataBaseSql bytes.Buffer
		dataBaseSql.WriteString(" CREATE database  IF NOT EXISTS  " + in.ProtocolCode)
		_, err := g.DB().Exec(context.Background(), dataBaseSql.String())
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the protocol is exists"))
	}
	return
}

func (p *cProtocol) Audit(ctx context.Context, req *v1.ProtocolAuditReq) (res *utils.ResponseRes, err error) {
	err = service.Protocol().UpdateValid(ctx, req.ProtocolCode, true, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	}
	return
}

// 查询已通过审核的协议列表
func (p *cProtocol) All(ctx context.Context, req *v1.ProtocolReq) (res *utils.ResponseRes, err error) {
	result, err := service.Protocol().QueryList(ctx, req.ProtocolCode, req.IsValid, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(result, len(result)))
	}
	return
}

//查询已通过审核的协议列表
//func (p *cProtocol) ToAuditList(ctx context.Context, req *v1.ProtocolToAuditListReq) (res *utils.ResponseRes, err error) {
//	result, err := service.Protocol().ToAuditList(ctx, req.ChainName)
//	if err != nil {
//		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
//	} else {
//		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(result, len(result)))
//	}
//	return
//}

func (p *cProtocol) Del(ctx context.Context, req *v1.ProtocolDelReq) (res *utils.ResponseRes, err error) {
	err = service.Protocol().DeleteByProtocolCode(ctx, req.ProtocolCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	}
	return
}
