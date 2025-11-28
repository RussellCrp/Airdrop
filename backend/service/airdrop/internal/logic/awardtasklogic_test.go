package logic

import (
	"context"
	"testing"

	"airdrop/internal/testutil"
	"airdrop/internal/types"
)

func TestAwardTaskLogic(t *testing.T) {
	svcCtx := testutil.NewTestServiceContext(t)
	req := &types.AwardTaskRequest{
		Wallet: "0x9876543210000000000000000000000000000000",
		Task:   "promo",
		Amount: 0,
	}
	resp, err := NewAwardTaskLogic(context.Background(), svcCtx).AwardTask(req)
	if err != nil {
		t.Fatalf("award task: %v", err)
	}
	if resp.Data.Points != 1000 {
		t.Fatalf("expected 1000 points, got %d", resp.Data.Points)
	}
}
