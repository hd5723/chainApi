package transactionLog

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/service"
	"context"
)

type (
	sTransactionLog struct{}
)

func (s *sTransactionLog) ExecInsert(ctx context.Context, in model.TransactionLogCreateInput, chainName string) (err error) {
	return dao.TransactionLog.ExecInsert(ctx, in, chainName)
}

func (s *sTransactionLog) ExecBatchInsert(ctx context.Context, ins []model.TransactionLogCreateInput, chainName string) (err error) {
	return dao.TransactionLog.ExecBatchInsert(ctx, ins, chainName)
}

func (s *sTransactionLog) QueryExistsByTxHashAndIndex(ctx context.Context, txHash string, logIndex int32, chainName string) (exists bool, err error) {
	return dao.TransactionLog.QueryExistsByTxHashAndIndex(ctx, txHash, logIndex, chainName)
}

func init() {
	service.RegisterTransactionLog(New())
}

func New() *sTransactionLog {
	return &sTransactionLog{}
}
