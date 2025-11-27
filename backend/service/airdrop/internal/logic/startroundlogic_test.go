package logic

import (
	"context"
	"strings"
	"testing"
	"time"

	"airdrop/internal/entity"
	"airdrop/internal/testutil"
	"airdrop/internal/types"
)

func TestStartRoundLogic(t *testing.T) {
	svcCtx := testutil.NewTestServiceContext(t)
	user := entity.User{
		Wallet:        strings.ToLower("0x1230000000000000000000000000000000000001"),
		PointsBalance: 800,
	}
	if err := svcCtx.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	req := &types.StartRoundRequest{
		RoundName:     "Genesis",
		MerkleRoot:    "0xabc",
		TokenAddress:  "0xdef",
		ClaimDeadline: time.Now().Add(1 * time.Hour).Unix(),
	}
	resp, err := NewStartRoundLogic(context.Background(), svcCtx).StartRound(req)
	if err != nil {
		t.Fatalf("start round: %v", err)
	}
	if resp.CurrentRoundId == 0 {
		t.Fatal("round id missing")
	}
	var snapshot entity.RoundPoint
	if err := svcCtx.DB.Where("round_id = ? AND user_id = ?", resp.CurrentRoundId, user.ID).First(&snapshot).Error; err != nil {
		t.Fatalf("round snapshot missing: %v", err)
	}
	if snapshot.Points != 800 {
		t.Fatalf("expected snapshot 800, got %d", snapshot.Points)
	}
	if err := svcCtx.DB.First(&user, user.ID).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if user.PointsBalance != 0 {
		t.Fatalf("expected user points zero after snapshot, got %d", user.PointsBalance)
	}
}
