package middleware

import (
	"net/http"

	"log-system-backend/common/ctxutils"
	"log-system-backend/common/errorx"
	"log-system-backend/common/rpc/auth"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AppAuthMiddleware struct {
	authRpc auth.Auth
}

func NewAppAuthMiddleware(authRpc auth.Auth) *AppAuthMiddleware {
	return &AppAuthMiddleware{
		authRpc: authRpc,
	}
}

func (m *AppAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appId := r.Header.Get("X-App-Id")
		appSecret := r.Header.Get("X-App-Secret")

		if appId == "" || appSecret == "" {
			httpx.Error(w, errorx.NewCodeError(401, "Missing App-Id or App-Secret header"))
			return
		}

		resp, err := m.authRpc.VerifyAppSecret(r.Context(), &auth.VerifyAppSecretRequest{
			AppId:     appId,
			AppSecret: appSecret,
		})
		if err != nil {
			httpx.Error(w, errorx.NewCodeError(500, "Internal Server Error during auth"))
			return
		}

		if !resp.IsValid {
			httpx.Error(w, errorx.NewCodeError(401, "Invalid App-Id or App-Secret"))
			return
		}

		// 将 appCode 放入 Context
		ctx := ctxutils.SetAppCodeToCtx(r.Context(), resp.AppCode)
		next(w, r.WithContext(ctx))
	}
}
