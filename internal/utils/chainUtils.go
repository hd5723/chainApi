package utils

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"math/big"
	"reflect"
	"strconv"
)

func AdaptFunctionType(ftype int) string {
	if ftype == 3 {
		return "call"
	} else {
		return "evt"
	}
}

func AdaptBigIntResult(temp gdb.Result) []gmap.Map {
	if temp == nil || len(temp) == 0 {
		return nil
	}

	resultList := make([]gmap.Map, len(temp))
	for i := 0; i < len(temp); i++ {
		resultList[i] = AdaptBigIntRecord(temp[i])
	}
	return resultList
}

func AdaptBigIntRecord(record gdb.Record) gmap.Map {
	var gRecord gmap.Map
	for key, val := range record.Map() {
		if "big.Int" == reflect.TypeOf(val).String() {
			v := record[key].Interface().(big.Int)
			gRecord.SetIfNotExist(key, v.String())
		} else {
			gRecord.SetIfNotExist(key, val)
		}
	}
	return gRecord
}

func AdaptAbiSimpleData(abiType string, value common.Hash) (res string, err error) {
	if abiType == "address" {
		res = "'" + value.Hex() + "'"
	} else if abiType == "string" {
		res = "'" + value.String() + "'"
	} else if abiType == "uint256" || abiType == "uint160" || abiType == "uint128" || abiType == "uint112" || abiType == "uint32" || abiType == "uint16" || abiType == "uint" {
		res = value.Hex()
	} else if abiType == "bool" {
		res = value.Hex()
	} else if abiType == "uint8" {
		res = value.Hex()
	} else if abiType == "bytes" {
		res = value.Hex()
	} else if abiType == "bytes4" {
		//v := value.Bytes() //Uint8
		//var bytes4Value bytes.Buffer
		//for i := 0; i < len(v); i++ {
		//	bytes4Value.WriteString(strconv.Itoa(int(v[i])))
		//}
		//res = "'" + common.HexToAddress(bytes4Value.String()).Hex() + "'"
		res = value.Hex()
	} else if abiType == "bytes32" {
		//v := value.Bytes() //Uint8
		//var bytes32Value bytes.Buffer
		//for i := 0; i < len(v); i++ {
		//	bytes32Value.WriteString(strconv.Itoa(int(v[i])))
		//}
		//res = "'" + common.HexToAddress(bytes32Value.String()).Hex() + "'"
		res = value.Hex()
	} else {
		//# 其他类型数据待处理
		err = gerror.New("未知的类型，abiType: " + abiType)
		return
	}
	return
}

func AdaptAbiData(abiType string, receivedMap map[string]interface{}, methodName string) (res string, err error) {
	if abiType == "address[]" {
		addressHexValue := ""
		addressVals := receivedMap[methodName].([]common.Address)
		for i := 0; i < len(addressVals); i++ {
			addressHexValue += addressVals[i].Hex()
			if i < len(addressVals)-1 {
				addressHexValue += ","
			}
		}
		res = "'" + addressHexValue + "'"
	} else if abiType == "address" {
		res = "'" + (receivedMap[methodName].(common.Address).Hex()) + "'"
	} else if abiType == "string[]" {
		addressValue := ""
		addressVals := receivedMap[methodName].([]string)
		for i := 0; i < len(addressVals); i++ {
			addressValue += addressVals[i]
			if i < len(addressVals)-1 {
				addressValue += ","
			}
		}
		res = "'" + addressValue + "'"
	} else if abiType == "string" {
		res = "'" + receivedMap[methodName].(string) + "'"
	} else if abiType == "uint256[]" || abiType == "uint160[]" || abiType == "uint128[]" || abiType == "uint112[]" || abiType == "uint32[]" || abiType == "uint16[]" || abiType == "uint[]" {
		var value bytes.Buffer
		vals := receivedMap[methodName].([]*big.Int)
		for i := 0; i < len(vals); i++ {
			value.WriteString(vals[i].String())
			if i < len(vals)-1 {
				value.WriteString(",")
			}
		}
		res = "'" + value.String() + "'"
	} else if abiType == "uint8[]" {
		var value bytes.Buffer
		vals := receivedMap[methodName].([]uint8)
		for i := 0; i < len(vals); i++ {
			value.WriteString(strconv.Itoa(int(vals[i])))
			if i < len(vals)-1 {
				value.WriteString(",")
			}
		}
		res = "'" + value.String() + "'"
	} else if abiType == "uint256" || abiType == "uint160" || abiType == "uint128" || abiType == "uint112" || abiType == "uint32" || abiType == "uint16" || abiType == "uint" {
		res = receivedMap[methodName].(*big.Int).String()
		//res = fmt.Printf("%s\n", receivedMap[methodName])
	} else if abiType == "bool" {
		res = strconv.FormatBool(receivedMap[methodName].(bool))
	} else if abiType == "uint8" {
		res = strconv.Itoa(int(receivedMap[methodName].(uint8)))
	} else if abiType == "bytes" {
		v := receivedMap[methodName].([]uint8)
		b := make([]byte, 1)
		b[0] = v[0]
		res = "'" + hexutil.Encode(b) + "'"
	} else if abiType == "bytes4" {
		v := receivedMap[methodName].([4]uint8)
		b := make([]byte, 4)
		for i := 0; i < len(v); i++ {
			b[i] = v[i]
		}
		res = "'" + hexutil.Encode(b) + "'"
	} else if abiType == "bytes32" {
		v := receivedMap[methodName].([32]uint8)
		b := make([]byte, 32)
		for i := 0; i < len(v); i++ {
			b[i] = v[i]
		}
		res = "'" + hexutil.Encode(b) + "'"
	} else {
		//# 其他类型数据待处理
		err = gerror.New("未知的类型，abiType: " + abiType)
		return
	}
	return
}
