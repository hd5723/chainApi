package task

import (
	"OnchainParser/internal/chain/ethEvm/chain/baseConfig"
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/chain/ethEvm/utils/loader/chainLoader"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/web3/web3Client"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
)

func TaskManage(ctx context.Context) {

	cron := gcron.New()
	// ------------   BSC  -----------------------------------------------------------------
	// 每2秒执行一次 , 爬取最新区块任务，只有一个任务
	_, err := gcron.Add(ctx, "*/2 * * * * *", func(ctx context.Context) {
		chainLoader.DoScanLastBlock(baseConfig.CHAIN_NAME, baseConfig.CHAIN_ID)
	})
	if err != nil {
		g.Log().Error(ctx, "Add DoScanLastBlock cron fail ")
	}

	// ------------   监听任务完成状态  -------------------------------------------------------
	//每10秒执行一次 , 监听区块历史数据爬取任务完成状态
	//任务完成状态判断：cur_height >= end_height
	//通过cron搜索job
	taskListenerCtx := gctx.New()
	_, err = gcron.Add(taskListenerCtx, "*/10 * * * * *", func(ctx context.Context) {
		taskListener(ctx, cron, baseConfig.CHAIN_NAME)
	})
	if err != nil {
		g.Log().Error(ctx, "Add taskListener cron fail ")
	}

	// 更新当前公链的最新块高
	// 减少RPC请求
	_, err = gcron.Add(taskListenerCtx, "*/5 * * * * *", func(ctx context.Context) {
		SetChainLastHeight(ctx)
	})
	if err != nil {
		g.Log().Error(ctx, "Add taskListener cron fail ")
	}

	// ------------   从缓存记载protocolList  ------------------------------------------------
	// 每10秒执行一次 , 从缓存加载protocolList
	// 不中断程序的情况下，从缓存加载 protocolList ， 加载后删除
	addFromCacheCtx := gctx.New()
	_, err = gcron.Add(addFromCacheCtx, "*/10 * * * * *", func(ctx context.Context) {
		taskAddFromCache(ctx, cron, baseConfig.CHAIN_NAME, baseConfig.CHAIN_ID)
	})
	if err != nil {
		g.Log().Error(ctx, "Add taskAddFromCache cron fail ")
	}

	g.Log().Info(ctx, "TaskManage running ...")

}

func SetChainLastHeight(ctx context.Context) error {
	rcd := service.Chain().QueryOneByChainId(ctx, int32(baseConfig.CHAIN_ID), baseConfig.CHAIN_NAME)
	links := strings.Split(rcd.RpcUrls, ",")
	ethClient, _, err := web3Client.GetClientByLinks(ctx, links, baseConfig.CHAIN_NAME)
	if err != nil {
		return err
	}

	lastBlockNumber, err := web3Client.GetLastBlockNumber(ctx, ethClient)
	if err != nil {
		return err
	}
	baseConfig.CHAIN_LAST_HEIGHT = lastBlockNumber
	return nil
}

func doForScanBlockAllocation(ctx context.Context, cron *gcron.Cron, ent entity.ContractEntity, chainName string, chainId int) error {

	//注册job时，以 chainName_protocolCode_contractCode 为job的name
	jobName := chainName + "_" + ent.ProtocolCode + "_" + ent.ContractCode
	_, err := cron.Add(ctx, "*/3 * * * * *", func(ctx context.Context) {
		chainLoader.DoScanHistoryBlock(ent, chainName, chainId)
	}, jobName)
	if err != nil {
		return gerror.New("Add cron fail , jobName: " + jobName)
	}
	return nil
}

func doForSubscribeContractAllocation(ctx context.Context, cron *gcron.Cron, ent entity.ContractEntity, chainName string, chainId int) error {

	g.Log().Debug(ctx, "doForSubscribeContractAllocation ent:", ent.ContractAddress, " , chainName:", chainName)

	//注册job时，以 chainName_protocolCode_contractCode 为job的name
	jobName := chainName + "_" + ent.ProtocolCode + "_" + ent.ContractCode
	_, err := cron.Add(ctx, "*/3 * * * * *", func(ctx context.Context) {
		chainLoader.DoSubscribeHistoryContract(ent, chainName, chainId)
	}, jobName)
	if err != nil {
		return gerror.New("Add cron fail , jobName: " + jobName)
	}
	return nil
}

func taskListener(ctx context.Context, cron *gcron.Cron, chainName string) {
	//遍历已注册的任务
	entries := cron.Entries()
	g.Log().Info(ctx, " the activity task size:", len(entries), " , activity task:", entries)

	activityTaskList := make([]string, len(entries))
	key := chainName + consts.ACTIVITY_TASK_KEY
	if len(entries) == 0 {
		g.Redis().Do(ctx, "SET", key, "")
		return
	}

	//查询待删除的合约
	deleteingContractListKey := chainName + consts.DELETING_CONTRACT_LIST_KEY
	deleteingContractList, err := g.Redis().Do(ctx, "HGETALL", deleteingContractListKey)

	for i := 0; i < len(entries); i++ {
		v := entries[i]
		if err == nil && !deleteingContractList.IsEmpty() {
			num := 0
			for _, deleteingContract := range deleteingContractList.Array() {
				num++
				//HGETALL key value占2行，取其一
				if num%2 == 1 {
					continue
				}
				//匹配当前活跃的合约
				//从待删除列表、活跃任务列表中删除
				if v.Name == deleteingContract.(string) {
					cron.Remove(v.Name)
					g.Redis().Do(ctx, "HDEL", deleteingContractListKey, v.Name)
				}
			}
		}

		jobName := strings.Split(v.Name, "_")
		if chainName == jobName[0] {
			cent, err := service.Contract().QueryByContractCode(ctx, jobName[1], jobName[2], chainName)
			if err != nil {
				continue //跳过，执行下一个
			}

			//如果合约不存在
			if g.IsEmpty(cent) || g.IsEmpty(cent.ContractAddress) {
				cron.Remove(v.Name)
			}

			//如果合约存在,但数据有问题
			if g.IsEmpty(cent.RunHeight) {
				cron.Remove(v.Name)
			}

			//如果合约的当前爬取高度 大于等于 结束高度，从 cron 里移除
			if cent.CurrHeight >= cent.RunHeight {
				cron.Remove(v.Name)
			} else {
				activityTaskList[i] = v.Name
			}
		}
	}
	g.Redis().Do(ctx, "SET", key, strings.Join(activityTaskList, ","))
}

func taskAddFromCache(ctx context.Context, cron *gcron.Cron, chainName string, chainId int) error {

	protocolContractListkey := chainName + consts.PROTOCOL_CONTRACT_LIST_KEY
	result, err := g.Redis().Do(ctx, "HGETALL", protocolContractListkey)
	if err != nil {
		return err
	}
	num := 0
	for _, v := range result.Array() {
		num++
		//HGETALL key value占2行，取其一
		if num%2 == 1 {
			continue
		}
		if g.IsEmpty(v.(string)) {
			g.Redis().Do(ctx, "HDEL", protocolContractListkey, v)
			continue
		}
		rId := strings.Split(v.(string), "_")
		protocolCode := rId[0]
		contractCode := rId[1]

		pent, err := service.Protocol().QueryOneByProtocolCode(ctx, protocolCode, chainName)
		if err != nil {
			return err
		}
		//找不到数据（配置错误，或未通过审核)
		if g.IsEmpty(pent.ProtocolCode) || !pent.IsValid {
			g.Redis().Do(ctx, "HDEL", protocolContractListkey, v)
			continue
			//return gerror.New("the protocolCode is not exists or not pass audit")
		}

		ent, err := service.Contract().QueryByContractCode(ctx, protocolCode, contractCode, chainName)
		if err != nil {
			g.Redis().Do(ctx, "HDEL", protocolContractListkey, v)
			continue
			//return gerror.New("the contractCode is not exists or not pass audit")
		}

		//找不到数据（配置错误，或未通过审核)
		if g.IsEmpty(ent.ContractAddress) || !ent.IsValid {
			g.Redis().Do(ctx, "HDEL", protocolContractListkey, v)
			continue
		}

		// 分配[bscConfig.CHAIN_NAME]任务
		// 根据合约配置信息分配
		// 通过cron添加job
		if ent.DataType == consts.CONTRACT_DATA_BLOCK_TYPE {
			doForScanBlockAllocation(ctx, cron, ent, chainName, chainId)
		} else {
			doForSubscribeContractAllocation(ctx, cron, ent, chainName, chainId)
		}

		g.Redis().Do(ctx, "HDEL", protocolContractListkey, v)
	}
	return nil
}
