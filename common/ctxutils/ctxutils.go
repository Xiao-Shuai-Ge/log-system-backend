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
