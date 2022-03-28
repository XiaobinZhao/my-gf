package main

import (
	"myapp/internal/consts"
	"myapp/internal/controller"
	"myapp/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/protocol/goai"
)

// @title       基于‘GoFrame’定制化API Server框架
// @version     1.0
// @description 基于‘GoFrame’定制化API Server框架
// @schemes     http
func main() {
	s := g.Server()
	// 设置行号，日期，时间：日期+时间+毫秒，如：2009-01-23 01:23:23.675
	//g.Log().SetFlags(g.Log().GetFlags() | glog.F_FILE_SHORT | glog.F_TIME_DATE | glog.F_TIME_MILLI)  // 通过配置文件实现
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			service.Middleware().I18NMiddleware,
			//service.Middleware().Ctx,
			service.Middleware().ResponseHandler,
		)
		group.Bind(
			controller.User, // 用户
		)
	})
	// 自定义丰富文档
	enhanceOpenAPIDoc(s)
	// 启动Http Server
	s.Run()
}

func enhanceOpenAPIDoc(s *ghttp.Server) {
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`

	// API description.
	openapi.Info.Title = `MyApp`
	openapi.Info.Description = `基于'GoFrame'定制化API Server框架`

	// Sort the tags in custom sequence.
	openapi.Tags = &goai.Tags{
		{Name: consts.OpenAPITagNameUser},
		{Name: consts.OpenAPITagNameDesktop},
	}
}
