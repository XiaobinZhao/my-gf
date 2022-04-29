package service

import (
	"context"
	"myapp/internal/consts"
	"myapp/internal/model"
	"myapp/utility/token"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sContext struct{}

// Context 上下文管理服务
func Context() *sContext {
	return &sContext{}
}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *sContext) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.ContextKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func (s *sContext) Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// SetToken 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetToken(ctx context.Context, ctxToken *token.MyCacheToken) {
	s.Get(ctx).Token = ctxToken
}

// SetData 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}
