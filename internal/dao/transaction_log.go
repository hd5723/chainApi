// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"OnchainParser/internal/dao/internal"
)

// internalTransactionLogDao is internal type for wrapping internal DAO implements.
type internalTransactionLogDao = *internal.TransactionLogDao

// transactionLogDao is the data access object for table transaction_log.
// You can define custom methods on it to extend its functionality as you wish.
type transactionLogDao struct {
	internalTransactionLogDao
}

var (
	// TransactionLog is globally public accessible object for table transaction_log operations.
	TransactionLog = transactionLogDao{
		internal.NewTransactionLogDao(),
	}
)

// Fill with you ideas below.