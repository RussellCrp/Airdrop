// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserScoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserScoreLogic {
	return &GetUserScoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserScoreLogic) GetUserScore() (resp *types.ScoreResp, err error) {
	raw := l.ctx.Value(ctxKeyID{})
	idStr, ok := raw.(string)
	if !ok {
		return &types.ScoreResp{
			BaseResp: *ErrorResp(400, "missing user id"),
		}, nil
	}
	uid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return &types.ScoreResp{
			BaseResp: *ErrorResp(400, "invalid user id"),
		}, nil
	}

	score, tier, err := getUserScoreAndTier(l.svcCtx, uid)
	if err != nil {
		return nil, err
	}

	resp = &types.ScoreResp{
		BaseResp: *OkResp(),
		Data: types.ScoreData{
			Score: score,
			Tier:  tier,
		},
	}
	return resp, nil
}
