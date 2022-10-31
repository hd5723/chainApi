package transaction

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"context"
)

type (
	sTransaction struct{}
)

func (s *sTransaction) ExecInsert(ctx context.Context, in model.TransactionCreateInput, chainName string) (err error) {
	return dao.Transaction.ExecInsert(ctx, in, chainName)
}

func (s *sTransaction) UpdateGasprice(ctx context.Context, txHash string, gasPrice int64, chainName string) (err error) {
	return dao.Transaction.UpdateGasprice(ctx, txHash, gasPrice, chainName)
}

func (s *sTransaction) QueryLteListByGasPrice(ctx context.Context, chainName string) (rcd []entity.TransactionEntity, err error) {
	return dao.Transaction.QueryLteListByGasPrice(ctx, chainName)
}

func (s *sTransaction) QueryExistsByTxHash(ctx context.Context, txHash string, chainName string) (exists bool, err error) {
	return dao.Transaction.QueryExistsByTxHash(ctx, txHash, chainName)
}

func (s *sTransaction) QueryOneByTxHash(ctx context.Context, txHash string, chainName string) (rcd entity.TransactionEntity, err error) {
	return dao.Transaction.QueryOneByTxHash(ctx, txHash, chainName)
}

func (s *sTransaction) QueryListByTime(ctx context.Context, contractAddress string, blockBeginTime string, blockEndTime string, address string, pageSize int, pageNumber int, chainName string) (result []entity.TransactionEntity, err error) {
	return dao.Transaction.QueryListByTime(ctx, contractAddress, blockBeginTime, blockEndTime, address, pageSize, pageNumber, chainName)
}

func (s *sTransaction) QueryListByAddress(ctx context.Context, contractAddress string, pageSize int, pageNumber int, chainName string) (result []entity.TransactionEntity, err error) {
	return dao.Transaction.QueryListByAddress(ctx, contractAddress, pageSize, pageNumber, chainName)
}

func (s *sTransaction) QueryOneByAddressAndTxHash(ctx context.Context, contractAddress string, txHash string, chainName string) (result entity.TransactionEntity, err error) {
	return dao.Transaction.QueryOneByAddressAndTxHash(ctx, contractAddress, txHash, chainName)
}

func (s *sTransaction) QueryByAddressAndBlockNumber(ctx context.Context, contractAddress string, blockNumer int, pageSize int, pageNumber int, chainName string) (result []entity.TransactionEntity, err error) {
	return dao.Transaction.QueryByAddressAndBlockNumber(ctx, contractAddress, blockNumer, pageSize, pageNumber, chainName)
}

func init() {
	service.RegisterTransaction(New())
}

func New() *sTransaction {
	return &sTransaction{}
}
