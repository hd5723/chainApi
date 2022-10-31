package txLoader

import (
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/chain/ethEvm/utils/loader/eventLogsLoader"
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/utils"
	"bytes"
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"time"
)

// 扫描区块爬取方式，处理tx数据
func Do(ctx context.Context, ethClient *ethclient.Client, link string, tx *types.Transaction, block *types.Block, ent entity.ContractEntity, chainName string, task chan consts.TxExecuteInfo, execType int) {

	txHash := tx.Hash().Hex()
	var result consts.TxExecuteInfo
	result.TxHash = txHash

	var blockKey string
	if execType == 1 {
		blockKey = "ALL"
	} else {
		blockKey = ent.ProtocolCode + "_" + ent.ContractCode
	}
	//检查上一次是否执行成功，上一次未成功的继续执行；成功的标记err为nil返回
	blockTxsKey := chainName + "_" + blockKey + "_" + strconv.FormatInt(int64(block.NumberU64()), 10)
	res, err := g.Redis().Do(ctx, "HGET", blockTxsKey, txHash)
	if err != nil {
		result.Err = err
		task <- result
		return
	}

	// res == true, 已经处理过，忽略
	if res != nil && res.Bool() {
		result.Err = nil
		task <- result
		return
	}

	hexTxHash := common.HexToHash(txHash)

	curTime := time.Now()
	receipt, err := ethClient.TransactionReceipt(ctx, hexTxHash)
	g.Log().Info(ctx, "RPC tx.Do.ethClient.TransactionReceipt doing time:", time.Now().Sub(curTime), " ,link:", link)
	if err != nil {
		g.Log().Info(ctx, "RPC tx.Do.ethClient.TransactionReceipt hexTxHash ", hexTxHash, " , error:", err)
		result.Err = err
		task <- result
		return
	}

	//查询是否存在此交易信息，如果不存在就创建
	queryExistsByTxHashTime := time.Now()
	existsInfo, err := dao.Transaction.QueryExistsByTxHash(ctx, txHash, chainName)
	g.Log().Info(ctx, "db tx.Do.dao.Transaction.QueryExistsByTxHash doing time:", time.Now().Sub(queryExistsByTxHashTime))
	if err != nil {
		result.Err = err
		task <- result
		return
	}

	if !existsInfo {
		if tx == nil {
			curTime := time.Now()
			trans, _, err := ethClient.TransactionByHash(ctx, hexTxHash)
			g.Log().Info(ctx, "RPC tx.Do.ethClient.TransactionByHash doing time:", time.Now().Sub(curTime), " ,link:", link)
			if err != nil {
				result.Err = err
				task <- result
				return
			}
			tx = trans
		}

		var from = ""
		curTime := time.Now()
		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
		g.Log().Info(ctx, "RPC tx.Do.tx.AsMessage doing time:", time.Now().Sub(curTime))
		if err != nil {
			result.Err = err
			task <- result
			return
		}
		from = msg.From().Hex()

		//库里保存原始值，需要时再除以精度
		var transactionFee = int64(receipt.GasUsed) * tx.GasPrice().Int64()
		// var transactionFee = float64(int64(receipt.GasUsed) * trans.GasPrice().Int64())/ math.Pow(10, 18)
		var in model.TransactionCreateInput
		in.TxHash = txHash
		in.TxIndex = int32(receipt.TransactionIndex)
		in.Height = int32(receipt.BlockNumber.Int64())
		in.ContractAddress = ent.ContractAddress
		in.From = from
		if !g.IsNil(tx.To()) {
			in.To = tx.To().Hex()
		} else {
			g.Log().Warning(ctx, "to is nil,  txHash:", txHash)
		}
		in.Status = int32(receipt.Status)
		in.GasUsed = int64(receipt.GasUsed)
		in.GasPrice = tx.GasPrice().Int64()
		in.GasLimit = int64(tx.Gas())
		in.TransactionFee = transactionFee
		in.Value = tx.Value().Int64()
		in.CreateTime = int32(gtime.Timestamp())
		in.UpdateTime = int32(gtime.Timestamp())
		in.Timestamp = int32(int64(block.Time()))

		err = dao.Transaction.ExecInsertSql(ctx, in, chainName)
		if err != nil {
			result.Err = err
			task <- result
			return
		}

		g.Log().Debug(ctx, "txLoader Do txHash:", txHash)
		go doTxData(ctx, ent, from, receipt.BlockNumber.Int64(), int64(block.Time()), int64(receipt.Status), txHash, tx.Data())
	}

	err = eventLogsLoader.Do(ctx, receipt, txHash, ent, chainName)
	if err != nil {
		result.Err = err
		task <- result
		return
	}

	result.Err = nil
	task <- result
	return
}

// 监听事件方式，处理tx数据
func DoTx(ctx context.Context, ethClient *ethclient.Client, link string, newTx consts.TxInfo, toBlock int64, ent entity.ContractEntity, chainName string, task chan consts.TxExecuteInfo) {

	txHash := newTx.TxHash
	hexTxHash := common.HexToHash(txHash)

	var result consts.TxExecuteInfo
	result.TxHash = txHash
	blockKey := ent.ProtocolCode + "_" + ent.ContractCode

	//检查上一次是否执行成功，上一次未成功的继续执行；成功的标记err为nil返回
	blockTxsKey := chainName + "_" + blockKey + "_" + strconv.FormatInt(toBlock, 10)
	res, err := g.Redis().Do(ctx, "HGET", blockTxsKey, txHash)
	if err != nil {
		result.Err = err
		task <- result
		return
	}

	// res == true, 已经处理过，忽略
	if res != nil && res.Bool() {
		result.Err = nil
		task <- result
		return
	}

	curTime := time.Now()
	receipt, err := ethClient.TransactionReceipt(ctx, hexTxHash)
	g.Log().Info(ctx, "RPC txLoader.DoTx.ethClient.TransactionReceipt doing time:", time.Now().Sub(curTime), " , link:", link, " , contract_code:", ent.ContractCode)
	if err != nil {
		result.Err = nil
		task <- result
		return
	}

	//查询是否存在此交易信息，如果不存在就创建
	curTime = time.Now()
	existsInfo, err := dao.Transaction.QueryExistsByTxHash(ctx, txHash, chainName)
	g.Log().Info(ctx, "db txLoader.DoTx.Transaction.QueryExistsByTxHash doing time:", time.Now().Sub(curTime))
	if err != nil {
		result.Err = err
		task <- result
		return
	}

	if !existsInfo {
		var from = ""
		curTime = time.Now()
		block, err := ethClient.BlockByNumber(ctx, receipt.BlockNumber)
		g.Log().Info(ctx, "RPC txLoader.DoTx.ethClient.BlockByNumber doing time:", time.Now().Sub(curTime), " , link:", link)
		if err != nil {
			result.Err = err
			task <- result
			return
		}

		curTime = time.Now()
		tx, _, err := ethClient.TransactionByHash(ctx, hexTxHash)
		g.Log().Info(ctx, "RPC txLoader.DoTx.ethClient.TransactionByHash doing time:", time.Now().Sub(curTime), " , link:", link)
		if err != nil {
			result.Err = err
			task <- result
			return
		}

		curTime = time.Now()
		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
		g.Log().Info(ctx, " txLoader.DoTx.tx.AsMessage doing time:", time.Now().Sub(curTime))
		if err != nil {
			result.Err = err
			task <- result
			return
		}

		from = msg.From().Hex()
		//库里保存原始值，需要时再除以精度
		var transactionFee = int64(receipt.GasUsed) * tx.GasPrice().Int64()
		// var transactionFee = float64(int64(receipt.GasUsed) * trans.GasPrice().Int64())/ math.Pow(10, 18)
		var in model.TransactionCreateInput
		in.TxHash = txHash
		in.TxIndex = int32(receipt.TransactionIndex)
		in.Height = int32(receipt.BlockNumber.Int64())
		in.ContractAddress = ent.ContractAddress
		in.From = from
		if !g.IsNil(tx.To()) {
			in.To = tx.To().Hex()
		} else {
			g.Log().Warning(ctx, "to is nil,  txHash:", txHash)
		}
		in.Status = int32(receipt.Status)
		in.GasUsed = int64(receipt.GasUsed)
		in.GasPrice = tx.GasPrice().Int64()
		in.GasLimit = int64(tx.Gas())
		in.TransactionFee = transactionFee
		in.Value = tx.Value().Int64()
		in.CreateTime = int32(gtime.Timestamp())
		in.UpdateTime = int32(gtime.Timestamp())
		in.Timestamp = int32(int64(block.Time()))
		err = dao.Transaction.ExecInsertSql(ctx, in, chainName)
		if err != nil {
			result.Err = err
			task <- result
			return
		}
		go doTxData(ctx, ent, from, receipt.BlockNumber.Int64(), int64(block.Time()), int64(receipt.Status), txHash, tx.Data())
	}

	result.Err = nil
	task <- result
	return
}

func Undo(ctx context.Context) {
}

// 解析tx的Data数据
func doTxData(ctx context.Context, ent entity.ContractEntity, from string, blockNumber int64, blockTime int64, status int64, txHash string, txData []byte) (err error) {
	trasData := common.Bytes2Hex(txData)
	inputData, err := hex.DecodeString(trasData)
	if err != nil {
		g.Log().Error(ctx, "Contract:", ent.ContractAddress, ", txLoader.doTxData hex.DecodeString error:", err)
		return
	}
	abiStaking, err := abi.JSON(strings.NewReader(ent.AbiJson))
	if err != nil {
		g.Log().Error(ctx, ", txLoader.doTxData abi.JSON error:", err, ", Please check the abi data!")
		return
	}

	methodName, err := abiStaking.MethodById(inputData)
	if err != nil {
		//g.Log().Error(ctx, ", txLoader.doTxData abiStaking.MethodById, txHash:", txHash, " ,contractCode:", ent.ProtocolCode+"_"+ent.ContractCode, "error:", err.Error())
		return
	}

	var data []byte
	if inputData != nil && len(inputData) > 0 {
		data = inputData[4:]
	}

	isSuccess := false
	if status == 1 {
		isSuccess = true
	}
	method, _ := abiStaking.Methods[methodName.Name]
	receivedMap := map[string]interface{}{}
	err = method.Inputs.UnpackIntoMap(receivedMap, data)
	if err != nil {
		g.Log().Error(ctx, "Contract:", ent.ContractAddress, ", txLoader.doTxData Unpack deposit pubkey error:", err)
		return
	}

	table := ent.ProtocolCode + "_" + ent.ContractCode + "_" + utils.AdaptFunctionType(int(method.Type)) + "_" + method.Name
	var insertSql bytes.Buffer
	insertSql.WriteString("insert into ")
	insertSql.WriteString(ent.ProtocolCode + "." + table)
	var insertValue bytes.Buffer

	insertSql.WriteString(" (")
	insertSql.WriteString(" call_height , call_block_time , call_tx_hash , call_protocol_code , call_contract_code , call_is_success , call_from , call_function  ")

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
	insertValue.WriteString(strconv.FormatBool(isSuccess))
	insertValue.WriteString(",")
	insertValue.WriteString("'" + from + "'")
	insertValue.WriteString(",")
	insertValue.WriteString("'" + methodName.Name + "'")

	if !g.IsEmpty(receivedMap) {
		insertSql.WriteString(",")
		insertValue.WriteString(",")

		for i := 0; i < len(method.Inputs); i++ {
			m := method.Inputs[i]
			abiType := m.Type.String()
			insertSql.WriteString(m.Name)

			v, err := utils.AdaptAbiData(abiType, receivedMap, m.Name)
			if err != nil {
				//#TODO 其他类型数据待处理
				g.Log().Debug(ctx, "txLoader.doTxData  blockNumber:", blockNumber, " ,txHash:", txHash, ",未知的类型，abiType:", abiType)
				return err
			}
			insertValue.WriteString(v)

			if i < len(method.Inputs)-1 {
				insertSql.WriteString(",")
				insertValue.WriteString(",")
			}
		}
	}
	insertSql.WriteString(" ) values ")
	insertValue.WriteString(")")
	insertSql.WriteString(insertValue.String())
	_, err = g.DB().Exec(ctx, insertSql.String())
	if err != nil {
		g.Log().Debug(ctx, "txLoader.doTxData tx call data insert error:", err)
		return err
	}

	return
}
