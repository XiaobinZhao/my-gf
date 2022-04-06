package service

import (
	"context"
	"strings"
)

// 根据请求的排序字符串获取数据库对应的字段
// 为什么要做字段对应转换？因为前端传参是根据response字段传的，response字段一般使用request model的json格式。
// 而json格式字段和数据库字段以及entity的字段都不一样，比如：json(loginName)和数据库字段(login_name)以及entity的字段(LoginName)
// 字段格式不同，导致字符串比较不一致。所以需要做一下忽略格式的对比，返回对应的数据库字段
func GetSortField(ctx context.Context, inputColumn string, tableColumnsStr string) string {
	for _, column := range strings.Split(tableColumnsStr, ",") {
		// 比较方法：去除反引号`和_;忽略大小写
		tableColumn := column[1 : len(column)-1]                           // 去除字符串的反引号`
		column = strings.ReplaceAll(strings.ToLower(tableColumn), "_", "") // 忽略大小写，去除_
		inputColumn = strings.ReplaceAll(strings.ToLower(inputColumn), "_", "")
		if strings.Compare(column, inputColumn) == 0 {
			return tableColumn
		}
	}
	return ""
}
