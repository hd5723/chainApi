package utils

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// DefaultHandlerResponse is the default implementation of HandlerResponse.
type DefaultHandlerResponse struct {
	Code    int         `json:"code"    dc:"Error code"`
	Message string      `json:"message" dc:"Error message"`
	Data    interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

func MiddlewareAuth(r *ghttp.Request) {
	auth, er := g.Cfg().Get(context.Background(), "chain.auth")
	if er != nil {
		r.Response.WriteStatus(http.StatusForbidden)
	}

	apiKey := r.GetHeader("apiKey")
	if g.IsEmpty(apiKey) {
		r.Response.WriteStatus(http.StatusUnauthorized)
	}

	if auth.String() == apiKey {
		r.Middleware.Next()
	} else {
		r.Response.WriteJson(DefaultHandlerResponse{
			Code:    http.StatusUnauthorized,
			Message: "the apiKey is wrong",
			Data:    nil,
		})
	}

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		switch r.Response.Status {
		case http.StatusNotFound:
			code = gcode.CodeNotFound
		case http.StatusForbidden:
			code = gcode.CodeNotAuthorized
		default:
			code = gcode.CodeUnknown
		}
	} else {
		code = gcode.CodeOK
	}
	r.Response.WriteJson(DefaultHandlerResponse{
		Code:    code.Code(),
		Message: msg,
		Data:    res,
	})
}
