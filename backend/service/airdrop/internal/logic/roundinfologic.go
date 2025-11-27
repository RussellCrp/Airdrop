// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"airdrop/internal/entity"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RoundInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoundInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoundInfoLogic {
	return &RoundInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoundInfoLogic) RoundInfo() (*types.RoundInfoResponse, error) {
	var round entity.AirdropRound
	err := l.svcCtx.DB.WithContext(l.ctx).Where("status = ?", "active").Order("id DESC").First(&round).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.RoundInfoResponse{}, nil
		}
		return nil, err
	}
	var total int64
	l.svcCtx.DB.WithContext(l.ctx).Model(&entity.RoundPoint{}).Where("round_id = ?", round.ID).Select("COALESCE(SUM(points),0)").Scan(&total)
	return &types.RoundInfoResponse{
		CurrentRoundId: int64(round.ID),
		RoundName:      round.Name,
		ClaimDeadline:  round.ClaimDeadline.Unix(),
		MerkleRoot:     round.MerkleRoot,
		TokenAddress:   round.TokenAddress,
		TotalPoints:    total,
	}, nil
}
