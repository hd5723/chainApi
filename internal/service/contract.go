// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type IContract interface {
	UpdateInfo(ctx context.Context, protocolCode string, contractCode string, onceHeight string, deployHeight string, chainName string) (err error)
	QueryEnableListByProtocolCode(ctx context.Context, protocolCode string, chainName string, dataType int) (rcd []entity.ContractEntity, err error)
	QueryEnableList(ctx context.Context, chainName string) (rcd []entity.ContractEntity, err error)
	QueryList(ctx context.Context, protocolCode string, contractCode string, isValid string, chainName string) (rcd []entity.ContractViewEntity, err error)
	QueryEnableListWithDataType(ctx context.Context, chainName string, dataType int) (rcd []entity.ContractEntity, err error)
	QueryToAuditList(ctx context.Context, chainName string) (rcd []entity.ContractEntity, err error)
	UpdateCurHeight(ctx context.Context, contractCode string, protocolCode string, curHeight int32, chainName string) (err error)
	UpdateValid(ctx context.Context, protocolCode string, contractCode string, isVaild bool, onceHeight int, runHeight int, dataType int, chainName string) (err error)
	Create(ctx context.Context, in model.ContractCreateInput, chainName string) (err error)
	DeleteByCode(ctx context.Context, protocolCode string, contractCode string, chainName string) (err error)
	QueryByContractAddress(ctx context.Context, protocolCode string, contractAddress string, chainName string) (rcd entity.ContractEntity, err error)
	QueryOneByContractAddress(ctx context.Context, contractAddress string, chainName string) (rcd entity.ContractEntity, err error)
	QueryByContractCode(ctx context.Context, protocolCode string, contractCode string, chainName string) (rcd entity.ContractEntity, err error)
	DoEventContractSign(ctx context.Context, contractIn model.ContractCreateInput, abiEvents map[string]abi.Event, chainName string) (events string, eventSigs string, err error)
	DoFunctionContractSign(ctx context.Context, contractIn model.ContractCreateInput, abiFunctions map[string]abi.Method, chainName string) (functions string, functionSigs string, err error)
}

var localContract IContract

func Contract() IContract {
	if localContract == nil {
		panic("implement not found for interface IContract, forgot register?")
	}
	return localContract
}

func RegisterContract(i IContract) {
	localContract = i
}
