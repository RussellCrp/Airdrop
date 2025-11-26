// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"sort"

	"airdrop/internal/merkle"
	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetSnapshotProofLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSnapshotProofLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSnapshotProofLogic {
	return &GetSnapshotProofLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSnapshotProofLogic) GetSnapshotProof(req *types.GetProofReq) (resp *types.GetProofResp, err error) {
	// 1. 找到 snapshot
	var snap model.AirdropSnapshot
	if err = l.svcCtx.DB.First(&snap, req.SnapshotId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.GetProofResp{
				BaseResp: *ErrorResp(404, "snapshot not found"),
			}, nil
		}
		return nil, err
	}

	// 2. 取出所有 items，按 idx 排序
	var items []model.AirdropSnapshotItem
	if err = l.svcCtx.DB.
		Where("snapshot_id = ?", snap.ID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &types.GetProofResp{
			BaseResp: *ErrorResp(400, "no items in snapshot"),
		}, nil
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Idx < items[j].Idx })

	// 3. 找到目标地址的 item
	var target *model.AirdropSnapshotItem
	var leafIndex int
	for i, it := range items {
		if it.WalletAddr == req.Address {
			target = &items[i]
			leafIndex = i
			break
		}
	}
	if target == nil {
		return &types.GetProofResp{
			BaseResp: *ErrorResp(404, "address not in snapshot"),
		}, nil
	}

	// 4. 重新构建 leaves
	leaves := make([][32]byte, len(items))
	for i, it := range items {
		leaf, err2 := merkle.FromHex32(it.LeafHash)
		if err2 != nil {
			return nil, err2
		}
		leaves[i] = leaf
	}

	tree := merkle.Build(leaves)
	proofHashes := tree.Proof(leafIndex)

	proofStr := make([]string, len(proofHashes))
	for i, p := range proofHashes {
		proofStr[i] = merkle.ToHex32(p)
	}

	resp = &types.GetProofResp{
		BaseResp: *OkResp(),
		Data: types.ProofData{
			Address: target.WalletAddr,
			Index:   int64(target.Idx),
			Amount:  target.Amount,
			Proof:   proofStr,
		},
	}
	return resp, nil
}
