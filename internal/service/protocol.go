// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"context"
)

type IProtocol interface {
	QueryOneByProtocolCode(ctx context.Context, protocolCode string, chainName string) (rcd entity.ProtocolEntity, err error)
	ToAuditList(ctx context.Context, chainName string) (result []entity.ProtocolEntity, err error)
	DeleteByProtocolCode(ctx context.Context, protocolCode string, chainName string) (err error)
	QueryEnableList(ctx context.Context, chainName string) (result []entity.ProtocolEntity, err error)
	QueryList(ctx context.Context, protocol_code string, is_valid string, chainName string) (result []entity.ProtocolEntity, err error)
	ExecInsert(ctx context.Context, chainName string, in model.ProtocolCreateInput) (err error)
	UpdateValid(ctx context.Context, protocolCode string, isVaild bool, chainName string) (err error)
}

var localProtocol IProtocol

func Protocol() IProtocol {
	if localProtocol == nil {
		panic("implement not found for interface IProtocol, forgot register?")
	}
	return localProtocol
}

func RegisterProtocol(i IProtocol) {
	localProtocol = i
}
