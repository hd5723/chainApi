// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"OnchainParser/internal/dao/internal"
)

// internalContractSignDao is internal type for wrapping internal DAO implements.
type internalContractSignDao = *internal.ContractSignDao

// contractSignDao is the data access object for table contract_sign.
// You can define custom methods on it to extend its functionality as you wish.
type contractSignDao struct {
	internalContractSignDao
}

var (
	// ContractSign is globally public accessible object for table contract_sign operations.
	ContractSign = contractSignDao{
		internal.NewContractSignDao(),
	}
)

// Fill with you ideas below.
