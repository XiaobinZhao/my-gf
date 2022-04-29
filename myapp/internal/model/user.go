package model

import (
	"myapp/api"

	"github.com/gogf/gf/v2/os/gtime"
)

type UserGetOutput struct {
	Uuid      string      // uuid
	Password  string      // 密码
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 最后修改时间
	api.UserBase
}

type UserCreateInput struct {
	Password string // 密码
	api.UserBase
}

type UserLoginInput struct {
	Password string // 密码
	UserName string
}

type UserUpdateInput struct {
	api.UserBase
	Uuid string // userUuid
}

type UserListInput struct {
	SearchStr string `json:"searchStr" in:"query" dc:"模糊查询，匹配登录名/姓名/电话/描述"`
	Page      int    `json:"page" description:"分页码"`
	Size      int    `json:"size" description:"分页数量"`
	SortKey   string `json:"sortKey"   in:"query" dc:"排序字段"`
	SortValue string `json:"sortValue"   in:"query" dc:"排序顺序，ASC升序，DESC降序"`
}

type UserListOutput struct {
	api.UserListRes
}
