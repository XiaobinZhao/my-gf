package service

import (
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/utility/response"
	"myapp/utility/token"

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
	configLang := g.Cfg().MustGet(r.Context(), "server.lang", "zh-CN").String()

	langInHeader := r.GetHeader("Lang")        // 获取不到返回""
	langInQuery := r.GetQuery("Lang").String() // 获取不到返回 nil
	// url参数Lang优先级高于header的Lang
	requestLang := ""
	if gconv.Bool(langInHeader) {
		requestLang = langInHeader
	}
	if gconv.Bool(langInQuery) {
		requestLang = langInQuery
	}
	if requestLang != "" && requestLang != configLang {
		g.Log().Debugf(r.Context(), "切换当前语言为：%s", requestLang)
		r.SetCtx(gi18n.WithLanguage(r.Context(), requestLang))
	}
	r.Middleware.Next()
}

// 返回处理中间件
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有返回内容，那么该中间件什么也不做
	if r.Response.BufferLength() > 0 {
		g.Log().Warningf(r.GetCtx(), "response exists something, skip ResponseHandler middleware")
		return
	}

	res, err := r.GetHandlerResponse()

	formatResponse(r, res, err)
}

func formatResponse(r *ghttp.Request, res interface{}, err error) {

	var (
		code gcode.Code = gcode.CodeOK
	)
	if err != nil {
		code = gerror.Code(err)
		if code == errorCode.CodeNil { // code是可比较的结构体
			code = errorCode.CodeInternalError
		}
		if detail, ok := code.Detail().(errorCode.MyCodeDetail); ok {
			r.Response.WriteStatus(detail.HttpCode) // 修改默认的状态码，并清除已经写入的response内容
			r.Response.ClearBuffer()                // gf 会自动往response追加http.StatusText。此处不需要，所以删除掉。
		}
		g.Log().Errorf(r.GetCtx(), "%+v", err)
		response.JsonExit(r, code.Code(), gerror.Current(err).Error()) // 只暴露当前error给调用者
	} else {
		response.JsonExit(r, code.Code(), "", res)
	}
}

func (s *sMiddleware) MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *sMiddleware) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Data: make(g.Map),
	}
	Context().Init(r, customCtx)

	r.Middleware.Next()
}

/**
 * @description 认证token;1. 从request获取token 2. 从cache获取对应token，进行校验  3. 校验成功，把cache缓存的数据放入context
 * @param r *ghttp.Request
 **/
func (s *sMiddleware) TokenAuth(r *ghttp.Request) {
	// 1. 从request获取token
	tokenStr, err := token.GetRequestToken(r)
	if err != nil {
		formatResponse(r, nil, err)
	}
	// 2. 从cache获取对应token，进行校验
	myCacheToken, err := token.Instance().ValidToken(r.Context(), tokenStr)
	if err != nil {
		formatResponse(r, nil, err)
	}
	// 3. 校验成功，把cache缓存的数据放入context
	Context().SetToken(r.Context(), myCacheToken)
	r.Middleware.Next()
}
