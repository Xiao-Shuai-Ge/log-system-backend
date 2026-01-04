package log

import (
	"context"

	"log-system-backend/application/log-api/api/internal/svc"
	"log-system-backend/application/log-api/api/internal/types"
	"log-system-backend/common/ctxutils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushLogLogic {
	return &PushLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushLogLogic) PushLog(req *types.WriteLogReq) (resp *types.WriteLogResp, err error) {
	// 从 Context 中获取 appCode，确保日志来源是经过认证的应用
	appCode, err := ctxutils.GetAppCodeFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LogApiService.WriteAppLog(l.ctx, appCode, req.Level, req.Content, req.Metadata)
	if err != nil {
		return nil, err
	}
	return &types.WriteLogResp{
		Ok: true,
	}, nil
}
