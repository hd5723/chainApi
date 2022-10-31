package cmd

import (
	"OnchainParser/internal/utils"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"OnchainParser/internal/controller"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Index,
				)
			})
			s.Group("/protocol", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.Protocol,
				)
			})
			s.Group("/contract", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.Contract,
				)
			})
			s.Group("/functionInfo", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.FunctionInfo,
				)
			})
			s.Group("/abi", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.Abi,
				)
			})
			s.Group("/eventInfo", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.EventInfo,
				)
			})
			s.Group("/transaction", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.Transaction,
				)
			})
			s.Group("/chain", func(group *ghttp.RouterGroup) {
				group.Middleware(utils.MiddlewareAuth)
				group.Bind(
					controller.Chain,
				)
			})
			s.Run()
			return nil
		},
	}
)
