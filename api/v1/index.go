package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type WelcomeReq struct {
	g.Meta `path:"/" tags:"Welcome" method:"get" summary:"Welcome page""`
}
