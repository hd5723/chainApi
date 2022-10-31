package web3Client

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"math/rand"
)

func GetClientByLinks(ctx context.Context, links g.ArrayStr, chainName string) (cli *ethclient.Client, link string, err error) {

	linksLen := rand.Intn(len(links))
	if linksLen == 0 {
		linksLen = 1
	}
	// 从redis 顺序取之：1、2、3、4、5、6、7、8...∞
	key := chainName + "_num"
	res, err := g.Redis().Do(ctx, "incr", key)
	if err != nil {
		g.Log().Error(ctx, "redis do fail")
		link = links[0]
	} else {
		link = links[res.Int()%linksLen]
	}
	//通过rpc长度取模，0、1、2、3...rpc长度-1；0、1、2、3...rpc长度-1
	client, err := ethclient.Dial(link)
	if err != nil {
		g.Log().Error(ctx, "web3 rpc.Dial err")
		return nil, link, err
	}
	return client, link, nil
}

func GetLastBlockNumber(ctx context.Context, client *ethclient.Client) (result int64, err error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	result = header.Number.Int64()
	return result, nil
}
