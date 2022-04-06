package controller

import (
	"context"
	"myapp/api"
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/internal/service"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/jinzhu/copier"
)

var Desktop = cDesktop{}

type cDesktop struct{}

func (a *cDesktop) Get(ctx context.Context, req *api.DesktopGetReq) (res *api.DesktopGetRes, err error) {
	res = &api.DesktopGetRes{}
	getDesktop, err := service.Desktop().GetDesktopById(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	if getDesktop == nil {
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#desktopNotExists}`, req.Uuid))
	}
	if err = copier.Copy(res, getDesktop); err != nil {
		return nil, err
	}
	return
}

func (a *cDesktop) Delete(ctx context.Context, req *api.DesktopDeleteReq) (res *api.EmptyRes, err error) {
	rowsAffected, err := service.Desktop().DeleteDesktopById(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 { // uuid唯一值，删除成功只会删除一条数据
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#desktopNotExists}`, req.Uuid))
	}
	return
}

func (a *cDesktop) Create(ctx context.Context, req *api.DesktopCreateReq) (res *api.DesktopGetRes, err error) {
	res = &api.DesktopGetRes{}
	userCreateInput := &model.DesktopCreateInput{}
	copier.Copy(userCreateInput, req)
	uuid, err := service.Desktop().CreateDesktop(ctx, userCreateInput)
	if err != nil {
		return nil, err
	}
	userGetReq := &api.DesktopGetReq{Uuid: uuid}
	return a.Get(ctx, userGetReq)
}

func (a *cDesktop) Update(ctx context.Context, req *api.DesktopUpdateReq) (res *api.DesktopGetRes, err error) {
	res = &api.DesktopGetRes{}
	userUpdateInput := &model.DesktopUpdateInput{}
	copier.Copy(userUpdateInput, req)
	rowsAffected, err := service.Desktop().UpdateDesktop(ctx, userUpdateInput)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 1 { // 更新成功，应只会更新一条数据
		userGetReq := &api.DesktopGetReq{Uuid: req.Uuid}
		return a.Get(ctx, userGetReq)
	} else {
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#desktopNotExists}`, req.Uuid))
	}
}

func (a *cDesktop) List(ctx context.Context, req *api.DesktopListReq) (res *api.DesktopListRes, err error) {
	res = &api.DesktopListRes{}
	sortKey := "createdAt" // 默认使用创建时间排序
	sortValue := "asc"
	if len(req.Sort) > 0 {
		if req.Sort[:1] == "-" {
			sortValue = "desc"
		}
		sortKey = req.Sort[1:]
	}
	listDesktops, err := service.Desktop().QueryDesktops(ctx, model.DesktopListInput{
		SearchStr: req.SearchStr,
		Page:      req.Page,
		Size:      req.Size,
		SortKey:   sortKey,
		SortValue: sortValue,
	})
	if err != nil {
		return nil, err
	}
	copier.Copy(res, listDesktops)
	return
}
