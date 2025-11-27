// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"airdrop/internal/entity"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type StartRoundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartRoundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartRoundLogic {
	return &StartRoundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartRoundLogic) StartRound(req *types.StartRoundRequest) (*types.RoundInfoResponse, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.RoundName) == "" {
		return nil, errors.New("round name required")
	}
	deadline := time.Unix(req.ClaimDeadline, 0)
	if deadline.Before(time.Now()) {
		return nil, errors.New("deadline must be in future")
	}
	round := &entity.AirdropRound{
		Name:          req.RoundName,
		MerkleRoot:    strings.ToLower(req.MerkleRoot),
		TokenAddress:  strings.ToLower(req.TokenAddress),
		ClaimDeadline: deadline,
		Status:        "active",
	}

	err := l.svcCtx.RunTx(l.ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&entity.AirdropRound{}).Where("status = ?", "active").Update("status", "archived").Error; err != nil {
			return err
		}
		if err := tx.Create(round).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	round.SnapshotAt = &now
	if err := l.svcCtx.SnapshotRound(l.ctx, round); err != nil {
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
