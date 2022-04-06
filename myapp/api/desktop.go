package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type DesktopGetReq struct {
	g.Meta `path:"/desktop/{uuid}" method:"get" summary:"获取单个桌面" tags:"云桌面"`
	Uuid   string `json:"uuid" p:"uuid"  v:"required" in:"path" dc:"桌面UUID"`
}

type DesktopListReq struct {
	g.Meta `path:"/desktop" method:"get" summary:"获取桌面列表" tags:"云桌面"`
	CommonPaginationReq
	SearchStr string `json:"searchStr" p:"searchStr" in:"query" dc:"模糊查询，匹配名称/描述"`
	Sort      string `json:"sort" p:"sort" v:"regex:^[+-].+" in:"query" dc:"排序字段，+表示升序ASC，-表示降序DESC,比如+name：表示按照name升序排列"`
}

type DesktopCreateReq struct {
	g.Meta `path:"/desktop" method:"post" summary:"创建桌面" tags:"云桌面"`
	DesktopBase
}

type DesktopDeleteReq struct {
	g.Meta `path:"/desktop/{uuid}" method:"delete" summary:"删除单个桌面" tags:"云桌面"`
	Uuid   string `json:"uuid" p:"uuid"  v:"required" in:"path" dc:"桌面UUID"`
}

type DesktopUpdateReq struct {
	g.Meta          `path:"/desktop/{uuid}" method:"patch" summary:"更新单个桌面" tags:"云桌面"`
	Uuid            string `json:"uuid" p:"uuid"  v:"required" in:"path" dc:"桌面UUID"`
	VmUuid          string `json:"vmUuid" p:"vmUuid" v:"max-length:32" dc:"虚拟化平台上虚机的uuid"`
	DisplayName     string `json:"displayName" p:"displayName" dc:"桌面的显示名称"`
	Enabled         string `json:"enabled" p:"enabled" v:"in:enabled,disabled" dc:"桌面的启用状态，enabled表示启用，disabled表示禁用"`
	GpuAttachStatus string `json:"gpuAttachStatus" p:"gpuAttachStatus" v:"in:preAttached,attached,unattached" dc:"gpu挂载状态"`
	IsDefault       int    `json:"isDefault" p:"isDefault" v:"in:0,1" dc:"是否是默认桌面，0表示否，1表示是"`
	NodeUuid        string `json:"nodeUuid" p:"nodeUuid" v:"max-length:32" dc:"云桌面所在的物理机节点uuid"`
	NodeName        string `json:"nodeName" p:"nodeName" dc:"云桌面所在的物理机节点名称"`
	Desc            string `json:"desc" p:"desc" v:"max-length:255"  dc:"描述信息"`
}

type DesktopBase struct {
	VmUuid          string `json:"vmUuid" p:"vmUuid" v:"required|max-length:32" dc:"虚拟化平台上虚机的uuid"`
	DisplayName     string `json:"displayName" p:"displayName" v:"required" dc:"桌面的显示名称"`
	Enabled         string `json:"enabled" p:"enabled" d:"enabled" v:"required|in:enabled,disabled" dc:"桌面的启用状态，enabled表示启用，disabled表示禁用"`
	GpuAttachStatus string `json:"gpuAttachStatus" p:"gpuAttachStatus" d:"unattached" v:"required|in:preAttached,attached,unattached" dc:"gpu挂载状态"`
	IsDefault       int    `json:"isDefault" p:"isDefault" v:"required|in:0,1" dc:"是否是默认桌面，0表示否，1表示是"`
	NodeUuid        string `json:"nodeUuid" p:"nodeUuid" v:"max-length:32" dc:"云桌面所在的物理机节点uuid"`
	NodeName        string `json:"nodeName" p:"nodeName" dc:"云桌面所在的物理机节点名称"`
	Desc            string `json:"desc" p:"desc" v:"max-length:255"  dc:"描述信息"`
}

type DesktopGetRes struct {
	Uuid            string      `json:"uuid"            description:"桌面uuid"`
	VmUuid          string      `json:"vmUuid"          description:"虚拟化平台上虚机的uuid"`
	DisplayName     string      `json:"displayName"     description:"桌面的显示名称"`
	GpuAttachStatus string      `json:"gpuAttachStatus" description:"桌面挂载gpu状态，'preAttached'表示预挂载，即关联了GPU规格的关机态虚机,'attached'表示已经挂载,'unattached'表示未挂载"`
	Enabled         string      `json:"enabled"         description:"桌面的启用状态，enabled表示启用，disabled表示桌面已禁用"`
	NodeUuid        string      `json:"nodeUuid"        description:"物理机uuid"`
	NodeName        string      `json:"nodeName"        description:"物理机名称"`
	IsDefault       int         `json:"isDefault"       description:"是否是默认桌面，False表示不是，True表示是。默认False。"`
	Desc            string      `json:"desc"            description:"桌面的描述信息"`
	CreatedAt       *gtime.Time `json:"createdAt"       description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       description:"最后修改时间"`
}

type DesktopListRes struct {
	List  []DesktopGetRes `json:"desktops" description:"桌面列表"`
	Page  int             `json:"page" description:"分页码"`
	Size  int             `json:"size" description:"分页数量"`
	Total int             `json:"total" description:"数据总数"`
}
