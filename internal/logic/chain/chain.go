package contract

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"context"
)

type (
	sChain struct{}
)

func init() {
	service.RegisterChain(New())
}

func New() *sChain {
	return &sChain{}
}

func (s *sChain) QueryOneByChainId(ctx context.Context, chianId int32, chainName string) (rcd entity.BlockChainEntity) {
	return dao.BlockChain.QueryOneByChainId(ctx, chianId, chainName)
}

func (s *sChain) ExecInsert(ctx context.Context, chainName string, in model.BlockChainCreateInput) error {
	return dao.BlockChain.ExecInsert(ctx, chainName, in)
}

func (s *sChain) Update(ctx context.Context, chainName string, in model.BlockChainCreateInput) error {
	return dao.BlockChain.Update(ctx, chainName, in)
}

func (s *sChain) QueryOneByChainName(ctx context.Context, chainName string) (rcd entity.BlockChainEntity) {
	return dao.BlockChain.QueryOneByChainName(ctx, chainName)
}
