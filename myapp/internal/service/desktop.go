package service

import (
	"context"
	"myapp/api"
	"myapp/internal/model"
	"myapp/internal/model/entity"
	"myapp/internal/service/internal/dao"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/util/gconv"
)

type sDesktop struct{}

var insDesktop = sDesktop{}

func Desktop() *sDesktop {
	return &insDesktop
}

func (s *sDesktop) GetDesktopById(ctx context.Context, desktopId string) (out *model.DesktopGetOutput, err error) {
	if err = dao.Desktop.Ctx(ctx).WherePri(desktopId).Scan(&out); err != nil {
		return nil, err
	}
	// 需要判断nil是否存在,不存在需要判断为空,以防后续nil
	if out == nil {
		return nil, nil
	}
	return
}

func (s *sDesktop) DeleteDesktopById(ctx context.Context, desktopId string) (rowsAffected int64, err error) {
	result, err := dao.Desktop.Ctx(ctx).Where(dao.Desktop.Columns().Uuid, desktopId).Delete()
	if err != nil {
		return 0, err
	}
	rowsAffected, err = result.RowsAffected()

	return
}

func (s *sDesktop) CreateDesktop(ctx context.Context, input *model.DesktopCreateInput) (desktopUuid string, err error) {
	desktop := &entity.Desktop{}
	err = dao.Desktop.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if err := gconv.Struct(input, desktop); err != nil {
			return err
		}
		desktop.Uuid = guid.S()
		_, err := dao.Desktop.Ctx(ctx).Data(desktop).OmitEmpty().Insert()
		return err
	})
	return desktop.Uuid, err
}

func (s *sDesktop) UpdateDesktop(ctx context.Context, input *model.DesktopUpdateInput) (rowsAffected int64, err error) {
	err = dao.Desktop.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		result, err := dao.Desktop.Ctx(ctx).OmitEmpty().Data(input).
			FieldsEx(dao.Desktop.Columns().Uuid).
			Where(dao.Desktop.Columns().Uuid, input.Uuid).
			Update()
		if err != nil {
			return err
		}
		rowsAffected, err = result.RowsAffected()
		return err
	})
	return
}

func (s *sDesktop) QueryDesktops(ctx context.Context, input model.DesktopListInput) (out *model.DesktopListOutput, err error) {
	var (
		m           = dao.Desktop.Ctx(ctx)
		likePattern = `%` + input.SearchStr + `%`
	)
	out = &model.DesktopListOutput{}
	out.Page = input.Page
	out.Size = input.Size
	out.List = []api.DesktopGetRes{}

	// 模糊查询
	if len(input.SearchStr) > 0 {
		m = m.WhereLike(dao.Desktop.Columns().DisplayName, likePattern).
			WhereOrLike(dao.Desktop.Columns().Desc, likePattern)
	}
	// 分页
	listModel := m.Page(input.Page, input.Size)
	// 排序
	sortField := GetSortField(ctx, input.SortKey, dao.Desktop.Ctx(ctx).GetFieldsStr())
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
