// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-03-18 18:02:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Desktop is the golang structure of table desktop for DAO operations like Where/Data.
type Desktop struct {
	g.Meta          `orm:"table:desktop, do:true"`
	Uuid            interface{} // 桌面uuid
	VmUuid          interface{} // 虚拟化平台上虚机的uuid
	DisplayName     interface{} // 桌面的显示名称
	GpuAttachStatus interface{} // 桌面挂载gpu状态，'pre_attached'表示预挂载，即关联了GPU规格的关机态虚机,'attached'表示已经挂载,'unattached'表示未挂载
	Enabled         interface{} // 桌面的启用状态，enabled表示启用，disabled表示桌面已禁用
	NodeUuid        interface{} // 物理机uuid
	NodeName        interface{} // 物理机名称
	IsDefault       interface{} // 是否是默认桌面，False表示不是，True表示是。默认False。
	Desc            interface{} // 桌面的描述信息
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 最后修改时间
}
