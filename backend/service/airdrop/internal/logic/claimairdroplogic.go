// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"airdrop/internal/authctx"
	"airdrop/internal/entity"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClaimAirdropLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimAirdropLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimAirdropLogic {
	return &ClaimAirdropLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimAirdropLogic) ClaimAirdrop(req *types.ClaimRequest) (*types.ClaimResponse, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if req.RoundId <= 0 {
		return nil, errors.New("invalid round id")
	}
	claims, ok := authctx.FromContext(l.ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	if !strings.EqualFold(claims.Wallet, req.Wallet) && req.Wallet != "" {
		return nil, errors.New("wallet mismatch")
	}
	roundID := uint64(req.RoundId)
	var round entity.AirdropRound
	if err := l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", roundID).First(&round).Error; err != nil {
		return nil, err
	}
	if round.Status != "active" {
		return nil, errors.New("round not active")
	}
	if time.Now().After(round.ClaimDeadline) {
		return nil, errors.New("claim deadline passed")
	}
	var snapshot entity.RoundPoint
	if err := l.svcCtx.DB.WithContext(l.ctx).Where("round_id = ? AND user_id = ?", roundID, claims.UID).First(&snapshot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no allocation for user")
		}
		return nil, err
	}
	remaining := snapshot.Points - snapshot.ClaimedPoints
	if req.Amount <= 0 || req.Amount > remaining {
		return nil, errors.New("amount exceeds allocation")
	}

	err := l.svcCtx.RunTx(l.ctx, func(tx *gorm.DB) error {
		newClaim := entity.Claim{
			RoundID: roundID,
			UserID:  claims.UID,
			Wallet:  claims.Wallet,
			Amount:  req.Amount,
			Status:  entity.ClaimStatusPending,
		}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "round_id"}, {Name: "wallet"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": req.Amount, "status": entity.ClaimStatusPending}),
		}).Create(&newClaim).Error
	})
	if err != nil {
		return nil, err
	}

	return &types.ClaimResponse{
		TxHash: "",
		Status: entity.ClaimStatusPending,
	}, nil
}
