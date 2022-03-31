package service

import (
	"context"
	"myapp/api"
	"myapp/internal/errorCode"
	"myapp/internal/model"
	"myapp/internal/model/entity"
	"myapp/internal/service/internal/dao"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/crypto/gmd5"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gerror"
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

// 检测给定的账号是否唯一
func (s *sUser) CheckLoginNameUnique(ctx context.Context, loginName string) error {
	g.Log().Infof(ctx, "===========:%s \n\n", ctx.Value("I18nLanguage").(string))
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().LoginName, loginName).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return gerror.NewCode(errorCode.LoginNameConflicted, g.I18n().Tf(ctx, `{#loginNameConflicted}`, loginName))
	}
	return nil
}

// 将密码按照内部算法进行加密:md5（登录名+密码）
func (s *sUser) EncryptPassword(loginName, password string) string {
	return gmd5.MustEncrypt(loginName + password)
}

func (s *sUser) CreateUser(ctx context.Context, input *model.UserCreateInput) (userUuid string, err error) {
	var user *entity.User
	err = dao.User.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := gconv.Struct(input, &user); err != nil {
			return err
		}
		if err := s.CheckLoginNameUnique(ctx, user.LoginName); err != nil {
			return err
		}
		user.Password = s.EncryptPassword(user.LoginName, user.Password)
		user.Uuid = guid.S()
		_, err := dao.User.Ctx(ctx).Data(user).OmitEmpty().Save()
		return err
	})
	return user.Uuid, err
}

// 根据请求的排序字符串获取数据库对应的字段
// 为什么要做字段对应转换？因为前端传参是根据response字段传的，response字段一般使用request model的json格式。
// 而json格式字段和数据库字段以及entity的字段都不一样，比如：json(loginName)和数据库字段(login_name)以及entity的字段(LoginName)
// 字段格式不同，导致字符串比较不一致。所以需要做一下忽略格式的对比，返回对应的数据库字段
func (s *sUser) getSortField(ctx context.Context, inputColumn string) string {
	tableColumnsStr := dao.User.Ctx(ctx).GetFieldsStr()
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
	sortField := s.getSortField(ctx, input.SortKey)
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
	//if err := all.ScanList(&out.List, "User"); err != nil {
	//	return nil, err
	//}
	if err := all.Structs(&out.List); err != nil {
		return nil, err
	}
	return out, nil
}
