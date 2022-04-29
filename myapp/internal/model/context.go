package model

import (
	"myapp/utility/token"

	"github.com/gogf/gf/v2/frame/g"
)

// Context 请求上下文结构
type Context struct {
	Token *token.MyCacheToken // token信息，包含上下文用户信息
	Data  g.Map               // 自定KV变量，业务模块根据需要设置，不固定
}
