package service

import (
	"fmt"
	"myapp/internal/errorCode"
	"myapp/utility/response"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sMiddleware struct{}

var insMiddleware = sMiddleware{}

// 中间件管理服务
func Middleware() *sMiddleware {
	return &insMiddleware
}

/**
 * @Description I18N中间件，根据Header上的Lang参数或者Query参数来设置当前的I18N.Query参数优先级高于header。
 **/
func (s *sMiddleware) I18NMiddleware(r *ghttp.Request) {
	configLang, _ := g.Cfg().Get(r.Context(), "server.lang", "zh_CN")
	lang := fmt.Sprint(configLang)
	lang1 := r.GetHeader("Lang") // 获取不到返回""
	lang2 := r.GetQuery("Lang")  // 获取不到返回 nil
	// url参数Lang优先级高于header的Lang
	if gconv.Bool(lang1) {
		lang = lang1
	}
	if lang2 != nil {
		lang = fmt.Sprint(lang2.Val())
	}
	g.Log().Infof(r.Context(), "切换当前语言为：%s \n", lang)
	r.SetCtx(gi18n.WithLanguage(r.Context(), lang))
	g.Log().Info(r.Context(), g.I18n().Tf(r.Context(), `{#hello}`, "beep"))

	r.Middleware.Next()
}

// 返回处理中间件
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有返回内容，那么该中间件什么也不做
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err  error
		res  interface{}
		code gcode.Code = gcode.CodeOK
	)
	res, err = r.GetHandlerResponse()
	if err != nil {

		code = gerror.Code(err)
		if code == errorCode.CodeNil { // code是可比较的结构体
			code = errorCode.CodeInternalError
		}
		if detail, ok := code.Detail().(errorCode.MyCodeDetail); ok {
			r.Response.WriteStatus(detail.HttpCode)
			r.Response.ClearBuffer() // gf 会自动往response追加http.StatusText。此处不需要，所以删除掉。
		}
		g.Log().Errorf(r.GetCtx(), "%+v", err)
		response.JsonExit(r, code.Code(), err.Error())
	} else {
		response.JsonExit(r, code.Code(), "", res)
	}
}

//// 自定义上下文对象
//func (s *sMiddleware) Ctx(r *ghttp.Request) {
//	// 初始化，务必最开始执行
//	customCtx := &model.Context{
//		Session: r.Session,
//		Data:    make(g.Map),
//	}
//	Context().Init(r, customCtx)
//	if userEntity := Session().GetUser(r.Context()); userEntity.Id > 0 {
//		adminId := g.Cfg().MustGet(r.Context(), "setting.adminId", consts.DefaultAdminId).Uint()
//		customCtx.User = &model.ContextUser{
//			Id:       userEntity.Id,
//			Passport: userEntity.Passport,
//			Nickname: userEntity.Nickname,
//			Avatar:   userEntity.Avatar,
//			IsAdmin:  userEntity.Id == adminId,
//		}
//	}
//	// 将自定义的上下文对象传递到模板变量中使用
//	r.Assigns(g.Map{
//		"Context": customCtx,
//	})
//	// 执行下一步请求逻辑
//	r.Middleware.Next()
//}
//
//// 前台系统权限控制，用户必须登录才能访问
//func (s *sMiddleware) Auth(r *ghttp.Request) {
//	user := Session().GetUser(r.Context())
//	if user.Id == 0 {
//		_ = Session().SetNotice(r.Context(), &model.SessionNotice{
//			Type:    consts.SessionNoticeTypeWarn,
//			Content: "未登录或会话已过期，请您登录后再继续",
//		})
//		// 只有GET请求才支持保存当前URL，以便后续登录后再跳转回来。
//		if r.Method == "GET" {
//			_ = Session().SetLoginReferer(r.Context(), r.GetUrl())
//		}
//		// 根据当前请求方式执行不同的返回数据结构
//		if r.IsAjaxRequest() {
//			response.JsonRedirectExit(r, 1, "", s.LoginUrl)
//		} else {
//			r.Response.RedirectTo(s.LoginUrl)
//		}
//	}
//	r.Middleware.Next()
//}
