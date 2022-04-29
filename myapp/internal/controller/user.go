package controller

import (
	"context"
	"myapp/api"
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/internal/service"
	"myapp/utility/utils"
)

var User = cUser{}

type cUser struct{}

func (a *cUser) Get(ctx context.Context, req *api.UserGetReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	getUser, err := service.User().GetUserById(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	if getUser == nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.UserNotFound, "uuid", req.Uuid)
	}
	if err = utils.MyCopy(ctx, res, getUser); err != nil {
		return nil, err
	}
	return
}

func (a *cUser) Delete(ctx context.Context, req *api.UserDeleteReq) (res *api.EmptyRes, err error) {
	rowsAffected, err := service.User().DeleteUserById(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 { // uuid唯一值，删除成功只会删除一条数据
		return nil, errorCode.NewMyErr(ctx, errorCode.UserNotFound, "uuid", req.Uuid)
	}
	return
}

func (a *cUser) Create(ctx context.Context, req *api.UserCreateReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	userCreateInput := &model.UserCreateInput{}
	utils.MyCopy(ctx, userCreateInput, req)
	userUuid, err := service.User().CreateUser(ctx, userCreateInput)
	if err != nil {
		return nil, err
	}
	userGetReq := &api.UserGetReq{Uuid: userUuid}
	return a.Get(ctx, userGetReq)
}

func (a *cUser) Update(ctx context.Context, req *api.UserUpdateReq) (res *api.UserGetRes, err error) {
	res = &api.UserGetRes{}
	userUpdateInput := &model.UserUpdateInput{}
	utils.MyCopy(ctx, userUpdateInput, req)
	rowsAffected, err := service.User().UpdateUser(ctx, userUpdateInput)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 1 { // 更新成功，应只会更新一条数据
		userGetReq := &api.UserGetReq{Uuid: req.Uuid}
		return a.Get(ctx, userGetReq)
	} else {
		return nil, errorCode.NewMyErr(ctx, errorCode.UserNotFound, "uuid", req.Uuid)
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
	utils.MyCopy(ctx, res, listUsers)
	return
}
