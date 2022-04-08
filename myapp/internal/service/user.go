package service

import (
	"context"
	"myapp/api"
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/internal/model/entity"
	"myapp/internal/service/internal/dao"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/crypto/gmd5"

	"github.com/gogf/gf/v2/util/gconv"
)

type sUser struct{}

var insUser = sUser{}

func User() *sUser {
	return &insUser
}

func (s *sUser) GetUserById(ctx context.Context, userId string) (out *model.UserGetOutput, err error) {
	if err = dao.User.Ctx(ctx).WherePri(userId).Scan(&out); err != nil {
		return nil, err
	}
	// 需要判断nil是否存在,不存在需要判断为空,以防后续nil
	if out == nil {
		return nil, nil
	}
	return
}

func (s *sUser) DeleteUserById(ctx context.Context, userId string) (rowsAffected int64, err error) {
	result, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Uuid, userId).Delete()
	if err != nil {
		return 0, err
	}
	rowsAffected, err = result.RowsAffected()

	return
}

// 检测给定的账号是否唯一
func (s *sUser) CheckLoginNameUnique(ctx context.Context, loginName string) error {
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().LoginName, loginName).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errorCode.NewMyErr(ctx, errorCode.LoginNameConflicted, loginName)
	}
	return nil
}

// 将密码按照内部算法进行加密:md5（登录名+密码）
func (s *sUser) EncryptPassword(loginName, password string) string {
	return gmd5.MustEncrypt(loginName + password)
}

func (s *sUser) CreateUser(ctx context.Context, input *model.UserCreateInput) (userUuid string, err error) {
	user := &entity.User{}
	err = dao.User.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := gconv.Struct(input, user); err != nil {
			return err
		}
		if err := s.CheckLoginNameUnique(ctx, user.LoginName); err != nil {
			return err
		}
		user.Password = s.EncryptPassword(user.LoginName, user.Password)
		user.Uuid = guid.S()
		_, err := dao.User.Ctx(ctx).Data(user).OmitEmpty().Insert()
		return err
	})
	return user.Uuid, err
}

func (s *sUser) UpdateUser(ctx context.Context, input *model.UserUpdateInput) (rowsAffected int64, err error) {
	err = dao.User.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if len(input.LoginName) > 0 {
			if err := s.CheckLoginNameUnique(ctx, input.LoginName); err != nil {
				return err
			}
		}
		result, err := dao.User.Ctx(ctx).OmitEmpty().Data(input).
			FieldsEx(dao.User.Columns().Uuid).
			Where(dao.User.Columns().Uuid, input.UserUuid).
			Update()
		if err != nil {
			return err
		}
		rowsAffected, err = result.RowsAffected()
		return err
	})
	return
}

func (s *sUser) QueryUsers(ctx context.Context, input model.UserListInput) (out *model.UserListOutput, err error) {
	var (
		m           = dao.User.Ctx(ctx)
		likePattern = `%` + input.SearchStr + `%`
	)
	out = &model.UserListOutput{}
	out.Page = input.Page
	out.Size = input.Size
	out.List = []api.UserGetRes{}

	// 模糊查询
	if len(input.SearchStr) > 0 {
		m = m.WhereLike(dao.User.Columns().DisplayName, likePattern).
			WhereOrLike(dao.User.Columns().LoginName, likePattern).
			WhereOrLike(dao.User.Columns().Phone, likePattern).
			WhereOrLike(dao.User.Columns().Email, likePattern).
			WhereOrLike(dao.User.Columns().Desc, likePattern)
	}
	// 分页
	listModel := m.Page(input.Page, input.Size)
	// 排序
	sortField := GetSortField(ctx, input.SortKey, dao.User.Ctx(ctx).GetFieldsStr())
	if len(input.SortKey) > 0 && sortField != "" {
		listModel = listModel.Order(sortField, input.SortValue)
	}
	all, err := listModel.All()
	if err != nil {
		return nil, err
	}
	if all.IsEmpty() {
		return out, nil
	}
	// total
	out.Total, err = m.Count()
	// 数据转换，db->model
	if err := all.Structs(&out.List); err != nil {
		return nil, err
	}
	return out, nil
}
