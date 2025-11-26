// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"airdrop/internal/merkle"
	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSnapshotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSnapshotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSnapshotLogic {
	return &CreateSnapshotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSnapshotLogic) CreateSnapshot(req *types.CreateSnapshotReq) (resp *types.CreateSnapshotResp, err error) {
	// 1. 读取所有有积分的用户
	var scores []model.UserScore
	if err = l.svcCtx.DB.Find(&scores).Error; err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return &types.CreateSnapshotResp{
			BaseResp: *ErrorResp(400, "no users to snapshot"),
		}, nil
	}

	// 2. 建立 snapshot 记录
	snap := model.AirdropSnapshot{
		Name:       req.Name,
		TokenAddr:  req.TokenAddress,
		TotalUsers: 0,
		Status:     0,
		// total_amount、merkle_root 稍后更新
	}
	if err = l.svcCtx.DB.Create(&snap).Error; err != nil {
		return nil, err
	}

	// 3. 为每个用户根据积分找到档位 & 空投数量，生成 items 和 leaves
	var items []model.AirdropSnapshotItem
	var leaves [][32]byte
	var totalAmountBig int64

	for idx, us := range scores {
		// 找到用户信息
		var user model.User
		if err = l.svcCtx.DB.First(&user, us.UserID).Error; err != nil {
			continue
		}
		// 通过 score_tiers 查 amount
		var tier model.ScoreTier
		if err = l.svcCtx.DB.
			Where("min_score <= ? AND max_score >= ?", us.TotalScore, us.TotalScore).
			First(&tier).Error; err != nil {
			continue
		}

		itemIdx := idx
		leaf := merkle.EncodeLeaf(uint64(itemIdx), user.WalletAddr, tier.AirdropAmount)
		items = append(items, model.AirdropSnapshotItem{
			SnapshotID: snap.ID,
			UserID:     user.ID,
			WalletAddr: user.WalletAddr,
			Idx:        itemIdx,
			Score:      us.TotalScore,
			Amount:     tier.AirdropAmount,
			LeafHash:   merkle.ToHex32(leaf),
		})
		leaves = append(leaves, leaf)
		// 为简单起见，总额使用 int64 累加（示例用）
		fmt.Sscan(tier.AirdropAmount, &totalAmountBig)
	}

	if len(items) == 0 {
		return &types.CreateSnapshotResp{
			BaseResp: *ErrorResp(400, "no snapshot items generated"),
		}, nil
	}

	if err = l.svcCtx.DB.Create(&items).Error; err != nil {
		return nil, err
	}

	// 4. 生成 Merkle Tree
	tree := merkle.Build(leaves)
	root, ok := tree.Root()
	if !ok {
		return nil, fmt.Errorf("failed to build merkle root")
	}

	// 5. 更新 snapshot 汇总信息
	snap.TotalUsers = len(items)
	snap.TotalAmount = fmt.Sprintf("%d", totalAmountBig)
	snap.MerkleRoot = merkle.ToHex32(root)
	snap.Status = 1
	if err = l.svcCtx.DB.Save(&snap).Error; err != nil {
		return nil, err
	}

	resp = &types.CreateSnapshotResp{
		BaseResp: *OkResp(),
		Data: types.SnapshotInfo{
			Id:           int64(snap.ID),
			Name:         snap.Name,
			MerkleRoot:   snap.MerkleRoot,
			TokenAddress: snap.TokenAddr,
			TotalUsers:   int64(snap.TotalUsers),
			TotalAmount:  snap.TotalAmount,
			Status:       int64(snap.Status),
		},
	}
	return resp, nil
}
