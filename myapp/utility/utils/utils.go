package utils

import (
	"context"
	"myapp/internal/errorCode"

	"github.com/jinzhu/copier"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
)

func CfgGetIgnoreError(ctx context.Context, name string) *gvar.Var {
	gVar, _ := g.Config().Get(ctx, name)
	return gVar
}

func MyCopy(ctx context.Context, toValue interface{}, fromValue interface{}) (err error) {
	if err = copier.Copy(toValue, fromValue); err != nil {
		return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
	} else {
		return nil
	}
}
