
-- 统计查询 某个合约地址的各区块交易数量、logs数量
SELECT t.height , t.contract_address , t.trans_num , tl.log_um
FROM
    (
    SELECT count(t.tx_hash) trans_num , t.height height , MAX(t.contract_address)  contract_address
        FROM `transaction` t
        WHERE
            t.contract_address  = '0x10ED43C718714eb63d5aA57B78B54704E256024E'
        GROUP BY  t.height
     ) AS t
LEFT JOIN
    (
        SELECT count(tl.tx_hash) log_um , tl.height
        FROM  `transaction_log` tl
        WHERE
             tl.tx_hash  in (  select tx_hash  from `transaction`  where contract_address  = '0x10ED43C718714eb63d5aA57B78B54704E256024E' )
        GROUP BY tl.height
    ) AS tl
ON t.height = tl.height
ORDER BY t.height desc ;


