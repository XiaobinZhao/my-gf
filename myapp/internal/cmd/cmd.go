package cmd

import (
	"context"
	"myapp/internal/consts"
	"myapp/internal/controller"
	"myapp/internal/service"

	"github.com/gogf/gf/v2/os/glog"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/protocol/goai"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start MyGoFrame server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			SetLoggerDefaultHandler() // 替代默认的log

			s := g.Server()

			err = g.I18n().SetPath("resource/i18n") // i18n目录默认是gres资源目录或者根目录；在研发阶段需要重设一下i18n目录
			if err != nil {
				panic(err)
			}
			lang := g.Cfg().MustGet(context.Background(), "server.lang", "zh_CN").String()
			g.I18n().SetLanguage(lang) // 设置全局i18N语言
			g.Log().Infof(context.Background(), "全局设置当前语言为：%s", lang)
			// 设置行号，日期，时间：日期+时间+毫秒，如：2009-01-23 01:23:23.675
			//g.Log().SetFlags(g.Log().GetFlags() | glog.F_FILE_SHORT | glog.F_TIME_DATE | glog.F_TIME_MILLI) // 此处注释掉，通过配置文件实现

			//全局中间件注册
			s.Use(
				//service.Middleware().MiddlewareCORS,
				service.Middleware().Ctx,
				service.Middleware().I18NMiddleware,
				service.Middleware().ResponseHandler,
			)

			// 不需要token校验的路由
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.POST("/login", controller.Login.Login)
			})

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().TokenAuth,
				)
				//group.Bind(
				//	controller.Token,    // 用户
				//	controller.Desktop, // 桌面
				//)
				// 官方文档建议使用对象注册（如上）的规范路由：path和method写到API的struct中，但是并没有写到一个地方感觉便于管理
				group.DELETE("/logout/{uuid}", controller.Login.Logout)
				group.Group("/users", func(group *ghttp.RouterGroup) {
					group.GET("/", controller.User.List)
					group.POST("/", controller.User.Create)
					group.GET("/{uuid}", controller.User.Get)
					group.PATCH("/{uuid}", controller.User.Update)
					group.DELETE("/{uuid}", controller.User.Delete)
				})
				group.Group("/desktops", func(group *ghttp.RouterGroup) {
					group.GET("/", controller.Desktop.List)
					group.POST("/", controller.Desktop.Create)
					group.GET("/{uuid}", controller.Desktop.Get)
					group.PATCH("/{uuid}", controller.Desktop.Update)
					group.DELETE("/{uuid}", controller.Desktop.Delete)
				})
			})

			// 自定义丰富文档
			enhanceOpenAPIDoc(s)
			// 启动Http Server
			s.Run()
			return
		},
	}
)

func enhanceOpenAPIDoc(s *ghttp.Server) {
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`

	// API description.
	openapi.Info.Title = `MyApp`
	openapi.Info.Version = "1.0.0"
	openapi.Info.Description = `基于'GoFrame'定制化API Server框架`

	// Sort the tags in custom sequence.
	openapi.Tags = &goai.Tags{
		{Name: consts.OpenAPITagNameUser, Description: "用户管理"},
		{Name: consts.OpenAPITagNameDesktop, Description: "云桌面管理"},
		{Name: consts.OpenAPITagNameAuthorization, Description: "认证管理"},
	}

	openapi.Components = goai.Components{
		SecuritySchemes: goai.SecuritySchemes{
			"BearerAuth": goai.SecuritySchemeRef{
				Ref: "", // 暂时还不知道该值是干什么用的
				Value: &goai.SecurityScheme{
					Type:   "http",
					In:     "header",
					Name:   "Authorization",
					Scheme: "bearer",
				},
			},
		},
	}
}

// 替代默认的日志handler，禁止控制台输出，全部输出到文件
func SetLoggerDefaultHandler() {
	glog.SetDefaultHandler(func(ctx context.Context, in *glog.HandlerInput) {
		m := map[string]interface{}{
			"stdout": g.Config().MustGet(ctx, "logger.globalStdout", true).Bool(), // 使用自定义的全局字段
			"path":   g.Config().MustGet(ctx, "logger.path", "log/").String(),     // 此处必须重新设置，才可以实现db的log写入文件
		}
		in.Logger.SetConfigWithMap(m)
		in.Next(ctx)
	})
}
