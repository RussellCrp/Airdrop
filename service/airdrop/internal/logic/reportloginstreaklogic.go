// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ReportLoginStreakLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportLoginStreakLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportLoginStreakLogic {
	return &ReportLoginStreakLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportLoginStreakLogic) ReportLoginStreak(req *types.ReportTaskReq) (resp *types.BaseResp, err error) {
	var user model.User
	if err = l.svcCtx.DB.First(&user, req.UserId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorResp(404, "user not found"), nil
		}
		return nil, err
	}

	if req.Value == "" {
		return ErrorResp(400, "login streak required"), nil
	}
	streak, err := strconv.Atoi(req.Value)
	if err != nil {
		return ErrorResp(400, "invalid login streak"), nil
	}
	if streak < 7 {
		return ErrorResp(400, "login streak must be >= 7"), nil
	}

	extra := map[string]any{
		"evidence": req.Evidence,
		"streak":   streak,
	}
	if err = completeTaskAndUpdateScore(l.svcCtx, uint64(req.UserId), taskCodeLogin7Days, extra); err != nil {
		return nil, err
	}
	return OkResp(), nil
}
