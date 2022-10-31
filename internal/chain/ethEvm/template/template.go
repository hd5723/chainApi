package template

var (
	USER_FUNCTION_TEMP string
	USER_EVENT_TEMP    string
)

// 初始化配置
func init() {
	// 联合索引，不作额外处理
	USER_FUNCTION_TEMP = "CREATE TABLE IF NOT EXISTS ${.tableName} (    " +
		"`call_height`     Int32  COMMENT '区块号',      " +
		"`call_block_time` Int64  COMMENT '交易时间',    " +
		"`call_tx_hash` String COMMENT '当前交易的哈希值', " +
		"`call_protocol_code`  String COMMENT '协议编码', " +
		"`call_contract_code`  String  COMMENT '合约编码', " +
		"`call_is_success`   Bool COMMENT '是否成功' ,    " +
		"`call_from`  String COMMENT 'address' ,			" +
		"`call_function`  String COMMENT '方法名' 		" +
		"${.hasInputColumn}						   " +
		"${.inputColumn}						   " +
		" ) ENGINE = MergeTree                     " +
		" ORDER BY (  call_height, call_block_time, call_protocol_code , call_contract_code) " +
		" SETTINGS index_granularity = 8192   COMMENT '${.tableComment}';"

	USER_EVENT_TEMP = "CREATE TABLE IF NOT EXISTS ${.tableName} (    " +
		"`evt_height`     Int32  COMMENT '区块号',       " +
		"`evt_block_time` Int64  COMMENT '交易时间',     " +
		"`evt_tx_hash` String COMMENT '当前交易的哈希值', " +
		"`evt_protocol_code`  String COMMENT '协议编码', " +
		"`evt_contract_code`  String  COMMENT '合约编码', " +
		"`evt_removed`   Bool COMMENT '是否已删除',       " +
		"`evt_log_index`  Int32 COMMENT '当前交易的index'," +
		"`evt_from`  String COMMENT 'address', 			" +
		"`evt_event`  String COMMENT '方法名' 			" +
		"${.hasInputColumn}						   " +
		"${.inputColumn}						    " +
		" ) ENGINE = MergeTree                      " +
		" ORDER BY (  evt_height, evt_block_time, evt_protocol_code , evt_contract_code) " +
		" SETTINGS index_granularity = 8192   COMMENT '${.tableComment}';"
}
