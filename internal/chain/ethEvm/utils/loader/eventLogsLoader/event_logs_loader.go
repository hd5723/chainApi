package eventLogsLoader

import (
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
)

func init() {

}

// 扫描区块爬取方式，处理logs数据
func Do(ctx context.Context, receipt *types.Receipt, txHash string, ent entity.ContractEntity, chainName string) (err error) {
	LogsLen := len(receipt.Logs)
	ins := make([]model.TransactionLogCreateInput, LogsLen)
	num := 0
	for i := 0; i < LogsLen; i++ {
		vLog := receipt.Logs[i]
		hashLog, err := service.TransactionLog().QueryExistsByTxHashAndIndex(ctx, txHash, int32(vLog.Index), chainName)
		if err != nil {
			return err
		}

		//如果不存在，就创建
		if !hashLog {
			var in model.TransactionLogCreateInput
			in.TxHash = receipt.TxHash.Hex()
			in.BlockHash = receipt.BlockHash.Hex()
			in.Data = common.Bytes2Hex(vLog.Data)
			in.Address = vLog.Address.Hex()
			in.Removed = vLog.Removed
			if vLog.Topics != nil {
				var topiclen = len(vLog.Topics)
				topics := make([]string, topiclen)
				for i := 0; i < topiclen; i++ {
					var hash = vLog.Topics[i].Hex()
					topics[i] = hash
				}
				in.Topics = topics
			}
			in.Height = int32(vLog.BlockNumber)
			in.LogIndex = int32(vLog.Index)
			in.UpdateTime = int32(gtime.Timestamp())

			ins[num] = in
			num++
			doEventLogData(ctx, ent, vLog.Address.Hex(), int64(vLog.BlockNumber), int64(vLog.Index), int64(in.UpdateTime), vLog.Removed, txHash, vLog.Data, vLog.Topics)
		}
	}
	if num > 0 {
		err := service.TransactionLog().ExecBatchInsert(ctx, ins, chainName)
		if err != nil {
			return err
		}
	}
	return
}

// 监听event事件爬取方式，处理logs数据
func DoEventLogs(ctx context.Context, logs []types.Log, ent entity.ContractEntity, chainName string) (err error) {
	g.Log().Info(ctx, "Contract:", ent.ContractAddress, ", LogsLen:", len(logs))

	ins := make([]model.TransactionLogCreateInput, len(logs))
	num := 0
	executeNum := 1
	for _, vLog := range logs {
		if g.IsEmpty(vLog) || vLog.BlockNumber == 0 {
			continue
		}
		txHash := vLog.TxHash.Hex()
		hashLog, err := service.TransactionLog().QueryExistsByTxHashAndIndex(ctx, txHash, int32(vLog.Index), chainName)
		if err != nil {
			return err
		}

		//如果不存在，就创建
		if !hashLog {
			//vLog 和 func do里的vLog类型不一致
			var in model.TransactionLogCreateInput
			in.TxHash = txHash
			in.BlockHash = vLog.BlockHash.Hex()
			in.Data = common.Bytes2Hex(vLog.Data)
			in.Address = vLog.Address.Hex()
			in.Removed = vLog.Removed
			if vLog.Topics != nil {
				var topiclen = len(vLog.Topics)
				topics := make([]string, topiclen)
				for i := 0; i < topiclen; i++ {
					var hash = vLog.Topics[i].Hex()
					topics[i] = hash
				}
				in.Topics = topics
			}
			in.Height = int32(vLog.BlockNumber)
			in.LogIndex = int32(vLog.Index)
			in.UpdateTime = int32(gtime.Timestamp())
			ins[num] = in
			num++
			go doEventLogData(ctx, ent, vLog.Address.Hex(), int64(vLog.BlockNumber), int64(vLog.Index), int64(in.UpdateTime), vLog.Removed, txHash, vLog.Data, vLog.Topics)
		}
		executeNum++
	}

	if num > 1 {
		err := service.TransactionLog().ExecBatchInsert(ctx, ins, chainName)
		if err != nil {
			return err
		}
	}
	return
}

// 解析log的Data数据
func doEventLogData(ctx context.Context, ent entity.ContractEntity, address string, blockNumber int64, logIndex int64, blockTime int64, removed bool, txHash string, txData []byte, topics []common.Hash) (err error) {
	methodType := "evt"
	inputData := topics[0]
	abiStaking, err := abi.JSON(strings.NewReader(ent.AbiJson))
	if err != nil {
		g.Log().Error(ctx, ", event_logs_loader.doEventLogData abi.JSON error:", err, ", Please check the abi data!")
		return
	}

	methodName, err := abiStaking.EventByID(inputData)
	if err != nil {
		//g.Log().Error(ctx, ", event_logs_loader.doEventLogData abiStaking.MethodById, txHash:", txHash, " ,contractCode:", ent.ProtocolCode+"_"+ent.ContractCode, "error:", err.Error())
		return
	}

	data := inputData[4:]

	method, _ := abiStaking.Events[methodName.Name]
	receivedMap := map[string]interface{}{}
	err = method.Inputs.UnpackIntoMap(receivedMap, txData)
	if err != nil {
		g.Log().Info(ctx, "Contract:", ent.ContractAddress, ", event_logs_loader.doEventLogData Unpack deposit pubkey error:", err, ",data:", data)
		return
	}

	table := ent.ProtocolCode + "_" + ent.ContractCode + "_" + methodType + "_" + method.Name
	var insertSql bytes.Buffer
	insertSql.WriteString("insert into ")
	insertSql.WriteString(ent.ProtocolCode + "." + table)
	var insertValue bytes.Buffer

	if !g.IsEmpty(receivedMap) || len(topics) > 1 {
		insertSql.WriteString(" (")
		insertSql.WriteString(" evt_height , evt_block_time , evt_tx_hash , evt_protocol_code , evt_contract_code , evt_removed , evt_log_index , evt_from , evt_event , ")

		insertValue.WriteString(" (")
		insertValue.WriteString(strconv.FormatInt(blockNumber, 10))
		insertValue.WriteString(",")
		insertValue.WriteString(strconv.FormatInt(blockTime, 10))
		insertValue.WriteString(",")
		insertValue.WriteString("'" + txHash + "'")
		insertValue.WriteString(",")
		insertValue.WriteString("'" + ent.ProtocolCode + "'")
		insertValue.WriteString(",")
		insertValue.WriteString("'" + ent.ContractCode + "'")
		insertValue.WriteString(",")
		insertValue.WriteString(strconv.FormatBool(removed))
		insertValue.WriteString(",")
		insertValue.WriteString(strconv.FormatInt(logIndex, 10))
		insertValue.WriteString(",")
		insertValue.WriteString("'" + address + "'")
		insertValue.WriteString(",")
		insertValue.WriteString("'" + methodName.Name + "'")
		insertValue.WriteString(",")

		paramLen := len(method.Inputs)
		paramLen1 := paramLen - 1
		index := 1
		for i := 0; i < paramLen; i++ {
			m := method.Inputs[i]
			insertSql.WriteString(m.Name)
			//如果有Indexed=true, 从topics下标1开始取数据
			if m.Indexed == true {
				abiType := m.Type.String()
				v, err := utils.AdaptAbiSimpleData(abiType, topics[index])
				if err != nil {
					//#其他类型数据待处理xw
					g.Log().Debug(ctx, "event_logs_loader.doEventLogData blockNumber:", blockNumber, " ,txHash:", txHash, ",topic 未知的类型，abiType:", abiType)
					return err
				}
				insertValue.WriteString(v)
				index++
			} else {
				abiType := m.Type.String()
				v, err := utils.AdaptAbiData(abiType, receivedMap, m.Name)
				if err != nil {
					//#其他类型数据待处理
					g.Log().Debug(ctx, "event_logs_loader.doEventLogData blockNumber:", blockNumber, " ,txHash:", txHash, ",未知的类型，abiType:", abiType)
					return err
				}
				insertValue.WriteString(v)
			}

			if i < paramLen1 {
				insertSql.WriteString(",")
				insertValue.WriteString(",")
			}
		}
		insertSql.WriteString(" ) values ")
		insertValue.WriteString(")")
		insertSql.WriteString(insertValue.String())
		_, err = g.DB().Exec(ctx, insertSql.String())
		if err != nil {
			g.Log().Debug(ctx, "event_logs_loader.doEventLogData evt data insert error:", err)
			err = gerror.New("event_logs_loader.doEventLogData evt data insert error , txHash: " + txHash)
			return
		}
	}
	return
}

func Undo(ctx context.Context) {

}
