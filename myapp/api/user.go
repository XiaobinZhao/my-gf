package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type UserGetReq struct {
	g.Meta   `path:"/user/{userUuid}" method:"get" summary:"获取单个用户" tags:"用户"`
	UserUuid string `json:"userUuid" p:"userUuid"  v:"required" in:"path" dc:"用户UUID"`
}

type UserListReq struct {
	g.Meta `path:"/user" method:"get" summary:"获取用户列表" tags:"用户"`
	CommonPaginationReq
	SearchStr string `json:"searchStr" p:"searchStr" in:"query" dc:"模糊查询，匹配登录名/姓名/电话/描述"`
	Sort      string `json:"sort" p:"sort" v:"regex:^[+-].+" in:"query" dc:"排序字段，+表示升序ASC，-表示降序DESC,比如+name：表示按照name升序排列"`
}

type UserBase struct {
	LoginName   string `json:"loginName" p:"loginName" v:"required|passport"  dc:"登录名,字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
	DisplayName string `json:"displayName" p:"displayName" v:"required" dc:"姓名"`
	Enabled     string `json:"enabled" p:"enabled" v:"required|in:enabled,disabled"  d:"enabled" dc:"用户的启用状态，enabled表示启用，disabled表示禁用"`
	Email       string `json:"email" p:"email" d:"" v:"email"  dc:"邮箱"`
	Phone       string `json:"phone" p:"phone" d:"" v:"phone" dc:"电话"`
	Desc        string `json:"desc" p:"desc" d:"" v:"max-length:255"  dc:"描述信息"`
}

type UserCreateReq struct {
	g.Meta   `path:"/user" method:"post" summary:"创建用户" tags:"用户"`
	Password string `json:"password" p:"password"  v:"required"  dc:"密码"`
	UserBase
}

type UserUpdateReq struct {
	g.Meta      `path:"/user/{userUuid}" method:"patch" summary:"更新单个用户" tags:"用户"`
	UserUuid    string `json:"userUuid" p:"userUuid"  v:"required" in:"path" dc:"用户UUID"`
	LoginName   string `json:"loginName" p:"loginName" v:"passport"  dc:"登录名,字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
	DisplayName string `json:"displayName" p:"displayName" dc:"姓名"`
	Enabled     string `json:"enabled" p:"enabled" v:"in:enabled,disabled" dc:"用户的启用状态，enabled表示启用，disabled表示禁用"`
	Email       string `json:"email" p:"email" v:"email"  dc:"邮箱"`
	Phone       string `json:"phone" p:"phone" v:"phone" dc:"电话"`
	Desc        string `json:"desc" p:"desc" v:"max-length:255"  dc:"描述信息"`
}

type UserGetRes struct {
	Uuid        string      `json:"uuid"        dc:"uuid"`
	LoginName   string      `json:"loginName"   dc:"登录名"`
	DisplayName string      `json:"displayName" dc:"姓名"`
	Email       string      `json:"email"       dc:"邮箱"`
	Phone       string      `json:"phone"       dc:"电话"`
	Enabled     string      `json:"enabled"     dc:"用户的启用状态，enabled表示启用，disabled表示禁用"`
	Desc        string      `json:"desc"        dc:"描述信息"`
	CreatedAt   *gtime.Time `json:"createdAt"   dc:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   dc:"最后修改时间"`
}

type UserListRes struct {
	List  []UserGetRes `json:"users" description:"用户列表"`
	Page  int          `json:"page" description:"分页码"`
	Size  int          `json:"size" description:"分页数量"`
	Total int          `json:"total" description:"数据总数"`
}
