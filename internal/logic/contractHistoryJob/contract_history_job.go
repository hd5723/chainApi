package contractHistoryJob

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"context"
)

type (
	sContractHistoryJob struct{}
)

func init() {
	service.RegisterContractHistoryJob(New())
}

func New() *sContractHistoryJob {
	return &sContractHistoryJob{}
}

func (s *sContractHistoryJob) QueryByContract(ctx context.Context, protocolCode string, contractCode string, chainName string) (rcd entity.ContractHistoryJobEntity, err error) {
	return dao.ContractHistoryJob.QueryOneByContract(ctx, protocolCode, contractCode, chainName)
}

func (s *sContractHistoryJob) Create(ctx context.Context, in model.ContractHistoryJobCreateInput, chainName string) (err error) {
	return dao.ContractHistoryJob.ExecInsert(ctx, in, chainName)
}
