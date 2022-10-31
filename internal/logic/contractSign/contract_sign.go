package contractSign

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"context"
	"github.com/gogf/gf/v2/database/gdb"
)

type (
	sContractSign struct{}
)

func init() {
	service.RegisterContractSign(New())
}

func New() *sContractSign {
	return &sContractSign{}
}

func (s *sContractSign) QueryOneByName(ctx context.Context, contractAddress string, function string, chainName string) (ent entity.ContractSignEntity, err error) {
	return dao.ContractSign.QueryOneByName(ctx, contractAddress, function, chainName)
}

func (s *sContractSign) QueryListByName(ctx context.Context, contractAddress string, beginLimt int, pageSize int, chainName string) (ent gdb.Result, err error) {
	return dao.ContractSign.QueryListByName(ctx, contractAddress, beginLimt, pageSize, chainName)
}

func (s *sContractSign) QueryListByNameAndType(ctx context.Context, contractAddress string, beginLimt int, pageSize int, dataType string, chainName string) (ent gdb.Result, err error) {
	return dao.ContractSign.QueryListByNameAndType(ctx, contractAddress, beginLimt, pageSize, dataType, chainName)
}

func (s *sContractSign) AdaptFunctionType(ftype string) string {
	if ftype == "function" {
		return "call"
	} else {
		return "evt"
	}
}
