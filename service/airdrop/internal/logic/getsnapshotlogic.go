// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetSnapshotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSnapshotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSnapshotLogic {
	return &GetSnapshotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSnapshotLogic) GetSnapshot() (resp *types.GetSnapshotResp, err error) {
	idStr := l.ctx.Value("id").(string)
	sid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return &types.GetSnapshotResp{
			BaseResp: *ErrorResp(400, "invalid snapshot id"),
		}, nil
	}

	var snap model.AirdropSnapshot
	if err = l.svcCtx.DB.First(&snap, sid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.GetSnapshotResp{
				BaseResp: *ErrorResp(404, "snapshot not found"),
			}, nil
		}
		return nil, err
	}

	resp = &types.GetSnapshotResp{
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
