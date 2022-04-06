package response

import (
	"reflect"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// JsonRes 数据返回通用JSON数据结构
type JsonRes struct {
	Code    int         `json:"code"`    // 错误码((0:成功, 1:失败, >1:错误码))
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
}

// Json 返回标准JSON数据。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
		if reflect.ValueOf(responseData).IsNil() { // responseData是interface类型，判空需要使用反射；
			responseData = g.Map{} // response 为nil时，返回{}
		}
	} else {
		responseData = g.Map{}
	}
	r.Response.WriteJson(JsonRes{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8") // 重置response head增加charset=utf-8
}

// JsonExit 返回标准JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	Json(r, code, message, data...)
	r.Exit()
}
