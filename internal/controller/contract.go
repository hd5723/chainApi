package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/chain/ethEvm/chain/baseConfig"
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
)

var (
	Contract = cContract{}
)

type cContract struct{}

func (c *cContract) Create(ctx context.Context, req *v1.ContractCreateReq) (res *utils.ResponseRes, err error) {

	if utils.ExistsSpecialLetters(req.ProtocolCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("Cannot contain special characters"))
		return res, err
	}

	if utils.ExistsSpecialLetters(req.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("Cannot contain special characters"))
		return res, err
	}

	var in model.ContractCreateInput
	in.ProtocolCode = req.ProtocolCode
	in.ContractCode = req.ContractCode
	in.ContractAddress = common.HexToAddress(req.ContractAddress).Hex()
	in.AbiJson = req.AbiJson
	in.DeployHeight = int32(req.DeployHeight)
	in.IsValid = false //初始化设置，需要管理员审核是否通过验证
	in.CurrHeight = 0
	in.UpdateTime = int32(gtime.Timestamp())
	//in.OnceHeight = 10 //初期人工审核设置，每次任务爬取的区块数量
	//in.DataType = 0    //初期人工审核设置，0：通过区块爬取；1：通过查询区块区间爬取

	chainName := req.ChainName

	rcd, err := service.Contract().QueryByContractCode(ctx, in.ProtocolCode, in.ContractCode, chainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	if g.IsEmpty(rcd.ContractAddress) {
		temp, err := service.Contract().QueryOneByContractAddress(ctx, req.ContractAddress, chainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		if !g.IsEmpty(temp.ContractAddress) {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract address is exists"))
			return res, err
		}

		contractAbi, err := abi.JSON(strings.NewReader(in.AbiJson))
		if err != nil {
			g.Log().Error(context.Background(), "NewReader:", err)
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}

		//解析event
		events, eventSigs, err := service.Contract().DoEventContractSign(ctx, in, contractAbi.Events, chainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		in.ListenerEvent = events
		in.SignedListenerEvent = eventSigs

		//解析function
		functions, functionSigs, err := service.Contract().DoFunctionContractSign(ctx, in, contractAbi.Methods, chainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return res, err
		}
		in.ListenerFunction = functions
		in.SignedListenerFunction = functionSigs

		err = service.Contract().Create(context.Background(), in, chainName)
		if err == nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		} else {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		}
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is exists"))
	}
	return
}

func (p *cContract) queryByAddress(ctx context.Context, req *v1.ContractQueryByAddressReq) (res *utils.ResponseRes, err error) {
	address := common.HexToAddress(req.ContractAddress).Hex()
	result, err := service.Contract().QueryByContractAddress(ctx, req.ProtocolCode, address, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(result))
	}
	return
}

func (p *cContract) queryByCode(ctx context.Context, req *v1.ContractQueryByCodeReq) (res *utils.ResponseRes, err error) {
	result, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithData(result))
	}
	return
}

func (p *cContract) Audit(ctx context.Context, req *v1.ContractAuditReq) (res *utils.ResponseRes, err error) {

	ent, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	protocol, err := service.Protocol().QueryOneByProtocolCode(ctx, req.ProtocolCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	if g.IsEmpty(protocol.ProtocolCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the protocol is not exists"))
		return nil, err
	}

	if !protocol.IsValid {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the protocol is not approved"))
		return
	}

	if g.IsEmpty(ent.ContractCode) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
		return nil, err
	}

	if ent.IsValid {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract had audited"))
		return
	}
	//设置每次处理区块数量
	pnum, err := g.Cfg().Get(ctx, "chain.onceLength")
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return res, err
	}

	// 审核通过时
	// runHeight = chainNamde + _curr_height_key 缓存里的值
	// 如果缓存里没有（程序第一次部署），从Rpc获取最新节点（这里不需要once_height)
	var runHeight int
	chainHeightKey := req.ChainName + "_curr_height_key"
	v, err := g.Redis().Do(ctx, "GET", chainHeightKey)
	if err != nil || v.Val() == nil {
		//在存入时直接以blockchain的current_height + once_height 的值存入
		runHeight = int(baseConfig.CHAIN_LAST_HEIGHT) + pnum.Int()
	} else {
		runHeight = v.Int() + pnum.Int()
	}

	//endHeight , 审核时设置区块结束高度
	//审核通过，在根据公链抓取数据里，会自动加入此合约数据，抓取最新数据
	err = service.Contract().UpdateValid(ctx, req.ProtocolCode, req.ContractCode, true, req.OnceHeight, runHeight, req.DataType, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		protocolContractListkey := req.ChainName + consts.PROTOCOL_CONTRACT_LIST_KEY
		rId := req.ProtocolCode + "_" + req.ContractCode
		res, err := g.Redis().Do(ctx, "HGET", protocolContractListkey, rId)
		if err == nil && g.IsEmpty(res) {
			_, err = g.Redis().Do(ctx, "HSET", protocolContractListkey, rId, rId)
			if err != nil {
				g.Log().Error(ctx, "Redis error :", err)
			}
		}
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	}
	return
}

func (p *cContract) UpdateInfo(ctx context.Context, req *v1.ContractUpdateInfoReq) (res *utils.ResponseRes, err error) {

	if g.IsEmpty(req.DeployHeight) && g.IsEmpty(req.OnceHeight) {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("params error"))
	}

	ent, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		return
	}

	if !g.IsEmpty(ent.ContractAddress) {
		err = service.Contract().UpdateInfo(ctx, req.ProtocolCode, req.ContractCode, req.OnceHeight, req.DeployHeight, req.ChainName)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		} else {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		}
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
	}

	return
}

// 查询合约配置列表
func (p *cContract) All(ctx context.Context, req *v1.ContractReq) (res *utils.ResponseRes, err error) {

	result, err := service.Contract().QueryList(ctx, req.ProtocolCode, req.ContractCode, req.IsValid, req.ChainName)
	value := make([]entity.ContractViewEntity, len(result))

	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		chainHeightKey := req.ChainName + "_curr_height_key"
		v, err := g.Redis().Do(ctx, "GET", chainHeightKey)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return nil, err
		}

		key := req.ChainName + consts.ACTIVITY_TASK_KEY
		taskValue, err := g.Redis().Do(ctx, "GET", key)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
			return nil, err
		}

		taskSet := gset.NewSet(true)
		if !g.IsEmpty(v) {
			taskValueList := strings.Split(taskValue.String(), ",")
			for i := 0; i < len(taskValueList); i++ {
				if !g.IsEmpty(taskValueList[i]) && !taskSet.Contains(taskValueList[i]) {
					taskSet.Add(taskValueList[i])
				}
			}
		}

		for i := 0; i < len(result); i++ {
			cent := result[i]
			//如果当前块高 >= 结束块高
			if cent.CurrHeight >= cent.RunHeight {
				cent.CurrHeight = v.Int32()
			}
			taskName := req.ChainName + "_" + cent.ProtocolCode + "_" + cent.ContractCode
			if taskSet.Contains(taskName) {
				cent.TaskStatus = true
			} else {
				cent.TaskStatus = false
			}
			value[i] = cent
		}
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(value, len(value)))
	}
	return res, nil
}

// 查询待审核的合约配置列表
/*func (p *cContract) ToAuditList(ctx context.Context, req *v1.ContractToAuditListReq) (res *utils.ResponseRes, err error) {
	result, err := service.Contract().QueryToAuditList(ctx, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK.WithDataAndTotal(result, len(result)))
	}
	return
}
*/
func (p *cContract) Del(ctx context.Context, req *v1.ContractDelReq) (res *utils.ResponseRes, err error) {
	err = service.Contract().DeleteByCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
	} else {
		g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
	}
	return
}

// 合约任务操作
// 操作类型，1:pause, 2:recover
func (p *cContract) ContractTask(ctx context.Context, req *v1.ContractTaskReq) (res *utils.ResponseRes, err error) {

	if req.Type == 1 {
		//待删除的合约 （需要从活跃的任务列表中删除）
		deleteingVal := req.ChainName + "_" + req.ProtocolCode + "_" + req.ContractCode
		deleteingContractListKey := req.ChainName + consts.DELETING_CONTRACT_LIST_KEY
		_, err = g.Redis().Do(ctx, "HSET", deleteingContractListKey, deleteingVal, deleteingVal)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		} else {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		}
	} else if req.Type == 2 {
		//待恢复的合约 （需要加入到活跃的任务列表中）
		val := req.ProtocolCode + "_" + req.ContractCode
		protocolContractListkey := req.ChainName + consts.PROTOCOL_CONTRACT_LIST_KEY
		_, err = g.Redis().Do(ctx, "HSET", protocolContractListkey, val, val)
		if err != nil {
			g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
		} else {
			//如果有锁，释放锁
			contract, err := service.Contract().QueryByContractCode(ctx, req.ProtocolCode, req.ContractCode, req.ChainName)
			if err != nil {
				g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
				return res, err
			}
			if g.IsEmpty(contract.ContractCode) {
				g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg("the contract is not exists"))
				return nil, err
			}

			cacheKey := req.ChainName + "_lock_transaction_" + string(contract.DataType) + "_" + val
			_, err = g.Redis().Do(ctx, "DEL", cacheKey)
			if err != nil {
				g.RequestFromCtx(ctx).Response.WriteJson(utils.Err.WithMsg(err.Error()))
				return res, err
			}
			g.RequestFromCtx(ctx).Response.WriteJson(utils.OK)
		}
	}
	return
}
