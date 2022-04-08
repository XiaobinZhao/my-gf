package main

import (
	"myapp/internal/consts"
	"myapp/internal/controller"
	"myapp/internal/service"

	_ "myapp/packed"

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
	g.I18n().SetPath("resource/i18n") // i18n目录默认是根目录或者gres资源目录；在研发阶段需要重设一下i18n目录
	// 设置行号，日期，时间：日期+时间+毫秒，如：2009-01-23 01:23:23.675
	//g.Log().SetFlags(g.Log().GetFlags() | glog.F_FILE_SHORT | glog.F_TIME_DATE | glog.F_TIME_MILLI) // 此处注释掉，通过配置文件实现
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			//service.Middleware().Ctx,
			service.Middleware().I18NMiddleware,
			service.Middleware().ResponseHandler,
		)
		group.Bind(
			controller.User,    // 用户
			controller.Desktop, // 桌面
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
