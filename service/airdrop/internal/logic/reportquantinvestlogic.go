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

type ReportQuantInvestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportQuantInvestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportQuantInvestLogic {
	return &ReportQuantInvestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportQuantInvestLogic) ReportQuantInvest(req *types.ReportTaskReq) (resp *types.BaseResp, err error) {
	var user model.User
	if err = l.svcCtx.DB.First(&user, req.UserId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrorResp(404, "user not found"), nil
		}
		return nil, err
	}

	investAmount := req.Value
	// 可选：解析数值做校验
	if investAmount != "" {
		if _, err := strconv.ParseFloat(investAmount, 64); err != nil {
			return ErrorResp(400, "invalid invest amount"), nil
		}
	}

	extra := map[string]any{
		"evidence": req.Evidence,
		"amount":   investAmount,
	}
	if err = completeTaskAndUpdateScore(l.svcCtx, uint64(req.UserId), taskCodeQuantInvest, extra); err != nil {
		return nil, err
	}
	return OkResp(), nil
}
