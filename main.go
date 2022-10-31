package main

import (
	"OnchainParser/internal/chain/ethEvm/chain/baseConfig"
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/cmd"
	_ "OnchainParser/internal/logic"
	_ "OnchainParser/internal/packed"
	"OnchainParser/internal/task"
	"flag"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/housemecn/snowflake"
)

var (
	ctx        = gctx.New()
	runTypeVal = flag.String("runType", "", "api task")
)

func main() {

	//设置每次处理区块数量
	chainInfoName, err := g.Cfg().Get(ctx, "chain.run.name")
	if err != nil {
		return
	}

	//设置每次处理区块数量
	chainInfoId, err := g.Cfg().Get(ctx, "chain.run.id")
	if err != nil {
		return
	}

	baseConfig.CHAIN_ID = chainInfoId.Int()
	baseConfig.CHAIN_NAME = chainInfoName.String()
	// 设置 baseConfig.CHAIN_LAST_HEIGHT
	task.SetChainLastHeight(ctx)

	//重启时，释放 最新区块爬取任务 锁
	lockKey := chainInfoName.String() + "_" + consts.LOCK_TRANSACTION_KEY + "_*"
	v, err := g.Redis().Do(ctx, "keys", lockKey)
	if err == nil {
		for i := 0; i < len(v.Array()); i++ {
			key := v.Array()[i]
			g.Redis().Do(ctx, "DEL", key)
		}
	}

	flag.Parse()
	if g.IsEmpty(*runTypeVal) {
		//启动Task
		task.TaskManage(gctx.New())

		//启动Web
		cmd.Main.Run(gctx.New())
	}

	if *runTypeVal == "task" {
		//启动Task
		task.TaskManage(gctx.New())
		select {}
	}
	if *runTypeVal == "api" {
		//启动Web
		cmd.Main.Run(gctx.New())
	}

	g.Log().Info(ctx, "the application run success!")
}
