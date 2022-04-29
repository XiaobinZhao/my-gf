package errorCode

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type MyCode struct {
	code    int
	message string // message 设计为i18n的key
	detail  MyCodeDetail
}
type MyCodeDetail struct {
	HttpCode int
}

func (c MyCode) MyDetail() MyCodeDetail {
	return c.detail
}

func (c MyCode) Code() int {
	return c.code
}

func (c MyCode) Message() string {
	return c.message
}

func (c MyCode) Detail() interface{} {
	return c.detail
}

func New(httpCode int, code int, message string) gcode.Code {
	return MyCode{
		code:    code,
		message: message,
		detail: MyCodeDetail{
			HttpCode: httpCode,
		},
	}
}

func NewMyErr(ctx context.Context, code gcode.Code, params ...interface{}) error {
	tfStr := g.I18n().Tf(ctx, code.Message(), params...)
	return gerror.NewCode(code, tfStr)
}

func MyWrapCode(ctx context.Context, code gcode.Code, err error, params ...interface{}) error {
	tfStr := g.I18n().Tf(ctx, code.Message(), params...)
	return gerror.WrapCode(code, err, tfStr)
}

// code 码要大于1000,1000以内gf框架内使用
// code的设计关系到一个问题的争议：异常处理的HTTP响应状态码是否依然返回200？
// https://stackoverflow.com/questions/27921537/returning-http-200-ok-with-error-within-response-body
// 1. 如果我们明确API是REST的，而且API对外使用，应当使用合适的状态码来反映错误（建议控制在20个以内常用的），并且在文档中进行说明，
//    而且出错后需要在响应体补充细化的error信息（包含code和message）
// 2. 如果REST API对内使用，那么在客户端和服务端商量好统一标准的情况下可以对响应码类型进行收敛到几个，实现起来也方便
// 3. 如果API是内部使用的RPC over HTTP形式，甚至可以退化到业务异常也使用200响应返回
// 本项目希望尽可能的遵守RESTful规范，使用合适的状态码来反映错误，并且返回统一的response来进行错误说明。
var (
	// gf框架内置的，参见：github.com\gogf\gf\v2@v2.0.0-rc2\errors\gcode\gcode.go
	CodeNil           = New(200, -1, "")
	CodeNotFound      = New(404, 65, "Not Found")
	CodeInternalError = New(500, 50, "An error occurred internally")

	// 系统起始 10000
	MyInternalError = New(500, 10001, "{#myInternalError}")

	// token 20000起始
	AuthHeaderInvalidError     = New(401, 20001, `{#authHeaderInvalidError}`)
	NotSupportedCacheModeError = New(401, 20002, `{#notSupportedCacheModeError}`)
	TokenEmpty                 = New(401, 20003, `{#tokenEmpty}`)
	TokenKeyEmpty              = New(401, 20004, `{#tokenKeyEmpty}`)
	TokenInvalidError          = New(401, 20005, `{#tokenInvalidError}`)
	Unauthorized               = New(401, 20006, `{#unauthorized}`)
	AuthorizedFailed           = New(401, 20007, `{#authorizedFailed}`)

	//用户30000起始
	UserNotFound        = New(404, 30001, `{#userNotExists}`)
	LoginNameConflicted = New(403, 30002, `{#loginNameConflicted}`)
	PasswordError       = New(401, 30003, `{#passwordError}`)
	LoginFailed         = New(401, 30004, `{#loginFailed}`)

	// 桌面40000起始
	DesktopNotFound = New(404, 40001, `{#desktopNotExists}`)
)
