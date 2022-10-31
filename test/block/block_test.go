package main

import (
	"flag"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"math/big"
	"testing"
	"time"
)

var (
	ctx = gctx.New()
)
var cAddress = flag.String("address", "", "contract address")
var fromBlock = flag.Int("fromBlock", 0, "fromBlock")
var link = flag.String("link", "wss://bsc-mainnet.nodereal.io/ws/v1/271ce51c526244aa839f81e2ff2e92e3", "link")

func TestFilterLog(t *testing.T) {
	g.Log().Info(ctx, "TestFilterLog start")

	client, err := ethclient.Dial(*link)
	if err != nil {
		g.Log().Error(ctx, "web3 rpc.Dial err")
		return
	}

	toBlock := *fromBlock + 100
	contractAddress := common.HexToAddress(*cAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(*fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	curTime := time.Now()
	g.Log().Info(ctx, "RPC client.FilterLogs start fromBlock:", fromBlock, " ,toBlock:", toBlock, " , contract_code:", cAddress)
	logs, err := client.FilterLogs(ctx, query)
	g.Log().Info(ctx, "RPC block.subscribeContract.client.FilterLogs doing time:", time.Now().Sub(curTime), " ,link:", link, " , contract_code:", cAddress)
	if err != nil {
		g.Log().Error(ctx, "RPC block.subscribeContract.client.FilterLogs error:", err)
		return
	}
	g.Log().Info(ctx, "RPC block.subscribeContract.client.FilterLogs result:", logs)
	g.Log().Info(ctx, "TestFilterLog end")
}
