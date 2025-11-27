package logic

import (
	"context"
	"strings"
	"testing"
	"time"

	"airdrop/internal/authctx"
	"airdrop/internal/entity"
	"airdrop/internal/security"
	"airdrop/internal/testutil"
	"airdrop/internal/types"
)

func TestClaimAirdropLogic(t *testing.T) {
	svcCtx := testutil.NewTestServiceContext(t)
	user := entity.User{
		Wallet:        strings.ToLower("0x1110000000000000000000000000000000000000"),
		PointsBalance: 0,
		FrozenPoints:  500,
	}
	if err := svcCtx.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	round := entity.AirdropRound{
		Name:          "Round1",
		MerkleRoot:    "0xroot",
		TokenAddress:  "0xtoken",
		ClaimDeadline: time.Now().Add(time.Hour),
		Status:        "active",
	}
	if err := svcCtx.DB.Create(&round).Error; err != nil {
		t.Fatalf("create round: %v", err)
	}
	rp := entity.RoundPoint{
		RoundID: round.ID,
		UserID:  user.ID,
		Points:  500,
	}
	if err := svcCtx.DB.Create(&rp).Error; err != nil {
		t.Fatalf("create round point: %v", err)
	}
	ctx := authctx.WithClaims(context.Background(), &security.Claims{
		UID:    user.ID,
		Wallet: user.Wallet,
		Role:   security.RoleUser,
	})
	req := &types.ClaimRequest{
		RoundId: int64(round.ID),
		Wallet:  user.Wallet,
		Amount:  200,
	}
	resp, err := NewClaimAirdropLogic(ctx, svcCtx).ClaimAirdrop(req)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if resp.Status != entity.ClaimStatusPending {
		t.Fatalf("expected pending status, got %s", resp.Status)
	}
	var claim entity.Claim
	if err := svcCtx.DB.Where("round_id = ? AND wallet = ?", round.ID, user.Wallet).First(&claim).Error; err != nil {
		t.Fatalf("claim not persisted: %v", err)
	}
	if claim.Amount != req.Amount {
		t.Fatalf("expected claim amount %d, got %d", req.Amount, claim.Amount)
	}
}
