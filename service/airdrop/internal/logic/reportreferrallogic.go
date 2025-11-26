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

type ReportReferralLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportReferralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportReferralLogic {
	return &ReportReferralLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportReferralLogic) ReportReferral(req *types.ReportTaskReq) (resp *types.BaseResp, err error) {
	var user model.User
	if err = l.svcCtx.DB.First(&user, req.UserId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorResp(404, "user not found"), nil
		}
		return nil, err
	}

	if req.Value == "" {
		return ErrorResp(400, "referral count required"), nil
	}
	cnt, err := strconv.Atoi(req.Value)
	if err != nil {
		return ErrorResp(400, "invalid referral count"), nil
	}
	if cnt < 10 {
		return ErrorResp(400, "referral count must be >= 10"), nil
	}

	extra := map[string]any{
		"evidence": req.Evidence,
		"count":    cnt,
	}
	if err = completeTaskAndUpdateScore(l.svcCtx, uint64(req.UserId), taskCodeReferral10, extra); err != nil {
		return nil, err
	}
	return OkResp(), nil
}
