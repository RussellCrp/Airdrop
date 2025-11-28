// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegistrationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegistrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegistrationLogic {
	return &RegistrationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegistrationLogic) Registration(req *types.RegistrationRequest) (resp *types.RegistrationResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
