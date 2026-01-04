package logic

import (
	"context"

	"log-system-backend/application/auth/rpc/auth"
	"log-system-backend/application/auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	userId, err := l.svcCtx.AuthService.Register(l.ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		UserId: userId,
	}, nil
}
