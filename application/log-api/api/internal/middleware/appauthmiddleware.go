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

// NewAppAuthMiddleware 创建应用授权中间件实例
func NewAppAuthMiddleware(authRpc auth.Auth) *AppAuthMiddleware {
	return &AppAuthMiddleware{
		authRpc: authRpc,
	}
}

// Handle 处理请求，通过请求头中的 X-App-Id 和 X-App-Secret 验证应用身份
func (m *AppAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appId := r.Header.Get("X-App-Id")
		appSecret := r.Header.Get("X-App-Secret")

		if appId == "" || appSecret == "" {
			httpx.Error(w, errorx.NewCodeError(401, "请求头中缺少 X-App-Id 或 X-App-Secret"))
			return
		}

		resp, err := m.authRpc.VerifyAppSecret(r.Context(), &auth.VerifyAppSecretRequest{
			AppId:     appId,
			AppSecret: appSecret,
		})
		if err != nil {
			httpx.Error(w, errorx.NewCodeError(500, "认证过程中服务器内部错误"))
			return
		}

		if !resp.IsValid {
			httpx.Error(w, errorx.NewCodeError(401, "无效的 X-App-Id 或 X-App-Secret"))
			return
		}

		// 将 appCode 放入 Context
		ctx := ctxutils.SetAppCodeToCtx(r.Context(), resp.AppCode)
		next(w, r.WithContext(ctx))
	}
}
