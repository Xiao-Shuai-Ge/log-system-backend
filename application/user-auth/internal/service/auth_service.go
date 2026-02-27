package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log-system-backend/application/user-auth/internal/repository"
	"log-system-backend/common/errorx"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, bool, error)
}

type authService struct {
	repo         repository.UserRepository
	accessSecret string
	accessExpire int64
}

// NewAuthService 创建身份认证服务实例
func NewAuthService(repo repository.UserRepository, accessSecret string, accessExpire int64) AuthService {
	return &authService{
		repo:         repo,
		accessSecret: accessSecret,
		accessExpire: accessExpire,
	}
}

// Register 注册新用户，对密码进行哈希处理并存入数据库
func (s *authService) Register(ctx context.Context, username, password string) (string, error) {
	// 检查用户是否存在
	_, err := s.repo.FindOneByUsername(ctx, username)
	if err == nil {
		return "", errorx.NewCodeError(errorx.CodeParamError, "用户名已存在")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return "", errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeInternal, "密码加密失败")
	}

	user := &repository.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.repo.Insert(ctx, user)
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeInternal, "创建用户失败")
	}

	return user.ID, nil
}

// Login 用户登录，验证密码并生成 JWT 令牌
func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", errorx.NewCodeError(errorx.CodeAuthError, "无效的凭证")
		}
		return "", errorx.NewCodeError(errorx.CodeInternal, "数据库错误")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeAuthError, "无效的凭证")
	}

	// 生成 JWT
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + s.accessExpire
	claims["iat"] = now
	claims["userId"] = user.ID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeInternal, "生成令牌失败")
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT 令牌的有效性并返回用户 ID
func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (string, bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return "", false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var userId string
		if v, ok := claims["userId"].(string); ok {
			userId = v
		}
		return userId, true, nil
	}

	return "", false, nil
}
