package transactionTest

import (
	"OnchainParser/internal/chain/ethEvm/chain/baseConfig"
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"bytes"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"time"
)

var (
	ctx             = gctx.New()
	tempPath        = gfile.Pwd() + "/sql/temp.sql"
	batchInsertPath = gfile.Pwd() + "/sql/batchInsert.sql"
)

func BatchInsert() {
	baseConfig.CHAIN_NAME = "bsc"
	chainName := baseConfig.CHAIN_NAME

	if gfile.Exists(tempPath) {
		gfile.Create(tempPath)
	}

	for i := 0; i < 50000; i++ {
		var in model.TransactionCreateInput
		in.TxHash = "ox89435y8hf00055" + strconv.Itoa(i)
		in.TxIndex = 1
		in.Height = 1
		in.ContractAddress = "0xA0000000000000000000000000000"
		in.From = "0xA0000000000000000000000000000"
		in.To = "0xA0000000000000000000000000000"
		in.Status = int32(1)
		in.GasUsed = int64(200000000)
		in.GasPrice = int64(200000000)
		in.GasLimit = int64(6000)
		in.TransactionFee = int64(200000000)
		in.Value = int64(500000000)
		in.CreateTime = int32(gtime.Timestamp())
		in.UpdateTime = int32(gtime.Timestamp())
		in.Timestamp = int32(gtime.Timestamp())

		var dimSqlDML bytes.Buffer
		dimSqlDML.WriteString("(")
		dimSqlDML.WriteString("'" + in.TxHash + "'")
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.TxIndex), 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.Height), 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString("'" + in.ContractAddress + "'")
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.Status), 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString("'" + in.From + "'")
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString("'" + in.To + "'")
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(in.Value, 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(in.GasUsed, 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(in.GasLimit, 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(in.GasPrice, 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(in.TransactionFee, 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.CreateTime), 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.UpdateTime), 10))
		dimSqlDML.WriteString(",")
		dimSqlDML.WriteString(strconv.FormatInt(int64(in.Timestamp), 10))
		dimSqlDML.WriteString(")")
		gfile.PutContentsAppend(tempPath, dimSqlDML.String()+"$;")
	}

	err := gfile.Move(tempPath, batchInsertPath)
	if err != nil {
		g.Log().Error(ctx, err.Error())
	}

	curTime := time.Now()
	sqlContent := strings.ReplaceAll(gfile.GetContents(batchInsertPath), "\n", "")
	if g.IsEmpty(sqlContent) {
		return
	}
	gfile.Remove(batchInsertPath)
	sqlArray := strings.Split(sqlContent, "$;")

	var sqlDML bytes.Buffer
	sqlDML.WriteString("insert into ")
	sqlDML.WriteString(chainName + "." + dao.Transaction.Table())
	sqlDML.WriteString(" (tx_hash, tx_index,height,contract_address,status,from,to,value,gas_used,gas_limit,gas_price,transaction_fee,create_time,update_time,timestamp) values  ")

	for i := 0; i < len(sqlArray); i++ {
		sql := sqlArray[i]
		if g.IsEmpty(sql) {
			continue
		}

		sqlDML.WriteString(sql)
		if i != len(sqlArray) {
			sqlDML.WriteString(",")
		}
	}
	_, err = g.DB().Exec(ctx, sqlDML.String())
	if err != nil {
		g.Log().Error(ctx, err.Error())
	}
	g.Log().Info(ctx, " doing time:", time.Now().Sub(curTime))

}
