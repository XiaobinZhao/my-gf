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

var User = cUser{}

type cUser struct{}

func (a *cUser) Get(ctx context.Context, req *api.UserGetReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	getUser, err := service.User().GetUserById(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}
	if getUser == nil {
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#userNotExists}`, req.UserUuid))
	}
	if err = copier.Copy(res, getUser); err != nil {
		return nil, err
	}
	return
}

func (a *cUser) Delete(ctx context.Context, req *api.UserDeleteReq) (res *api.EmptyRes, err error) {
	rowsAffected, err := service.User().DeleteUserById(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 { // uuid唯一值，删除成功只会删除一条数据
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#userNotExists}`, req.UserUuid))
	}
	return
}

func (a *cUser) Create(ctx context.Context, req *api.UserCreateReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	userCreateInput := &model.UserCreateInput{}
	copier.Copy(userCreateInput, req)
	userUuid, err := service.User().CreateUser(ctx, userCreateInput)
	if err != nil {
		return nil, err
	}
	userGetReq := &api.UserGetReq{UserUuid: userUuid}
	return a.Get(ctx, userGetReq)
}

func (a *cUser) Update(ctx context.Context, req *api.UserUpdateReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	userUpdateInput := &model.UserUpdateInput{}
	copier.Copy(userUpdateInput, req)
	rowsAffected, err := service.User().UpdateUser(ctx, userUpdateInput)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 1 { // 更新成功，应只会更新一条数据
		userGetReq := &api.UserGetReq{UserUuid: req.UserUuid}
		return a.Get(ctx, userGetReq)
	} else {
		return nil, gerror.NewCode(errorCode.CodeNotFound, g.I18n().Tf(ctx, `{#userNotExists}`, req.UserUuid))
	}
}

func (a *cUser) List(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	res = &api.UserListRes{}
	sortKey := "createdAt" // 默认使用创建时间排序
	sortValue := "asc"
	if len(req.Sort) > 0 {
		if req.Sort[:1] == "-" {
			sortValue = "desc"
		}
		sortKey = req.Sort[1:]
	}
	listUsers, err := service.User().QueryUsers(ctx, model.UserListInput{
		SearchStr: req.SearchStr,
		Page:      req.Page,
		Size:      req.Size,
		SortKey:   sortKey,
		SortValue: sortValue,
	})
	if err != nil {
		return nil, err
	}
	copier.Copy(res, listUsers)
	return
}
