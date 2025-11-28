// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"airdrop/internal/svc"
	"airdrop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitTaskLogic) SubmitTask(req *types.SubmitTaskRequest) (resp *types.SubmitTaskResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
