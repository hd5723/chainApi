package blockLoader

import (
	"OnchainParser/internal/chain/ethEvm/chain/baseConfig"
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/chain/ethEvm/utils/loader/eventLogsLoader"
	"OnchainParser/internal/chain/ethEvm/utils/loader/txLoader"
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"OnchainParser/internal/web3/web3Client"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func DoScanLastBlock(chainId int, chainName string) (err error) {
	curTime := time.Now()
	//开启链路追踪
	ctx := gctx.New()

	lockKey := chainName + "_" + consts.LOCK_TRANSACTION_KEY + "_LAST_BLOCK"
	v, err := g.Redis().Do(ctx, "GET", lockKey)
	if err != nil {
		return err
	}

	// 如果没有获取到锁
	if v != nil && v.Bool() == true {
		g.Log().Info(ctx, "error:", "DoScanLastBlock the lock is not null or the lock is true , key:", lockKey)
		return gerror.New("DoScanLastBlock the lock is not null or the lock is true")
	}

	// 获取到锁后，加锁，设置锁超时时间,单位秒
	g.Redis().Do(ctx, "SETEX", lockKey, 600, true)
	// 业务处理完毕，释放锁
	defer g.Redis().Do(ctx, "DEL", lockKey)

	g.Log().Info(ctx, "DoScanLastBlock TransactionJob start")

	rcd := service.Chain().QueryOneByChainId(ctx, int32(chainId), chainName)
	links := strings.Split(rcd.RpcUrls, ",")
	ethClient, link, err := web3Client.GetClientByLinks(ctx, links, chainName)
	if err != nil {
		return
	}

	//search enable contract list
	result, err := service.Contract().QueryEnableList(ctx, chainName)
	if err != nil {
		return err
	}
	if g.IsEmpty(result) || len(result) == 0 {
		g.Log().Info(ctx, "DoScanLastBlock no task , chainName:", chainName)
		return nil
	}

	ctx = gctx.New()
	err = doScanLastBlock(ctx, ethClient, link, result, chainName)
	if err != nil {
		return err
	}

	ethClient.Close()
	g.Log().Info(ctx, "DoScanLastBlock TransactionJob doing time:", time.Now().Sub(curTime))
	return nil
}

func DoScanHistoryBlock(ctx context.Context, cli *ethclient.Client, link string, ent entity.ContractEntity, chainName string) (err error) {
	curTime := time.Now()

	lockKey := chainName + "_" + consts.LOCK_TRANSACTION_KEY + "_" + strconv.Itoa(consts.CONTRACT_DATA_BLOCK_TYPE) + "_" + ent.ProtocolCode + "_" + ent.ContractCode
	v, err := g.Redis().Do(ctx, "GET", lockKey)
	if err != nil {
		return err
	}

	// 如果没有获取到锁
	if v != nil && v.Bool() == true {
		g.Log().Info(ctx, "error:", "DoScanHistoryBlock the lock is not null or the lock is true , key:", lockKey)
		return gerror.New("DoScanHistoryBlock the lock is not null or the lock is true")
	}

	// 获取到锁后，加锁，设置锁超时时间,单位秒
	g.Redis().Do(ctx, "SETEX", lockKey, 600, true)
	// 业务处理完毕，释放锁
	defer g.Redis().Do(ctx, "DEL", lockKey)

	g.Log().Info(ctx, "DoScanHistoryBlock TransactionJob start , ContractCode", ent.ContractCode)
	err = doScanBlock(ctx, cli, link, ent, chainName)

	if err != nil {
		return err
	}
	g.Log().Info(ctx, "DoScanHistoryBlock TransactionJob doing time:", time.Now().Sub(curTime), " , ContractCode:", ent.ContractCode)
	return nil
}

func DoSubscribeHistoryContract(ent entity.ContractEntity, chainName string) error {
	curTime := time.Now()
	//开启链路追踪
	ctx := gctx.New()

	lockKey := chainName + "_" + consts.LOCK_TRANSACTION_KEY + "_" + strconv.Itoa(consts.CONTRACT_DATA_EVENTLOGS_TYPE) + "_" + ent.ProtocolCode + "_" + ent.ContractCode
	v, err := g.Redis().Do(ctx, "GET", lockKey)
	if err != nil {
		return err
	}

	//如果没有获取到锁
	if v != nil && v.Bool() == true {
		g.Log().Info(ctx, "error:", "DoSubscribeHistoryContract the lock is not null or the lock is true", " , key:", lockKey)
		return nil
	}

	rcd := service.Chain().QueryOneByChainName(ctx, chainName)
	links := strings.Split(rcd.RpcUrls, ",")
	ethClient, link, err := web3Client.GetClientByLinks(ctx, links, chainName)
	if err != nil {
		return err
	}
	g.Log().Info(ctx)

	//获取到锁后，加锁，设置锁超时时间,单位秒
	g.Redis().Do(ctx, "SETEX", lockKey, 600, true)

	//业务处理完毕，释放锁
	defer g.Redis().Do(ctx, "DEL", lockKey)

	g.Log().Info(ctx, "DoSubscribeHistoryContract TransactionJob start , ContractCode", ent.ContractCode)
	err = doSubscribeContract(ctx, ethClient, link, ent, chainName)
	if err != nil {
		return err
	}

	g.Log().Info(ctx, "DoSubscribeHistoryContract TransactionJob doing time:", time.Now().Sub(curTime), " , ContractCode:", ent.ContractCode)
	return nil
}

func Undo(ctx context.Context) {

}

func doScanLastBlock(ctx context.Context, cli *ethclient.Client, link string, ents []entity.ContractEntity, chainName string) (err error) {
	if ents == nil || len(ents) == 0 {
		return
	}

	curTime := time.Now()
	lastBlockNumber := baseConfig.CHAIN_LAST_HEIGHT

	chainHeightKey := chainName + "_curr_height_key"
	curTime = time.Now()
	v, err := g.Redis().Do(ctx, "GET", chainHeightKey)
	g.Log().Info(ctx, "redis doScanLastBlock.get_curr_height_key doing time:", time.Now().Sub(curTime))
	if err != nil {
		return err
	}

	//缓存里的最新区块，比最新区块小1；否则return
	if v.Int64()+1 >= lastBlockNumber {
		return
	}

	//设置每次处理区块数量
	pnum, err := g.Cfg().Get(ctx, "chain.onceLength")
	if err != nil {
		return
	}

	//每次处理区块数量20
	onceLength := pnum.Int64()

	//区块开始高度设置
	var fromBlock = lastBlockNumber - onceLength
	// 如果已处理高度 <  最新高度 -每次处理区块数量
	if v.Val() != nil {
		fromBlock = v.Int64()
	}

	//上一次区块高基础上 +1
	fromBlock = fromBlock + 1

	//区块结束高度设置
	var toBlock = fromBlock + onceLength - 1
	if toBlock > lastBlockNumber-1 {
		toBlock = lastBlockNumber - 1
	}

	nowTime := time.Now()

	var blockNumberList = make([]int64, toBlock-fromBlock+1)
	for i := 0; i < int(toBlock-fromBlock)+1; i++ {
		blockNumberList[i] = fromBlock + int64(i)
	}
	err = parseBlockByBlockNumListWithEnts(ctx, cli, link, blockNumberList, ents, chainName)
	if err != nil {
		return err
	}

	curTime = time.Now()
	// 批量执行SQL
	err = dao.Transaction.BatchExecSql(ctx, chainName)
	g.Log().Info(ctx, "db dao.Transaction , ToBlock:", toBlock, "，CurBlock:", fromBlock, ",time run:", time.Now().Sub(nowTime))
	if err != nil {
		return err
	}

	//任务执行完成，结束
	curTime = time.Now()
	_, err = g.Redis().Do(ctx, "SET", chainHeightKey, toBlock)
	g.Log().Info(ctx, "redis doScanLastBlock.set_curr_height_key doing time:", time.Now().Sub(curTime))
	if err != nil {
		return err
	}
	g.Log().Info(ctx, "doScanLastBlock finish  , ToBlock:", toBlock, "，CurBlock:", fromBlock, ",time run:", time.Now().Sub(nowTime))
	return nil
}

func doScanBlock(ctx context.Context, cli *ethclient.Client, link string, ent entity.ContractEntity, chainName string) (err error) {
	ent, err = service.Contract().QueryByContractCode(ctx, ent.ProtocolCode, ent.ContractCode, chainName)
	if err != nil {
		return
	}
	if g.IsEmpty(ent.ContractCode) {
		return
	}

	//一次最新区块爬取任务的区块高
	onceLength, err := g.Cfg().Get(ctx, "chain.onceLength")
	if err != nil {
		return
	}

	//冗余一次最新区块爬取任务的区块高
	//scan历史记录只到endHeight + 逐块处理的onceLength
	lastBlockNumber := int64(ent.RunHeight) + onceLength.Int64()

	//每次请求的区块数量
	var addLength = int64(ent.OnceHeight)

	//区块开始高度设置
	var fromBlock = int64(ent.DeployHeight)
	if ent.DeployHeight < ent.CurrHeight {
		fromBlock = int64(ent.CurrHeight)
	}

	// 终止条件
	if fromBlock >= int64(ent.RunHeight) {
		hrcd, err := service.ContractHistoryJob().QueryByContract(ctx, ent.ProtocolCode, ent.ContractCode, chainName)
		if err != nil {
			return err
		}
		//如果在已完成任务里查找不到数据，就插入一条数据
		if g.IsEmpty(hrcd.IsHistoryFinish) {
			var in model.ContractHistoryJobCreateInput
			in.IsHistoryFinish = true
			in.ProtocolCode = ent.ProtocolCode
			in.ContractCode = ent.ContractCode
			in.UpdateTime = int32(gtime.Timestamp())
			err = service.ContractHistoryJob().Create(ctx, in, chainName)
			if err != nil {
				return err
			}
		}
	}

	//开始区块+1
	fromBlock = fromBlock + 1

	//区块结束高度设置
	var toBlock = fromBlock + addLength
	if toBlock > lastBlockNumber-1 {
		toBlock = lastBlockNumber - 1
	}

	curTime := time.Now()
	var blockNumberList = make([]int64, toBlock-fromBlock+1)
	for i := 0; i < int(toBlock-fromBlock)+1; i++ {
		blockNumberList[i] = fromBlock + int64(i)
	}

	err = parseBlockByBlockNumList(ctx, cli, link, blockNumberList, ent, chainName)
	if err != nil {
		return err
	}
	g.Log().Debug(ctx, ",ToBlock:", toBlock, ",Contract:", ent.ContractAddress, "，CurBlock:", fromBlock, ",time run:", time.Now().Sub(curTime))

	// 批量执行SQL
	dao.Transaction.BatchExecSql(ctx, chainName)

	//任务执行完成后，修改合约配置里的当前块高，值：区块结束高度
	//当任务异常中断重启时，可以从最后一次成功记录里，获取到开始区块，重新开始任务
	//  redis 保存 cur_height 信息
	err = service.Contract().UpdateCurHeight(ctx, ent.ContractCode, ent.ProtocolCode, int32(toBlock), chainName)
	if err != nil {
		//如果有错误，或者执行失败，中断
		return err
	}

	g.Log().Info(ctx, "doScanBlock finish Contract:", ent.ContractAddress, ",time run:", time.Now().Sub(curTime))
	return nil
}

/*
*
带有合约配置信息的抓取区块交易信息
*/
func doSubscribeContract(ctx context.Context, ethClient *ethclient.Client, link string, ent entity.ContractEntity, chainName string) (err error) {
	g.Log().Info(ctx, "doSubscribeContract start ContractCode:", ent.ContractCode)
	ent, err = service.Contract().QueryByContractCode(ctx, ent.ProtocolCode, ent.ContractCode, chainName)
	if err != nil {
		return err
	}

	//设置每次处理区块数量
	pnum, err := g.Cfg().Get(ctx, "chain.onceLength")
	if err != nil {
		return err
	}

	//冗余一次最新区块爬取任务的区块高
	//scan历史记录只到endHeight + 逐块处理的onceLength
	lastBlockNumber := int64(ent.RunHeight) + pnum.Int64()

	//每次请求的区块数量
	var addLength = int64(ent.OnceHeight)

	//区块开始高度设置
	var fromBlock = int64(ent.DeployHeight)
	if ent.DeployHeight < ent.CurrHeight {
		fromBlock = int64(ent.CurrHeight)
	}

	// 终止条件
	if fromBlock >= int64(ent.RunHeight) {
		hrcd, err := service.ContractHistoryJob().QueryByContract(ctx, ent.ProtocolCode, ent.ContractCode, chainName)
		if err != nil {
			return err
		}
		//如果在已完成任务里查找不到数据，就插入一条数据
		// hrcd.IsHistoryFinish 是基础类型，一定会有值，不能用作判断
		if g.IsEmpty(hrcd.ContractCode) {
			var in model.ContractHistoryJobCreateInput
			in.IsHistoryFinish = true
			in.ProtocolCode = ent.ProtocolCode
			in.ContractCode = ent.ContractCode
			in.UpdateTime = int32(gtime.Timestamp())
			err = service.ContractHistoryJob().Create(ctx, in, chainName)
			if err != nil {
				return err
			}
		}
	}

	//开始区块+1
	fromBlock += 1

	//区块结束高度设置
	var toBlock = fromBlock + addLength
	if toBlock > lastBlockNumber-1 {
		toBlock = lastBlockNumber - 1
	}

	g.Log().Info(ctx, "doSubscribeContract subscribeContract start ContractCode:", ent.ContractCode)
	nowTime := time.Now()
	result := 0
	//当返回值为0时，继续执行
	for err == nil && result == 0 {
		result, err = subscribeContract(ctx, ethClient, link, ent, fromBlock, toBlock, chainName)
		//当返回值为0时, fromBlock、toBlock重置，追加一个onceLength
		if err == nil && result == 0 {
			// redis保存cur_height , 无数据时也更新cur_height
			err = service.Contract().UpdateCurHeight(ctx, ent.ContractCode, ent.ProtocolCode, int32(toBlock), chainName)
			if err != nil {
				return err
			}

			fromBlock = toBlock + 1
			//fromBlock > 结束区块时，跳出当前循环，设置curr_height = end_height
			if fromBlock >= int64(ent.RunHeight) {
				toBlock = int64(ent.RunHeight)
				break
			}
			toBlock = fromBlock + addLength
			if toBlock > lastBlockNumber-1 {
				toBlock = lastBlockNumber - 1
			}
		}
	}
	g.Log().Info(ctx, "doSubscribeContract subscribeContract end ContractCode:", ent.ContractCode)
	if err != nil {
		//失败之后中断执行
		return err
	} else {
		// 批量执行SQL
		err = dao.Transaction.BatchExecSql(ctx, chainName)
		if err != nil {
			return err
		}
		//任务执行完成后，修改合约配置里的当前块高，值：区块结束高度
		//当任务异常中断重启时，可以从最后一次成功记录里，获取到开始区块，重新开始任务
		err = service.Contract().UpdateCurHeight(ctx, ent.ContractCode, ent.ProtocolCode, int32(toBlock), chainName)
		if err != nil {
			return err
		}

		g.Log().Info(ctx, "doSubscribeContract finish Contract:", ent.ContractAddress, "FromBlock:", fromBlock, ",ToBlock:", toBlock, ",time run:", time.Now().Sub(nowTime))
	}
	return nil
}

// 并发解析block
func parseBlockByBlockNumList(ctx context.Context, client *ethclient.Client, link string, blockNumberList []int64, ent entity.ContractEntity, chainName string) (err error) {
	workList := make(chan consts.BlockExecuteInfo)
	//设置go pool的并发工作
	pnum, err := g.Cfg().Get(ctx, "chain.blockNum")
	if err != nil {
		return err
	}

	var poolNum int
	if pnum.Int() > len(blockNumberList) {
		poolNum = len(blockNumberList)
	} else {
		poolNum = pnum.Int()
	}
	pool := grpool.New(poolNum)

	for i := 0; i < len(blockNumberList); i++ {
		blockNumber := blockNumberList[i]
		pool.Add(ctx, func(ctx context.Context) {
			parseBlock(ctx, client, link, blockNumber, ent, chainName, workList)
		})
	}

	blockExecuteInfos := make([]consts.BlockExecuteInfo, len(blockNumberList))
	for i := 0; i < len(blockNumberList); i++ {
		blockExecuteInfos[i] = <-workList
	}
	close(workList)

	errCount := 0
	for i := 0; i < len(blockNumberList); i++ {
		if blockExecuteInfos[i].Err != nil {
			// 此条区块设置为false
			errCount++
		}
	}

	if errCount > 0 {
		return gerror.New("the blockList parser fail!")
	}
	return
}

func parseBlock(ctx context.Context, client *ethclient.Client, link string, blockNumber int64, ent entity.ContractEntity, chainName string, rootWorkList chan consts.BlockExecuteInfo) {

	var result consts.BlockExecuteInfo
	result.BlockNumber = blockNumber

	//检查上一次是否执行成功，上一次未成功的继续执行；成功的标记err为nil返回
	blockKey := ent.ProtocolCode + "_" + ent.ContractCode + "_" + strconv.Itoa(int(blockNumber))
	blocksKey := chainName + "_" + blockKey
	res, err := g.Redis().Do(ctx, "HGET", blocksKey, blockNumber)
	if err != nil {
		result.Err = err
		rootWorkList <- result
		return
	}

	// res == true, 已经处理过，忽略
	if res != nil && res.Bool() {
		result.Err = nil
		rootWorkList <- result
		return
	}

	curTime := time.Now()
	block, err := client.BlockByNumber(ctx, big.NewInt(blockNumber))
	g.Log().Info(ctx, "RPC block.parseBlock.client.BlockByNumber doing time:", time.Now().Sub(curTime), " ,link:", link)
	if err != nil {
		g.Log().Error(ctx, "RPC block.parseBlock.client.BlockByNumber error:", err)
		result.Err = err
		rootWorkList <- result
		return
	}

	var execType = 0 //合约类型

	txTodoTask := make(chan consts.TxExecuteInfo)

	start := gtime.TimestampMilli()

	var num = 0
	var notMatchNum = 0

	//设置go pool的并发工作
	pnum, err := g.Cfg().Get(ctx, "chain.poolNum")
	if err != nil {
		result.Err = err
		rootWorkList <- result
		return
	}

	var poolNum int
	if pnum.Int() > block.Transactions().Len() {
		poolNum = block.Transactions().Len()
	} else {
		poolNum = pnum.Int()
	}
	pool := grpool.New(poolNum)
	for i := 0; i < block.Transactions().Len(); i++ {
		tx := block.Transactions()[i]
		if tx.To() == nil {
			continue
		}
		to := tx.To().Hex()

		if strings.Compare(strings.ToLower(to), strings.ToLower(ent.ContractAddress)) != 0 {
			notMatchNum++
			continue
		} else {
			num++
			newTx := tx
			//使用通道+并发执行
			pool.Add(ctx, func(ctx context.Context) {
				txLoader.Do(ctx, client, link, newTx, block, ent, chainName, txTodoTask, execType)
			})
		}
	}

	txExecuteInfos := make([]consts.TxExecuteInfo, num)

	for i := 0; i < num; i++ {
		if txTodoTask != nil {
			txExecuteInfos[i] = <-txTodoTask
		}
	}
	close(txTodoTask)

	g.Log().Info(ctx, "Contract:", ent.ContractAddress, "blockNumber:", blockNumber, ",match Transaction num:", num, ",not match Transaction num:", notMatchNum, " ,time spent:", gtime.TimestampMilli()-start)
	blockTxsKey := baseConfig.CHAIN_NAME + "_" + ent.ProtocolCode + "_" + ent.ContractCode + "_" + strconv.FormatInt(int64(block.NumberU64()), 10)
	errCount := 0
	for i := 0; i < num; i++ {
		if txExecuteInfos[i].Err != nil {
			_, err = g.Redis().Do(ctx, "HSET", blockTxsKey, txExecuteInfos[i].TxHash, false)
			if err != nil {
				result.Err = err
				rootWorkList <- result
				return
			}
			errCount++
		} else {
			// 此条交易设置为true
			_, err = g.Redis().Do(ctx, "HSET", blockTxsKey, txExecuteInfos[i].TxHash, true)
			if err != nil {
				result.Err = err
				rootWorkList <- result
				return
			}
		}
	}
	if errCount > 0 {
		_, err = g.Redis().Do(ctx, "HSET", blocksKey, blockNumber, false)
		result.Err = gerror.New("the parseBlock execute fail")
		rootWorkList <- result
		return
	} else {
		//使用完后回收
		_, err = g.Redis().Do(ctx, "DEL", blockTxsKey)
	}
	result.Err = nil
	rootWorkList <- result
	return
}

// 并发解析block
func parseBlockByBlockNumListWithEnts(ctx context.Context, client *ethclient.Client, link string, blockNumberList []int64, ents []entity.ContractEntity, chainName string) (err error) {
	curTime := time.Now()

	workList := make(chan consts.BlockExecuteInfo)
	//设置go pool的并发工作
	pnum, err := g.Cfg().Get(ctx, "chain.blockNum")
	if err != nil {
		return err
	}

	var poolNum int
	if pnum.Int() > len(blockNumberList) {
		poolNum = len(blockNumberList)
	} else {
		poolNum = pnum.Int()
	}
	pool := grpool.New(poolNum)

	for i := 0; i < len(blockNumberList); i++ {
		var blockInfo consts.BlockInfo
		blockInfo.BlockNumber = blockNumberList[i]
		pool.Add(ctx, func(ctx context.Context) {
			parseBlockWithEnts(ctx, client, link, blockInfo, ents, chainName, workList)
		})
	}

	blockExecuteInfos := make([]consts.BlockExecuteInfo, len(blockNumberList))
	for i := 0; i < len(blockNumberList); i++ {
		blockExecuteInfos[i] = <-workList
	}
	close(workList)

	errCount := 0
	for i := 0; i < len(blockNumberList); i++ {
		if blockExecuteInfos[i].Err != nil {
			// 此条区块设置为false
			errCount++
		}
	}

	if errCount > 0 {
		return gerror.New("the blockList parser fail!")
	}
	//当前任务完成，删除缓存
	for i := 0; i < len(blockNumberList); i++ {
		blockTxsKey := baseConfig.CHAIN_NAME + "_ALL_" + strconv.FormatInt(blockNumberList[i], 10)
		curTime := time.Now()
		g.Redis().Do(ctx, "DEL", blockTxsKey)
		g.Log().Info(ctx, "redis parseBlockByBlockNumListWithEnts.del_blockTxsKey doing time:", time.Now().Sub(curTime))
	}
	g.Log().Info(ctx, "parseBlockByBlockNumListWithEnts 并发汇总 doing time:", time.Now().Sub(curTime), " , link:", link)
	return
}

func parseBlockWithEnts(ctx context.Context, client *ethclient.Client, link string, blockInfo consts.BlockInfo, ents []entity.ContractEntity, chainName string, rootWorkList chan consts.BlockExecuteInfo) (err error) {
	var result consts.BlockExecuteInfo
	result.BlockNumber = blockInfo.BlockNumber

	curTime := time.Now()
	block, err := client.BlockByNumber(ctx, big.NewInt(blockInfo.BlockNumber))
	g.Log().Info(ctx, "RPC block.parseBlockWithEnts.client.BlockByNumber doing time:", time.Now().Sub(curTime), " , link:", link)
	if err != nil {
		g.Log().Error(ctx, "RPC block.parseBlockWithEnts.client.BlockByNumber error:", err)
		result.Err = err
		rootWorkList <- result
		return
	}

	txTodoTask := make(chan consts.TxExecuteInfo)

	start := gtime.TimestampMilli()

	var contractSet gset.StrSet
	var contractMap gmap.Map
	for _, ent := range ents {
		contractSet.Add(strings.ToLower(ent.ContractAddress))
		contractMap.SetIfNotExist(strings.ToLower(ent.ContractAddress), ent)
	}

	var num = 0
	var notMatchNum = 0

	//设置go pool的并发工作
	pnum, err := g.Cfg().Get(ctx, "chain.poolNum")
	if err != nil {
		result.Err = err
		rootWorkList <- result
		return
	}

	var poolNum int
	if pnum.Int() > block.Transactions().Len() {
		poolNum = block.Transactions().Len()
	} else {
		poolNum = pnum.Int()
	}
	var execType = 1 //基于所有的合约数据
	pool := grpool.New(poolNum)
	for i := 0; i < block.Transactions().Len(); i++ {
		tx := block.Transactions()[i]
		if tx.To() == nil {
			g.Log().Debug(ctx, "QueryBlockInfoByBlockNumWithEnts to is nil , txHash :", tx.Hash().Hex())
			continue
		}
		to := tx.To().Hex()

		if !contractSet.Contains(strings.ToLower(to)) {
			notMatchNum++
			continue
		} else {
			num++
			newTx := tx
			var ent = contractMap.Get(strings.ToLower(to)).(entity.ContractEntity)
			//使用通道+并发执行
			pool.Add(ctx, func(ctx context.Context) {
				txLoader.Do(ctx, client, link, newTx, block, ent, chainName, txTodoTask, execType)
			})
		}
	}
	txExecuteInfos := make([]consts.TxExecuteInfo, num)
	for i := 0; i < num; i++ {
		txExecuteInfos[i] = <-txTodoTask
	}
	close(txTodoTask)
	g.Log().Info(ctx, "QueryBlockInfoByBlockNumWithEnts time spent:", gtime.TimestampMilli()-start)
	blockTxsKey := baseConfig.CHAIN_NAME + "_ALL_" + strconv.FormatInt(int64(block.NumberU64()), 10)
	errCount := 0
	for i := 0; i < num; i++ {
		if txExecuteInfos[i].Err != nil {
			_, err = g.Redis().Do(ctx, "HSET", blockTxsKey, txExecuteInfos[i].TxHash, false)
			if err != nil {
				result.Err = err
				rootWorkList <- result
				return
			}
			errCount++
		}
	}

	if errCount > 0 {
		result.Err = gerror.New("the parseBlockWithEnts execute fail")
		rootWorkList <- result
		return
	}
	g.Log().Info(ctx, "doScanLastBlockparseBlockWithEnts finish  , blockNumber:", blockInfo.BlockNumber, ",time run:", time.Now().Sub(curTime))
	result.Err = nil
	rootWorkList <- result
	return
}

func subscribeContract(ctx context.Context, client *ethclient.Client, link string, ent entity.ContractEntity, fromBlock int64, toBlock int64, chainName string) (result int, err error) {
	result = 200
	contractAddress := common.HexToAddress(ent.ContractAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	eClient, err := ethclient.Dial(link)
	if err != nil {
		g.Log().Error(ctx, "web3 rpc.Dial err")
		return
	}

	curTime := time.Now()
	g.Log().Info(ctx, "RPC block.subscribeContract.client.FilterLogs start fromBlock:", fromBlock, " ,toBlock:", toBlock, " , contract_code:", ent.ContractCode, "ent.ContractAddress:", ent.ContractAddress, " ,contractAddress:", contractAddress, " ,link:", link)
	logs, err := eClient.FilterLogs(ctx, query)
	g.Log().Info(ctx, "RPC block.subscribeContract.client.FilterLogs doing time:", time.Now().Sub(curTime))
	if err != nil {
		g.Log().Error(ctx, "RPC block.subscribeContract.client.FilterLogs error:", err)
		result = 1
		return result, err
	}

	logsLen := len(logs)
	if logsLen == 0 {
		result = 0
		g.Log().Info(ctx, "block.subscribeContract. result is nil , fromBlock:", fromBlock, " ,toBlock:", toBlock, "contract:", ent.ContractAddress, " ,chainName:", chainName)
		return result, nil
	}
	eventSigArray := strings.Split(ent.SignedListenerEvent, ", ")

	//设置go pool的并发工作
	pnum, err := g.Cfg().Get(ctx, "chain.poolNum")
	if err != nil {
		result = 0
		return result, err
	}

	var poolNum int
	if pnum.Int() > logsLen {
		poolNum = logsLen
	} else {
		poolNum = pnum.Int()
	}
	pool := grpool.New(poolNum)
	txTodoTask := make(chan consts.TxExecuteInfo)
	num := 0
	setNum := 0
	logsSet := make([]types.Log, logsLen)
	txSet := gset.NewSet(true)
	for _, vLog := range logs {
		sign := vLog.Topics[0].Hex()
		if utils.Contains(eventSigArray, sign) {
			logsSet[num] = vLog
			num++
			if !txSet.Contains(vLog.TxHash.Hex()) {
				txSet.Add(vLog.TxHash.Hex())
				setNum++
				var newTx consts.TxInfo
				newTx.TxHash = vLog.TxHash.String()
				//使用通道+并发执行
				pool.Add(ctx, func(ctx context.Context) {
					//处理交易信息
					txLoader.DoTx(ctx, client, link, newTx, toBlock, ent, chainName, txTodoTask)
				})

			}

		}
	}

	txExecuteInfos := make([]consts.TxExecuteInfo, setNum)
	for i := 0; i < setNum; i++ {
		txExecuteInfos[i] = <-txTodoTask
	}
	close(txTodoTask)
	blockKey := ent.ProtocolCode + "_" + ent.ContractCode
	blockTxsKey := baseConfig.CHAIN_NAME + "_" + blockKey + "_" + strconv.FormatInt(toBlock, 10)
	errCount := 0
	for i := 0; i < setNum; i++ {
		if txExecuteInfos[i].Err != nil {
			_, err = g.Redis().Do(ctx, "HSET", blockTxsKey, txExecuteInfos[i].TxHash, false)
			if err != nil {
				result = 0
				return result, err
			}
			errCount++
		}
	}

	if errCount > 0 {
		result = 0
		return result, err
		return
	}

	if num > 0 {
		err = eventLogsLoader.DoEventLogs(ctx, logsSet, ent, chainName)
		if err != nil {
			result = 3
			return result, err
		}
	}
	return result, nil
}
