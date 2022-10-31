# 链上数据解析【On-chain Data Parser】

## 项目版本
    Clickhouse Version:     22.1.3.7         [SELECT  VERSION()] 
    Go Version:             go1.17.11        [go version]
    GF Version:             v2.1.2           [gf version]

## 部署项目：
1、安装go、clickhouse、redis，可参照【项目版本】

2、安装项目框架工具(gf)
    手动编译安装: git clone https://github.com/gogf/gf && cd gf/cmd/gf && go install

3、启动clickhouse、redis

4、服务器拉取项目代码：
    OnchainParser: git clone https://github.com/Port3-Network/OnchainParser.git

5、修改配置文件
    OnchainParser > manifest > config > config.xml
    可修改服务端口、日志、数据库连接、redis连接

6、clickhouse数据库初始化步骤：  

    6.1 在数据库里执行: [CREATE DATABASE IF NOT EXISTS bsc], 初始化bsc公链数据结构；

    6.2 在数据库里执行 base-[最后日期].sql，初始化数据结构

    6.3 配置其他公链，替换6.1中的bsc，替换6.2中的block_chain初始化数据

7、启动项目
    nohup gf run main.go >[log path] &
    
8、区块节点信息

|   节点  |   类型  |   来源  |
|---  |--- | --- |
|  wss://bsc-mainnet.s.chainbase.online/v1/2FHzQi5jg5UH5SJZC6Njj9x5fFr   |  全节点   |  console.chainbase.online   |
|  wss://bsc-mainnet.s.chainbase.online/v1/2FORUqTpX1jkGExMIXanXv5vp40   |  全节点   |  console.chainbase.online   |
|  wss://bsc.getblock.io/mainnet/?api_key=ff78be83-cec1-4f04-8171-f95e25090d10   |  TOP200万   |  getblock.io   |
|  wss://bsc.getblock.io/mainnet/?api_key=4a9692ac-6ecd-4ed0-814e-9c541d25df2c   |  TOP200万   |  getblock.io   |
|  wss://bsc.getblock.io/mainnet/?api_key=3502c0be-b554-4c9f-9879-010d5915528c   |  TOP200万   |  getblock.io   |
|  wss://bsc.getblock.io/mainnet/?api_key=8e4937dc-a76e-403e-801d-e5686cdf667d   |  TOP200万   |  getblock.io   |



# 项目结构

## 项目主入口
    OnchainParser > main.go
## 配置文件路径
    OnchainParser > manifest > config
## Web配置路径
    OnchainParser > internal > cmd
## SQL文件存放路径
    OnchainParser > sql
## Job任务路径
    OnchainParser > internal > task
## Controller路径
    OnchainParser > internal > controller
## Web3配置路径
    OnchainParser > internal > web3
## 接口路径
    OnchainParser > internal > service
## 业务逻辑路径
    OnchainParser > internal > logic
## 数据Dao路径
    OnchainParser > internal > dao
## 数据Model路径
    OnchainParser > internal > model
## 公共数据路径
    OnchainParser > internal > consts
