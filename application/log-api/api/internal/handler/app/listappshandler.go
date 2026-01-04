// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"log-system-backend/application/log-api/api/internal/logic/app"
	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
)

func ListAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAppsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := app.NewListAppsLogic(r.Context(), svcCtx)
		resp, err := l.ListApps(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
