package dao

import (
	"OnchainParser/internal/dao/internal"
)

// internalContractDao is internal type for wrapping internal DAO implements.
type internalContractHistoryJobDao = *internal.ContractHistoryJobDao

// contractHistoryJobDao is the data access object for table contract.
// You can define custom methods on it to extend its functionality as you wish.
type contractHistoryJobDao struct {
	internalContractHistoryJobDao
}

var (
	// contractHistoryJobDao is globally public accessible object for table contract operations.
	ContractHistoryJob = contractHistoryJobDao{
		internal.NewContractHistoryJobDao(),
	}
)

// Fill with you ideas below.
