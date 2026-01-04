package ctxutils

import (
	"context"
	"encoding/json"
	"log-system-backend/common/errorx"
)

func GetUserIdFromCtx(ctx context.Context) (string, error) {
	userId := ctx.Value("userId")
	var userIdStr string
	if v, ok := userId.(string); ok {
		userIdStr = v
	} else if v, ok := userId.(json.Number); ok {
		userIdStr = v.String()
	} else {
		return "", errorx.NewCodeError(errorx.CodeAuthError, "invalid user id")
	}
	return userIdStr, nil
}

func GetAppCodeFromCtx(ctx context.Context) (string, error) {
	appCode := ctx.Value("appCode")
	if v, ok := appCode.(string); ok {
		return v, nil
	}
	return "", errorx.NewCodeError(errorx.CodeAuthError, "invalid app code")
}

func SetAppCodeToCtx(ctx context.Context, appCode string) context.Context {
	return context.WithValue(ctx, "appCode", appCode)
}
