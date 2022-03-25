package model

import (
	"myapp/api"

	"github.com/gogf/gf/v2/os/gtime"
)

type UserGetOutput struct {
	Uuid        string      // uuid
	Password    string      // 密码
	LoginName   string      // 登录名
	DisplayName string      // 姓名
	Email       string      // 邮箱
	Phone       string      // 电话
	Enabled     string      // 用户的启用状态，enabled表示启用，disabled表示禁用
	Desc        string      // 描述信息
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 最后修改时间
}

type UserCreateInput struct {
	api.UserCreateReq
}

type UserListInput struct {
	SearchStr string `json:"searchStr" in:"query" dc:"模糊查询，匹配登录名/姓名/电话/描述"`
	Page      int    `json:"page" description:"分页码"`
	Size      int    `json:"size" description:"分页数量"`
	SortKey   string `json:"sortKey"   in:"query" dc:"排序字段"`
	SortValue string `json:"sortValue"   in:"query" dc:"排序顺序，ASC升序，DESC降序"`
}

type UserItem struct {
	User *UserGetOutput `json:"user" description:"用户详情"`
}

type UserListOutput struct {
	List  []UserItem `json:"users" description:"用户列表"`
	Page  int        `json:"page" description:"分页码"`
	Size  int        `json:"size" description:"分页数量"`
	Total int        `json:"total" description:"数据总数"`
}
