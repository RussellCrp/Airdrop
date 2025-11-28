package handler

import (
	"airdrop/internal/svc"
	"context"
	"fmt"
)

const (
	InvestTaskCode = "INVEST"
	LoginTaskCode  = "LOGIN"
)

type TaskHandlerParams struct {
	TaskCode string
	Wallet   string
	Ctx      context.Context
	SvcCtx   *svc.ServiceContext
}

func NewTaskHandler(params *TaskHandlerParams) (ITaskHandler, error) {
	switch params.TaskCode {
	case LoginTaskCode:
		return &LoginTaskHandler{params: params}, nil
	case InvestTaskCode:
		return &InvestTaskHandler{params: params}, nil
	default:
		return nil, fmt.Errorf("unknown task code: %s", params.TaskCode)
	}
}

type ITaskHandler interface {
	Handle() error
}
