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

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser() (resp *types.GetUserResp, err error) {
	// id 从路由变量读取
	raw := l.ctx.Value(ctxKeyID{})
	idStr, ok := raw.(string)
	if !ok {
		return &types.GetUserResp{
			BaseResp: *ErrorResp(400, "missing user id"),
		}, nil
	}
	uid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return &types.GetUserResp{
			BaseResp: *ErrorResp(400, "invalid user id"),
		}, nil
	}

	var user model.User
	if err = l.svcCtx.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.GetUserResp{
				BaseResp: *ErrorResp(404, "user not found"),
			}, nil
		}
		return nil, err
	}

	score, tier, err := getUserScoreAndTier(l.svcCtx, uid)
	if err != nil {
		return nil, err
	}

	resp = &types.GetUserResp{
		BaseResp: *OkResp(),
		Data: types.GetUserData{
			User: types.UserInfo{
				Id:         int64(user.ID),
				WalletAddr: user.WalletAddr,
				Nickname:   user.Nickname,
			},
			Score: score,
			Tier:  tier,
		},
	}
	return resp, nil
}
