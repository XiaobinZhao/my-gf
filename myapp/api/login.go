package api

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `summary:"用户登录" tags:"认证"`
	UserName string `json:"userName" p:"userName" v:"required"  dc:"登录名"`
	Password string `json:"password" p:"password" v:"required"  dc:"密码"`
}

type LogoutReq struct {
	g.Meta `summary:"用户退出" tags:"认证"`
	Uuid   string `json:"uuid" p:"uuid"  v:"required" in:"path" dc:"用户UUID"`
}

type LoginRes struct {
	User  *UserGetRes `json:"user" description:"用户信息"`
	Token string      `json:"token" description:"token信息"`
}

type LogoutRes struct {
}
