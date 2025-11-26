// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"airdrop/internal/logic"
	"airdrop/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSnapshotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetSnapshotLogic(r.Context(), svcCtx)
		resp, err := l.GetSnapshot()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
