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

func (a *cUser) List(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	res = &api.UserListRes{}
	sortKey := ""
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
	res.Size = listUsers.Size
	res.Total = listUsers.Total
	res.Page = listUsers.Page
	res.List = []*api.UserGetRes{} // 默认为空数组而不是nil
	for _, user := range listUsers.List {
		userGetRes := &api.UserGetRes{}
		copier.Copy(userGetRes, user.User)
		res.List = append(res.List, userGetRes)
	}
	return
}
