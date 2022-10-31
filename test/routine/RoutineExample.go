package routineTest

import (
	"OnchainParser/internal/chain/ethEvm/consts"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
	"time"
)

var (
	taskList chan TranstTask
	wg       sync.WaitGroup
)

type TranstTask struct {
	Beigin_point int // 开始节点
	End_point    int // 结束节点
	Retry_number int // 重试次数
}

type TaskInfo struct {
	Root_point  int //根结点
	Child_point int //子节点
}

func doSubscribe(ctx context.Context, info TaskInfo, workList chan TranstTask) {

	var task TranstTask
	task.Beigin_point = info.Child_point
	task.End_point = 100 + info.Child_point
	task.Retry_number = 2
	if info.Child_point == 7 {
		task.End_point = 8888 + info.Child_point
		workList <- task
		return
	}

	g.Log().Info(ctx, "ping: Root_point ", info.Root_point, " , Child_point:", info.Child_point)
	time.Sleep(time.Second * 1)
	workList <- task
	return

	g.Log().Info(ctx, "xxxxxxxx")

}

func DoRuntine(ctx context.Context) {

	//pool := grpool.New(100)
	//for i := 0; i < 1000; i++ {
	//	pool.Add(ctx, job)
	//}
	//
	//g.Log().Info(ctx, "do something!!")

	start := gtime.TimestampMilli()

	n := 10
	workList := make(chan TranstTask)
	pool := grpool.New(5)
	for i := 0; i < n; i++ {
		var blockInfo consts.BlockInfo
		blockInfo.BlockNumber = int64(i)
		pool.Add(ctx, func(ctx context.Context) {
			doTask(ctx, blockInfo, workList)
		})
	}
	str := make([]TranstTask, n)
	for i := 0; i < n; i++ {
		str[i] = <-workList
	}
	close(workList)
	g.Log().Info(ctx, "time spent:", gtime.TimestampMilli()-start)

}

func doTask(ctx context.Context, blockInfo consts.BlockInfo, rootWorkList chan TranstTask) {
	rootIndex := int(blockInfo.BlockNumber)
	g.Log().Info(ctx, rootIndex)
	n := 100
	workList := make(chan TranstTask)
	pool := grpool.New(30)
	for i := 0; i < n; i++ {
		var v TaskInfo
		v.Child_point = i
		v.Root_point = rootIndex
		pool.Add(ctx, func(ctx context.Context) {
			doSubscribe(ctx, v, workList)
		})
		//go DoSubscribe(ctx, v, workList)
	}

	str := make([]TranstTask, n)
	for i := 0; i < n; i++ {
		str[i] = <-workList
	}
	close(workList)

	var task TranstTask
	task.Beigin_point = rootIndex
	task.End_point = 100 + rootIndex
	task.Retry_number = 2
	rootWorkList <- task
	//g.Log().Info(ctx, str)
}
