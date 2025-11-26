package logic

import (
	"context"
	"testing"

	"airdrop/internal/types"
)

func TestCreateAndGetUser(t *testing.T) {
	ctx := newCtx(t)
	logic := NewCreateUserLogic(ctx, testSvcCtx)

	req := &types.CreateUserReq{
		WalletAddr: "0x00000000000000000000000000000000000000aa",
		Nickname:   "test-user",
	}
	resp, err := logic.CreateUser(req)
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("CreateUser resp code != 0, got %d", resp.Code)
	}

	// 查询用户
	getLogic := NewGetUserLogic(ctx, testSvcCtx)
	// 模拟路由里的 id 值
	ctxWithID := contextWithID(context.Background(), resp.Data.Id)
	getLogic.ctx = ctxWithID

	uResp, err := getLogic.GetUser()
	if err != nil {
		t.Fatalf("GetUser error: %v", err)
	}
	if uResp.Code != 0 {
		t.Fatalf("GetUser resp code != 0, got %d", uResp.Code)
	}
	if uResp.Data.User.WalletAddr != req.WalletAddr {
		t.Fatalf("wallet mismatch")
	}
}
