package contract

import (
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/chain/ethEvm/template"
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/utils"
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
)

type (
	sContract struct{}
)

func (s *sContract) UpdateInfo(ctx context.Context, protocolCode string, contractCode string, onceHeight string, deployHeight string, chainName string) (err error) {
	return dao.Contract.UpdateInfo(ctx, protocolCode, contractCode, onceHeight, deployHeight, chainName)
}

func (s *sContract) QueryEnableListByProtocolCode(ctx context.Context, protocolCode string, chainName string, dataType int) (rcd []entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryEnableListByProtocolCode(ctx, protocolCode, chainName, dataType)
	if err != nil {
		return nil, err
	} else {
		return adaptCurHeightWithResultList(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryEnableList(ctx context.Context, chainName string) (rcd []entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryEnableList(ctx, chainName)
	if err != nil {
		return nil, err
	} else {
		return adaptCurHeightWithResultList(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryList(ctx context.Context, protocolCode string, contractCode string, isValid string, chainName string) (rcd []entity.ContractViewEntity, err error) {
	rcd, err = dao.Contract.QueryList(ctx, protocolCode, contractCode, isValid, chainName)
	if err != nil {
		return nil, err
	} else {
		return adaptCurHeightWithResultListAndView(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryEnableListWithDataType(ctx context.Context, chainName string, dataType int) (rcd []entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryEnableListWithDataType(ctx, chainName, dataType)
	if err != nil {
		return nil, err
	} else {
		return adaptCurHeightWithResultList(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryToAuditList(ctx context.Context, chainName string) (rcd []entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryToAuditList(ctx, chainName)
	if err != nil {
		return nil, err
	} else {
		return adaptCurHeightWithResultList(ctx, rcd, chainName)
	}
}

func (s *sContract) UpdateCurHeight(ctx context.Context, contractCode string, protocolCode string, curHeight int32, chainName string) (err error) {
	rcd, err := dao.Contract.QueryOneByCode(ctx, protocolCode, contractCode, chainName)
	if err != nil {
		return err
	} else {
		//合约在缓存中的当前爬取区块高度值
		redisCurHeightValueKey := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
		if !g.IsEmpty(rcd.ContractAddress) {
			// 更新缓存
			_, err := g.Redis().Do(ctx, "SET", redisCurHeightValueKey, curHeight)
			if err != nil {
				return err
			}
		}

		// 如果当前爬取区块 >= 结束区块
		if curHeight >= rcd.RunHeight {
			rcd, err := service.ContractHistoryJob().QueryByContract(ctx, protocolCode, contractCode, chainName)
			if err == nil && g.IsEmpty(rcd.ContractCode) {
				var in model.ContractHistoryJobCreateInput
				in.IsHistoryFinish = true
				in.ProtocolCode = protocolCode
				in.ContractCode = contractCode
				in.UpdateTime = int32(gtime.Timestamp())
				err = service.ContractHistoryJob().Create(ctx, in, chainName)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (s *sContract) UpdateValid(ctx context.Context, protocolCode string, contractCode string, isVaild bool, onceHeight int, runHeight int, dataType int, chainName string) (err error) {
	return dao.Contract.UpdateValid(ctx, protocolCode, contractCode, isVaild, onceHeight, runHeight, dataType, chainName)
}

func init() {
	service.RegisterContract(New())
}

func New() *sContract {
	return &sContract{}
}

func (s *sContract) Create(ctx context.Context, in model.ContractCreateInput, chainName string) (err error) {
	return dao.Contract.ExecInsert(ctx, in, chainName)
}

func (s *sContract) DeleteByCode(ctx context.Context, protocolCode string, contractCode string, chainName string) (err error) {
	rcd, err := dao.Contract.QueryOneByCode(ctx, protocolCode, contractCode, chainName)
	if err != nil {
		return err
	}

	if g.IsEmpty(rcd.ContractAddress) {
		return gerror.New("the contractCode is not exist")
	}

	//删除合约配置数据
	_, err = dao.Contract.CtxWithDatabase(ctx, chainName).Delete("contract_code", contractCode)
	if err != nil {
		return err
	}

	//删除合约配置下的交易数据
	_, err = dao.Transaction.CtxWithDatabase(ctx, chainName).Delete("contract_address", rcd.ContractAddress)
	if err != nil {
		return err
	}

	// 合约删除后，释放锁
	lockKey := chainName + "_" + consts.LOCK_TRANSACTION_KEY + "_" + strconv.Itoa(rcd.DataType) + "_" + protocolCode + "_" + contractCode
	defer g.Redis().Do(ctx, "DEL", lockKey)

	dm := dao.ContractHistoryJob.CtxWithDatabase(ctx, chainName)
	dm.Delete(dm.Builder().Where("contract_code", contractCode).Where("protocol_code", protocolCode))
	historyjobKey := chainName + "_" + consts.PROTOCOL_CONTRACT_HISTORY_FINISH_DATA_KEY + contractCode
	_, err = g.Redis().Do(ctx, "DEL", historyjobKey)
	if err != nil {
		return err
	}

	var delTracsLogsSql bytes.Buffer
	delTracsLogsSql.WriteString("  ALTER TABLE " + chainName + ".transaction_log  delete where tx_hash  in  (  select tx_hash  from bsc.transaction  where contract_address  = ")
	delTracsLogsSql.WriteString("'")
	delTracsLogsSql.WriteString(rcd.ContractAddress)
	delTracsLogsSql.WriteString("')")

	//删除合约配置下的eventLogs数据
	_, err = g.DB().Exec(ctx, delTracsLogsSql.String())
	if err != nil {
		return err
	}

	//删除合约在缓存中的当前爬取区块高度值
	redisCurHeightValue := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
	_, err = g.Redis().Do(ctx, "DEL", redisCurHeightValue)
	if err != nil {
		return err
	}

	//删除合约多线程使用的缓存
	blockkey := chainName + "_" + protocolCode + "_" + contractCode + "_*"
	v, err := g.Redis().Do(ctx, "keys", blockkey)
	if err != nil {
		return err
	}
	if len(v.Array()) > 0 {
		for i := 0; i < len(v.Array()); i++ {
			g.Redis().Do(ctx, "DEL", v.Array()[i])
		}
	}

	//查询需要删除的表数据，并拼接成drop table sql
	res, err := service.UserAction().QueryTablesByProtocolCodeAndContractCode(ctx, protocolCode, contractCode)
	if err != nil {
		return err
	}

	//执行 drop table sql
	err = service.UserAction().Excute(ctx, res)
	if err != nil {
		return err
	}
	return nil
}

func (s *sContract) QueryByContractAddress(ctx context.Context, protocolCode string, contractAddress string, chainName string) (rcd entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryOneByContract(ctx, protocolCode, contractAddress, chainName)
	if err != nil {
		return rcd, err
	} else {
		return adaptCurHeightWithResult(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryOneByContractAddress(ctx context.Context, contractAddress string, chainName string) (rcd entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryOneByContractAddress(ctx, contractAddress, chainName)
	if err != nil {
		return rcd, err
	} else {
		return adaptCurHeightWithResult(ctx, rcd, chainName)
	}
}

func (s *sContract) QueryByContractCode(ctx context.Context, protocolCode string, contractCode string, chainName string) (rcd entity.ContractEntity, err error) {
	rcd, err = dao.Contract.QueryOneByCode(ctx, protocolCode, contractCode, chainName)
	if err != nil {
		return rcd, err
	} else {
		if !g.IsEmpty(rcd.ContractAddress) {
			// 检查任务是否已完成, 如果已完成，设置 CurrHeight = RunHeight
			hrcd, err := service.ContractHistoryJob().QueryByContract(ctx, protocolCode, contractCode, chainName)
			if err == nil && !g.IsEmpty(hrcd.ContractCode) && hrcd.IsHistoryFinish {
				rcd.CurrHeight = rcd.RunHeight
				return rcd, nil
			} else {
				return adaptCurHeightWithResult(ctx, rcd, chainName)
			}
		}
		return rcd, nil
	}
}

func (s *sContract) DoEventContractSign(ctx context.Context, contractIn model.ContractCreateInput, abiEvents map[string]abi.Event, chainName string) (events string, eventSigs string, err error) {
	var eventBuf bytes.Buffer
	var eventSigBuf bytes.Buffer
	functiontType := "event"

	for key := range abiEvents {
		hash := crypto.Keccak256Hash([]byte(abiEvents[key].Sig))
		if g.IsEmpty(hash) == false {
			eventBuf.WriteString(abiEvents[key].Sig)
			eventBuf.WriteString(", ")
			eventSigBuf.WriteString(hash.Hex())
			eventSigBuf.WriteString(", ")
		}
		var in model.ContractSignCreateInput
		in.Address = contractIn.ContractAddress
		in.Name = key
		in.Type = functiontType
		in.Sign = utils.Substring(hash.Hex(), 0, 10)
		in.SignText = abiEvents[key].Sig
		in.SignTextView = abiEvents[key].String()

		in.UpdateTime = int32(gtime.Timestamp())
		dao.ContractSign.ExecInsert(context.Background(), in, chainName)

		method := abiEvents[key]
		var inputColumn bytes.Buffer
		if g.IsEmpty(hash) == false {
			for i := 0; i < len(method.Inputs); i++ {
				argument := method.Inputs[i]
				if !g.IsEmpty(argument.Type.String()) {
					inputColumn.WriteString("`" + argument.Name + "`   ")
					inputColumn.WriteString(adpatCkDataType(argument.Type.String()))
					if i < len(method.Inputs)-1 {
						inputColumn.WriteString(", ")
					}
				}
			}
		}

		table := contractIn.ProtocolCode + "_" + contractIn.ContractCode + "_" + adaptFunctionType(functiontType) + "_" + key
		tableComment := table

		hasInputColumn := ""
		if inputColumn.Len() > 0 {
			hasInputColumn = ","
		}

		v := g.View()
		v.SetDelimiters("${", "}")
		userFunctionSql, err := v.ParseContent(context.Background(),
			template.USER_EVENT_TEMP,
			map[string]interface{}{
				"tableName":      contractIn.ProtocolCode + "." + table,
				"tableComment":   tableComment,
				"inputColumn":    inputColumn.String(),
				"hasInputColumn": hasInputColumn,
			})
		if err == nil {
			//生成用户方法表
			_, err := g.DB().Exec(ctx, userFunctionSql)
			if err != nil {
				return "", "", err
			}
		}
	}

	events = eventBuf.String()
	eventSigs = eventSigBuf.String()
	return
}

func (s *sContract) DoFunctionContractSign(ctx context.Context, contractIn model.ContractCreateInput, abiFunctions map[string]abi.Method, chainName string) (functions string, functionSigs string, err error) {
	var functionBuf bytes.Buffer
	var functionSigBuf bytes.Buffer
	functiontType := "function"

	for key := range abiFunctions {
		hash := crypto.Keccak256Hash([]byte(abiFunctions[key].Sig))

		if g.IsEmpty(hash) == false {
			functionBuf.WriteString(abiFunctions[key].Sig)
			functionBuf.WriteString(", ")
			functionSigBuf.WriteString(hash.Hex())
			functionSigBuf.WriteString(", ")
		}
		var in model.ContractSignCreateInput
		in.Address = contractIn.ContractAddress
		in.Name = key
		in.Type = functiontType
		in.Sign = utils.Substring(hash.Hex(), 0, 10)
		in.SignText = abiFunctions[key].Sig
		in.SignTextView = abiFunctions[key].String()
		in.UpdateTime = int32(gtime.Timestamp())

		//合约签名表插入数据
		go dao.ContractSign.ExecInsert(context.Background(), in, chainName)

		method := abiFunctions[key]
		var inputColumn bytes.Buffer
		if g.IsEmpty(hash) == false {
			for i := 0; i < len(method.Inputs); i++ {
				argument := method.Inputs[i]
				if !g.IsEmpty(argument.Type.String()) && !g.IsEmpty(argument.Name) {
					inputColumn.WriteString("`" + argument.Name + "`   ")
					inputColumn.WriteString(adpatCkDataType(argument.Type.String()))
					if i < len(method.Inputs)-1 {
						inputColumn.WriteString(", ")
					}
				}
			}
		}

		hasInputColumn := ""
		if inputColumn.Len() > 0 {
			hasInputColumn = ","
		}

		table := contractIn.ProtocolCode + "_" + contractIn.ContractCode + "_" + adaptFunctionType(functiontType) + "_" + key
		tableComment := table

		v := g.View()
		v.SetDelimiters("${", "}")
		userFunctionSql, err := v.ParseContent(context.Background(),
			template.USER_FUNCTION_TEMP,
			map[string]interface{}{
				"tableName":      contractIn.ProtocolCode + "." + table,
				"tableComment":   tableComment,
				"inputColumn":    inputColumn.String(),
				"hasInputColumn": hasInputColumn,
			})
		if err == nil {
			//生成用户方法表
			_, err := g.DB().Exec(ctx, userFunctionSql)
			if err != nil {
				return "", "", err
			}
		}

	}

	functions = functionBuf.String()
	functionSigs = functionSigBuf.String()
	return
}

func adaptCurHeightWithResult(ctx context.Context, rcd entity.ContractEntity, chainName string) (res entity.ContractEntity, err error) {
	//合约在缓存中的当前爬取区块高度值
	redisCurHeightValueKey := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
	if !g.IsEmpty(rcd.ContractAddress) {
		v, err := g.Redis().Do(ctx, "GET", redisCurHeightValueKey)
		if err != nil {
			return rcd, err
		}
		rcd.CurrHeight = v.Int32()
	}
	return rcd, nil
}

/*func adaptCurHeightWithResultAndView(ctx context.Context, rcd entity.ContractViewEntity, chainName string) (res entity.ContractViewEntity, err error) {
	//合约在缓存中的当前爬取区块高度值
	redisCurHeightValueKey := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
	if !g.IsEmpty(rcd.ContractAddress) {
		v, err := g.Redis().Do(ctx, "GET", redisCurHeightValueKey)
		if err != nil {
			return rcd, err
		}
		rcd.CurrHeight = v.Int32()
	}
	return rcd, nil
}*/

func adaptCurHeightWithResultList(ctx context.Context, rcds []entity.ContractEntity, chainName string) (res []entity.ContractEntity, err error) {
	res = make([]entity.ContractEntity, len(rcds))

	for i := 0; i < len(rcds); i++ {
		rcd := rcds[i]
		//合约在缓存中的当前爬取区块高度值
		redisCurHeightValueKey := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
		if !g.IsEmpty(rcd.ContractAddress) {
			v, err := g.Redis().Do(ctx, "GET", redisCurHeightValueKey)
			if err != nil {
				return nil, err
			}
			rcd.CurrHeight = v.Int32()
		}
		res[i] = rcd
	}
	return res, nil
}

func adaptCurHeightWithResultListAndView(ctx context.Context, rcds []entity.ContractViewEntity, chainName string) (res []entity.ContractViewEntity, err error) {
	res = make([]entity.ContractViewEntity, len(rcds))

	for i := 0; i < len(rcds); i++ {
		rcd := rcds[i]
		//合约在缓存中的当前爬取区块高度值
		redisCurHeightValueKey := chainName + consts.PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY + rcd.ContractAddress
		if !g.IsEmpty(rcd.ContractAddress) {
			v, err := g.Redis().Do(ctx, "GET", redisCurHeightValueKey)
			if err != nil {
				return nil, err
			}
			rcd.CurrHeight = v.Int32()
		}
		res[i] = rcd
	}
	return res, nil
}

/*
func adaptViewEntity(rcd entity.ContractEntity) entity.ContractViewEntity {
	var res entity.ContractViewEntity
	res.RunHeight = rcd.RunHeight
	res.DeployHeight = rcd.DeployHeight
	res.CurrHeight = rcd.CurrHeight
	res.ContractAddress = rcd.ContractAddress
	res.ProtocolCode = rcd.ProtocolCode
	res.ContractCode = rcd.ContractCode
	res.OnceHeight = rcd.OnceHeight
	res.IsValid = rcd.IsValid
	res.DataType = rcd.DataType
	res.UpdateTime = rcd.UpdateTime
	return res
}

func adaptViewEntityList(rcds []entity.ContractEntity) []entity.ContractViewEntity {
	resList := make([]entity.ContractViewEntity, len(rcds))
	for i := 0; i < len(rcds); i++ {
		rcd := rcds[i]
		var res entity.ContractViewEntity
		res.RunHeight = rcd.RunHeight
		res.DeployHeight = rcd.DeployHeight
		res.CurrHeight = rcd.CurrHeight
		res.ContractAddress = rcd.ContractAddress
		res.ProtocolCode = rcd.ProtocolCode
		res.ContractCode = rcd.ContractCode
		res.OnceHeight = rcd.OnceHeight
		res.IsValid = rcd.IsValid
		res.DataType = rcd.DataType
		res.UpdateTime = rcd.UpdateTime
		resList[i] = res
	}
	return resList
}*/

func adaptFunctionType(ftype string) string {
	if ftype == "function" {
		return "call"
	} else if ftype == "event" {
		return "evt"
	}
	return ""
}

func adpatCkDataType(abiType string) string {
	if abiType == "address[]" || abiType == "address" {
		return "String"
	} else if abiType == "uint256[]" || abiType == "uint160[]" || abiType == "uint128[]" || abiType == "uint112[]" || abiType == "uint32[]" || abiType == "uint16[]" || abiType == "uint8[]" || abiType == "uint[]" || abiType == "string[]" {
		return "String"
	} else if abiType == "uint256" || abiType == "uint160" {
		return "UInt256"
	} else if abiType == "uint128" || abiType == "uint112" {
		return "UInt128"
	} else if abiType == "uint64" {
		return "UInt64"
	} else if abiType == "uint32" {
		return "UInt32"
	} else if abiType == "uint16" {
		return "UInt16"
	} else if abiType == "uint8" {
		return "UInt8"
	} else if abiType == "bool" {
		return "Bool"
	} else if abiType == "bytes32" || abiType == "bytes" || abiType == "bytes4" || abiType == "string" {
		return "String"
	}
	return abiType
}
