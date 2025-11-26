// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"

	"airdrop/internal/model"
	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	wallet := strings.TrimSpace(req.WalletAddr)
	if wallet == "" {
		return &types.CreateUserResp{
			BaseResp: *ErrorResp(400, "wallet_addr 不能为空"),
		}, nil
	}

	var user model.User
	err = l.svcCtx.DB.Where("wallet_addr = ?", wallet).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = model.User{
				WalletAddr: wallet,
				Nickname:   req.Nickname,
			}
			if err = l.svcCtx.DB.Create(&user).Error; err != nil {
				logx.Errorf("CreateUser create error: %v", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	resp = &types.CreateUserResp{
		BaseResp: *OkResp(),
		Data: types.UserInfo{
			Id:         int64(user.ID),
			WalletAddr: user.WalletAddr,
			Nickname:   user.Nickname,
		},
	}
	return resp, nil
}
