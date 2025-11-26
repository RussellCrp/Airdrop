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

type ReportTradeVolumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportTradeVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportTradeVolumeLogic {
	return &ReportTradeVolumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportTradeVolumeLogic) ReportTradeVolume(req *types.ReportTaskReq) (resp *types.BaseResp, err error) {
	var user model.User
	if err = l.svcCtx.DB.First(&user, req.UserId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorResp(404, "user not found"), nil
		}
		return nil, err
	}

	if req.Value == "" {
		return ErrorResp(400, "trade volume required"), nil
	}
	vol, err := strconv.ParseFloat(req.Value, 64)
	if err != nil {
		return ErrorResp(400, "invalid trade volume"), nil
	}
	if vol < 10000 {
		return ErrorResp(400, "trade volume must be >= 10000"), nil
	}

	extra := map[string]any{
		"evidence": req.Evidence,
		"volume":   vol,
	}
	if err = completeTaskAndUpdateScore(l.svcCtx, uint64(req.UserId), taskCodeTrade10k, extra); err != nil {
		return nil, err
	}
	return OkResp(), nil
}
