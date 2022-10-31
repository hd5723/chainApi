package chainLoader

import (
	"OnchainParser/internal/chain/ethEvm/utils/loader/blockLoader"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"OnchainParser/internal/web3/web3Client"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
)

// 爬取历史区块
func DoScanHistoryBlock(ent entity.ContractEntity, chainName string, chainId int) {
	//开启链路追踪
	ctx := gctx.New()
	rcd := service.Chain().QueryOneByChainId(ctx, int32(chainId), chainName)
	links := strings.Split(rcd.RpcUrls, ",")
	ethClient, link, err := web3Client.GetClientByLinks(ctx, links, chainName)
	if err != nil {
		return
	}
	blockLoader.DoScanHistoryBlock(ctx, ethClient, link, ent, chainName)
	ethClient.Close()
}

// 爬取最新区块
func DoScanLastBlock(chainName string, chainId int) {
	blockLoader.DoScanLastBlock(chainId, chainName)
}

// 订阅历史区块
func DoSubscribeHistoryContract(ent entity.ContractEntity, chainName string, chainId int) {
	blockLoader.DoSubscribeHistoryContract(ent, chainName)
}

func Undo() {

}
