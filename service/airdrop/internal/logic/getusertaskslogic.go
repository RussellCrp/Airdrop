// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"
	"time"

	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTasksLogic {
	return &GetUserTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserTasksLogic) GetUserTasks() (resp *types.UserTasksResp, err error) {
	raw := l.ctx.Value(ctxKeyID{})
	idStr, ok := raw.(string)
	if !ok {
		return &types.UserTasksResp{
			BaseResp: *ErrorResp(400, "missing user id"),
		}, nil
	}
	uid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return &types.UserTasksResp{
			BaseResp: *ErrorResp(400, "invalid user id"),
		}, nil
	}

	var tasks []model.Task
	if err = l.svcCtx.DB.Where("enabled = 1").Find(&tasks).Error; err != nil {
		return nil, err
	}

	var uts []model.UserTask
	if err = l.svcCtx.DB.Where("user_id = ?", uid).Find(&uts).Error; err != nil {
		return nil, err
	}
	utMap := make(map[uint64]model.UserTask)
	for _, ut := range uts {
		utMap[ut.TaskID] = ut
	}

	items := make([]types.UserTaskItem, 0, len(tasks))
	for _, t := range tasks {
		item := types.UserTaskItem{
			TaskId: int64(t.ID),
			Code:   t.Code,
			Name:   t.Name,
			Status: 0,
		}
		if ut, ok := utMap[t.ID]; ok {
			item.Status = int64(ut.Status)
			if ut.CompletedAt != nil && !ut.CompletedAt.IsZero() {
				item.CompletedAt = ut.CompletedAt.Format(time.RFC3339)
			}
		}
		items = append(items, item)
	}

	resp = &types.UserTasksResp{
		BaseResp: *OkResp(),
		Data:     items,
	}
	return resp, nil
}
