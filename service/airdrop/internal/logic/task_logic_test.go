package logic

import (
	"testing"

	"airdrop/internal/types"
)

func TestReportTasksAndScore(t *testing.T) {
	ctx := newCtx(t)

	// 使用 seed.sql 中的用户 1
	userID := int64(1)

	// 宣传任务
	if _, err := NewReportPromoteLogic(ctx, testSvcCtx).ReportPromote(&types.ReportTaskReq{
		UserId:   userID,
		Evidence: "link-1",
	}); err != nil {
		t.Fatalf("ReportPromote error: %v", err)
	}

	// 交易量任务（>=10000）
	if _, err := NewReportTradeVolumeLogic(ctx, testSvcCtx).ReportTradeVolume(&types.ReportTaskReq{
		UserId: userID,
		Value:  "10000",
	}); err != nil {
		t.Fatalf("ReportTradeVolume error: %v", err)
	}

	// 连续登录任务
	if _, err := NewReportLoginStreakLogic(ctx, testSvcCtx).ReportLoginStreak(&types.ReportTaskReq{
		UserId: userID,
		Value:  "7",
	}); err != nil {
		t.Fatalf("ReportLoginStreak error: %v", err)
	}

	// 查询积分
	scoreLogic := NewGetUserScoreLogic(ctx, testSvcCtx)
	ctxWithID := contextWithID(ctx, userID)
	scoreLogic.ctx = ctxWithID
	sResp, err := scoreLogic.GetUserScore()
	if err != nil {
		t.Fatalf("GetUserScore error: %v", err)
	}
	if sResp.Code != 0 {
		t.Fatalf("GetUserScore code != 0, got %d", sResp.Code)
	}
	if sResp.Data.Score <= 0 {
		t.Fatalf("expected positive score, got %d", sResp.Data.Score)
	}
}

