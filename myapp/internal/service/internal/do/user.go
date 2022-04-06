// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-04-06 16:59:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta      `orm:"table:user, do:true"`
	Uuid        interface{} // uuid
	LoginName   interface{} // 登录名
	DisplayName interface{} // 姓名
	Password    interface{} // 密码
	Email       interface{} // 邮箱
	Phone       interface{} // 电话
	Enabled     interface{} // 用户的启用状态，1表示enabled：启用，0表示disabled：禁用
	Desc        interface{} // 描述信息
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 最后修改时间
}
