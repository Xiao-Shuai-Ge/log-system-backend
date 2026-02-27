package ctxutils

import (
	"context"
	"encoding/json"
	"log-system-backend/common/errorx"
)

// GetUserIdFromCtx 从上下文中获取用户 ID
func GetUserIdFromCtx(ctx context.Context) (string, error) {
	userId := ctx.Value("userId")
	var userIdStr string
	if v, ok := userId.(string); ok {
		userIdStr = v
	} else if v, ok := userId.(json.Number); ok {
		userIdStr = v.String()
	} else {
		return "", errorx.NewCodeError(errorx.CodeAuthError, "无效的用户 ID")
	}
	return userIdStr, nil
}

// GetAppCodeFromCtx 从上下文中获取应用代码
func GetAppCodeFromCtx(ctx context.Context) (string, error) {
	appCode := ctx.Value("appCode")
	if v, ok := appCode.(string); ok {
		return v, nil
	}
	return "", errorx.NewCodeError(errorx.CodeAuthError, "无效的应用代码")
}

// SetAppCodeToCtx 将应用代码设置到上下文中
func SetAppCodeToCtx(ctx context.Context, appCode string) context.Context {
	return context.WithValue(ctx, "appCode", appCode)
}
