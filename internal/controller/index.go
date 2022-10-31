package controller

import (
	"OnchainParser/api/v1"
	"OnchainParser/internal/utils"
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	Index = cIndex{}
)

type cIndex struct{}

func (c *cIndex) Welcome(ctx context.Context, req *v1.WelcomeReq) (res *utils.ResponseRes, err error) {
	re := "Welcome!"
	g.RequestFromCtx(ctx).Response.WriteJson(re)
	return
}
