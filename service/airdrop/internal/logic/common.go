package logic

import (
	"context"
	"strconv"

	"airdrop/internal/types"
)

func OkResp() *types.BaseResp {
	return &types.BaseResp{Code: 0, Msg: "ok"}
}

func ErrorResp(code int64, msg string) *types.BaseResp {
	return &types.BaseResp{Code: code, Msg: msg}
}

// 帮助测试用：在 context 里存放路径参数 id
type ctxKeyID struct{}

func contextWithID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, ctxKeyID{}, strconv.FormatInt(id, 10))
}
