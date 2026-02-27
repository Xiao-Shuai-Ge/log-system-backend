package logic

import (
	"context"

	"log-system-backend/application/user-auth/rpc/auth"
	"log-system-backend/application/user-auth/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 处理用户登录请求，验证身份并返回 JWT 令牌
func (l *LoginLogic) Login(in *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, err := l.svcCtx.AuthService.Login(l.ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}
