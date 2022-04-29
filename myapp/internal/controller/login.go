package controller

import (
	"context"
	"myapp/api"
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/internal/service"
	"myapp/utility/token"
	"myapp/utility/utils"
)

var Login = cLogin{}

type cLogin struct{}

func (a *cLogin) Login(ctx context.Context, req *api.LoginReq) (res *api.LoginRes, err error) {
	res = &api.LoginRes{}
	userLoginInput := &model.UserLoginInput{}
	utils.MyCopy(ctx, userLoginInput, req)
	user, err := service.User().Login(ctx, userLoginInput)
	if err != nil {
		return nil, errorCode.MyWrapCode(ctx, errorCode.LoginFailed, err)
	}
	res.User = &api.UserGetRes{}
	if err = utils.MyCopy(ctx, res.User, user); err != nil {
		return nil, err
	}
	myCacheToken, err := token.Instance().GenerateToken(ctx, user.Uuid, user)
	if err != nil {
		return nil, err
	}
	res.Token = myCacheToken.Token
	return res, nil
}

func (a *cLogin) Logout(ctx context.Context, req *api.LogoutReq) (res *api.LogoutRes, err error) {
	user, err := service.User().GetUserById(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.UserNotFound, "uuid", req.Uuid)
	}
	err = token.Instance().RemoveCache(ctx, token.CacheKeyPrefix+req.Uuid)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
