package logic

import (
	"testing"

	"airdrop/internal/types"
)

func TestSnapshotAndProof(t *testing.T) {
	ctx := newCtx(t)

	// 生成快照
	createLogic := NewCreateSnapshotLogic(ctx, testSvcCtx)
	snapResp, err := createLogic.CreateSnapshot(&types.CreateSnapshotReq{
		Name:         "test-snapshot",
		TokenAddress: "0x00000000000000000000000000000000000000ff",
	})
	if err != nil {
		t.Fatalf("CreateSnapshot error: %v", err)
	}
	if snapResp.Code != 0 {
		t.Fatalf("CreateSnapshot code != 0, got %d", snapResp.Code)
	}

	sid := snapResp.Data.Id

	// 查询快照
	getLogic := NewGetSnapshotLogic(ctx, testSvcCtx)
	getLogic.ctx = contextWithID(ctx, sid)
	gResp, err := getLogic.GetSnapshot()
	if err != nil {
		t.Fatalf("GetSnapshot error: %v", err)
	}
	if gResp.Code != 0 {
		t.Fatalf("GetSnapshot code != 0, got %d", gResp.Code)
	}
	if gResp.Data.MerkleRoot == "" {
		t.Fatalf("expected non-empty merkle root")
	}

	// 取出某个地址的 proof（使用 seed 里的用户 1）
	proofLogic := NewGetSnapshotProofLogic(ctx, testSvcCtx)
	pResp, err := proofLogic.GetSnapshotProof(&types.GetProofReq{
		SnapshotId: sid,
		Address:    "0x0000000000000000000000000000000000000001",
	})
	if err != nil {
		t.Fatalf("GetSnapshotProof error: %v", err)
	}
	if pResp.Code != 0 {
		t.Fatalf("GetSnapshotProof code != 0, got %d", pResp.Code)
	}
	if len(pResp.Data.Proof) == 0 && gResp.Data.TotalUsers > 1 {
		t.Fatalf("expected non-empty proof when many leaves")
	}
}

