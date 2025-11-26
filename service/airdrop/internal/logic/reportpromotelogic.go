// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ReportPromoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportPromoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportPromoteLogic {
	return &ReportPromoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportPromoteLogic) ReportPromote(req *types.ReportTaskReq) (resp *types.BaseResp, err error) {
	var user model.User
	if err = l.svcCtx.DB.First(&user, req.UserId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorResp(404, "user not found"), nil
		}
		return nil, err
	}

	extra := map[string]any{
		"evidence": req.Evidence,
	}
	if err = completeTaskAndUpdateScore(l.svcCtx, uint64(req.UserId), taskCodePromote, extra); err != nil {
		return nil, err
	}
	return OkResp(), nil
}
