package utils

import (
	"context"
	"myapp/internal/errorCode"

	"github.com/jinzhu/copier"
)

func MyCopy(ctx context.Context, toValue interface{}, fromValue interface{}) (err error) {
	if err = copier.Copy(toValue, fromValue); err != nil {
		return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
	} else {
		return nil
	}
}
