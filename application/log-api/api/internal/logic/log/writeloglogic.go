// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"
	"encoding/json"
	"fmt"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
	"log-system-backend/common/errorx"
	"log-system-backend/common/rpc/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWriteLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteLogLogic {
	return &WriteLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WriteLogLogic) WriteLog(req *types.WriteLogReq) (resp *types.WriteLogResp, err error) {
	userId := l.ctx.Value("userId")
	var userIdStr string
	if v, ok := userId.(string); ok {
		userIdStr = v
	} else if v, ok := userId.(json.Number); ok {
		userIdStr = v.String()
	} else {
		return nil, errorx.NewCodeError(errorx.CodeAuthError, "invalid user id")
	}

	// Verify access
	accessResp, err := l.svcCtx.AuthRpc.VerifyAppAccess(l.ctx, &auth.VerifyAppAccessRequest{
		UserId:  userIdStr,
		AppCode: req.Source,
	})
	if err != nil {
		return nil, err
	}
	if !accessResp.HasAccess {
		return nil, errorx.NewCodeError(errorx.CodeForbidden, fmt.Sprintf("no access to app: %s", req.Source))
	}

	err = l.svcCtx.LogApiService.WriteLog(l.ctx, req.Source, req.Level, req.Content, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &types.WriteLogResp{
		Ok: true,
	}, nil
}
