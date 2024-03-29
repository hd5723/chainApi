// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

type IUserAction interface {
	QueryTablesByProtocolCodeAndContractCode(ctx context.Context, protocolCode string, contractCode string) (result []gdb.Value, err error)
	Excute(ctx context.Context, sqls []gdb.Value) (err error)
}

var localUserAction IUserAction

func UserAction() IUserAction {
	if localUserAction == nil {
		panic("implement not found for interface IUserAction, forgot register?")
	}
	return localUserAction
}

func RegisterUserAction(i IUserAction) {
	localUserAction = i
}
