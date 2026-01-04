// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"log-system-backend/application/log-api/api/internal/logic/log"
	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
)

func SearchLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchLogReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := log.NewSearchLogLogic(r.Context(), svcCtx)
		resp, err := l.SearchLog(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
