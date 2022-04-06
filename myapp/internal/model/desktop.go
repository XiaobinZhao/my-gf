package model

import (
	"myapp/api"

	"github.com/gogf/gf/v2/os/gtime"
)

type DesktopGetOutput struct {
	Uuid      string      // uuid
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 最后修改时间
	api.DesktopBase
}

type DesktopCreateInput struct {
	api.DesktopBase
}

type DesktopUpdateInput struct {
	api.DesktopBase
	Uuid string
}

type DesktopListInput struct {
	SearchStr string `json:"searchStr" in:"query" dc:"模糊查询，匹配名称/描述"`
	Page      int    `json:"page" description:"分页码"`
	Size      int    `json:"size" description:"分页数量"`
	SortKey   string `json:"sortKey"   in:"query" dc:"排序字段"`
	SortValue string `json:"sortValue"   in:"query" dc:"排序顺序，ASC升序，DESC降序"`
}

type DesktopListOutput struct {
	api.DesktopListRes
}
