package consts

// 常量统一用大写下划线：LOCK_TRANSACTION_KEY
var (
	LOCK_TRANSACTION_KEY                      string //
	TRANSACTION_LOG_KEY                       string //
	TRANSACTION_KEY                           string //
	CONTRACT_VALID_LIST_KEY                   string
	CONTRACT_DATA_BLOCK_TYPE                  int
	CONTRACT_DATA_EVENTLOGS_TYPE              int
	BLOCK_CHAIN_DATA_KEY                      string
	PROTOCOL_DATA_KEY                         string
	PROTOCOL_VALID_LIST_DATA_KEY              string
	PROTOCOL_CONTRACT_DATA_KEY                string
	PROTOCOL_CONTRACT_HISTORY_FINISH_DATA_KEY string
	PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY      string
	COONTRACT_DATA_KEY                        string
	BATCH_TRANSACTION_DATA_KEY                string
	SIMPLE_TX_CACHE_KEY                       string

	TRANSACTION_LOG_DATA_KEY   string
	TRANSACTION_DATA_KEY       string
	PROTOCOL_CONTRACT_LIST_KEY string
	DELETING_CONTRACT_LIST_KEY string

	ACTIVITY_TASK_KEY string
	RPC_ERR_NUM       string
	RPC_SUCCESS_NUM   string
)

func init() {
	LOCK_TRANSACTION_KEY = "lock_transaction"

	TRANSACTION_LOG_KEY = "key_transaction_log_"
	TRANSACTION_KEY = "key_transaction_"
	CONTRACT_VALID_LIST_KEY = "contract_valid_list"
	BLOCK_CHAIN_DATA_KEY = "block_chain_data_"
	PROTOCOL_DATA_KEY = "protocol_data_"
	PROTOCOL_VALID_LIST_DATA_KEY = "protocol_valid_list"
	PROTOCOL_CONTRACT_DATA_KEY = "protocol_contract_data_"
	PROTOCOL_CONTRACT_HISTORY_FINISH_DATA_KEY = "protocol_contract_history_finish_data_"

	COONTRACT_DATA_KEY = "contract_data_"
	TRANSACTION_LOG_DATA_KEY = "key_transaction_log_data_"
	TRANSACTION_DATA_KEY = "key_transaction_data_"
	BATCH_TRANSACTION_DATA_KEY = "_batch_key_transaction_data"
	SIMPLE_TX_CACHE_KEY = "_simple_tx_cache_key"
	PROTOCOL_CONTRACT_LIST_KEY = "_protocol_contract_list"
	PROTOCOL_CONTRACT_CURRENT_HEIGHT_KEY = "_protocol_contract_"

	DELETING_CONTRACT_LIST_KEY = "_deleting_contract_list"
	ACTIVITY_TASK_KEY = "_activity_task"
	RPC_ERR_NUM = "_rpc_err_num"
	RPC_SUCCESS_NUM = "_rpc_success_num "

	CONTRACT_DATA_BLOCK_TYPE = 0
	CONTRACT_DATA_EVENTLOGS_TYPE = 1
}
