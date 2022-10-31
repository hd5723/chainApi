
CREATE TABLE ${.database}.block
(
    `height` Int32 COMMENT '区块高度',
    `hash` String COMMENT '当前区块的哈希值',
    `parent_hash` String COMMENT '当前区块的上一个区块的哈希值',
    `noce` String COMMENT '用于记录在该区块的矿工做了多少次哈希才成功计算出胜出区块',
    `size` Int32 COMMENT '区块大小',
    `status` Int32 COMMENT '处理状态，1:processing 2:completed 3:orphan',
    `timestamp` Int32 COMMENT '区块的生成时间',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    PRIMARY KEY height
    ORDER BY (height,
     hash)
    SETTINGS index_granularity = 8192
    COMMENT '区块信息表';


CREATE TABLE ${.database}.block_chain
(
    `chain_name` String COMMENT '区块链名',
    `base_coin` String COMMENT 'e.g BNB',
    `chain_id` Int32 COMMENT '链编号',
    `rpc_urls` String COMMENT 'RPC配置',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    ORDER BY chain_id
    SETTINGS index_granularity = 8192
    COMMENT '区块链信息表';


CREATE TABLE ${.database}.contract
(
    `protocol_code` String COMMENT '协议编码',
    `contract_code` String COMMENT '合约编码，用于分库作为数据库名',
    `contract_address` String COMMENT '合约地址',
    `abi_json` String COMMENT 'abi的JSON数据',
    `deploy_height` Int32 COMMENT '起始区块号',
    `run_height` Int32 COMMENT '结束区块号',
    `curr_height` Int32 COMMENT '当前爬取到的区块号',
    `once_height` Int32 COMMENT '一次爬取到的区块',
    `listener_event` String COMMENT '监听的事件（transfer , approve等）',
    `signed_listener_event` String COMMENT '签名好的监听事件数据',
    `listener_function` String COMMENT '监听的方法（swap , dispreseEther等）',
    `signed_listener_function` String COMMENT '签名好的监听方法数据',
    `is_valid` Bool COMMENT '是否已启用',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）',
    `data_type` Int32 DEFAULT 0 COMMENT '数据解析类型，0:通过区块解析 1:通过event事件解析'
)
    ENGINE = MergeTree
    PRIMARY KEY contract_address
    ORDER BY (contract_address,
    contract_code,
    protocol_code)
    SETTINGS index_granularity = 8192
    COMMENT '合约表';


CREATE TABLE ${.database}.contract_history_job
(
    `protocol_code` String COMMENT '协议编码',
    `contract_code` String COMMENT '合约编码，用于分库作为数据库名',
    `is_history_finish` Bool COMMENT '历史记录已爬取完',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    ORDER BY (protocol_code,
    contract_code,
    is_history_finish)
    SETTINGS index_granularity = 8192
    COMMENT '合约历史区块爬取任务';



CREATE TABLE ${.database}.contract_sign
(
    `id` Int64 COMMENT '主键ID',
    `address` String COMMENT '合约地址',
    `name` String COMMENT '方法/Event名称',
    `sign` String COMMENT '签名原始数据',
    `sign_text` String COMMENT '签名数据',
    `sign_text_view` String COMMENT '类型，function、event',
    `type` String COMMENT '类型，function、event',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    PRIMARY KEY id
    ORDER BY (id,
    address,
    name,
    sign)
    SETTINGS index_granularity = 8192
    COMMENT '合约签名表';


CREATE TABLE ${.database}.protocol
(
    `protocol_code` String COMMENT '编码，如pancakeswap',
    `protocol_name` String COMMENT '名称，如PancakeSwap',
    `is_valid` Bool COMMENT '是否有效,1:有效 0:无效',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    PRIMARY KEY protocol_code
    ORDER BY protocol_code
    SETTINGS index_granularity = 8192
    COMMENT '协议表';



CREATE TABLE ${.database}.transaction
(
    `tx_hash` String COMMENT '当前交易的哈希值',
    `tx_index` Int32 COMMENT 'index',
    `height` Int32 COMMENT '区块号',
    `contract_address` String COMMENT '合约地址',
    `status` Int32 COMMENT '交易状态 1:成功 0:失败',
    `from` String COMMENT '交易发起者的地址',
    `to` String COMMENT '交易接收者的地址',
    `value` Int64 COMMENT '交易金额',
    `gas_used` Int64 COMMENT '燃料费',
    `gas_limit` Int64 COMMENT 'the gas limit in jager',
    `gas_price` Int64 COMMENT 'the gas price in jager',
    `transaction_fee` Int64 COMMENT 'Amount paid to the miner for processing the transaction.',
    `create_time` Int32 COMMENT '创建时间（时间戳，秒级）',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）',
    `timestamp` Int32 COMMENT '交易时间'
)
    ENGINE = MergeTree
    PRIMARY KEY tx_hash
    ORDER BY (tx_hash,
    height,
    from)
    SETTINGS index_granularity = 8192
    COMMENT '交易信息表';


CREATE TABLE ${.database}.transaction_log
(
    `removed` Bool COMMENT '是否已移除',
    `log_index` Int32 COMMENT 'index',
    `tx_hash` String COMMENT '当前交易的哈希值',
    `block_hash` String COMMENT '当前区块的哈希值',
    `height` Int32 COMMENT '当前区块号',
    `address` String COMMENT '合约地址',
    `data` String COMMENT 'log data',
    `type` String COMMENT 'type',
    `topics` String COMMENT 'topics，Array<String>类型',
    `update_time` Int32 COMMENT '修改时间（时间戳，秒级）'
)
    ENGINE = MergeTree
    ORDER BY (tx_hash,
     height)
    SETTINGS index_granularity = 8192
    COMMENT '交易日志表';