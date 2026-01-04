package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log-system-backend/application/auth/internal/repository"
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

func NewAuthService(repo repository.UserRepository, accessSecret string, accessExpire int64) AuthService {
	return &authService{
		repo:         repo,
		accessSecret: accessSecret,
		accessExpire: accessExpire,
	}
}

func (s *authService) Register(ctx context.Context, username, password string) (string, error) {
	// Check if user exists
	_, err := s.repo.FindOneByUsername(ctx, username)
	if err == nil {
		return "", errorx.NewCodeError(errorx.CodeParamError, "username already exists")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		fmt.Printf("find user error: %v\n", err)
		return "", errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeInternal, "failed to hash password")
	}

	user := &repository.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.repo.Insert(ctx, user)
	if err != nil {
		fmt.Printf("insert user error: %v\n", err)
		return "", errorx.NewCodeError(errorx.CodeInternal, "failed to create user")
	}

	return user.ID, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", errorx.NewCodeError(errorx.CodeAuthError, "invalid credentials")
		}
		return "", errorx.NewCodeError(errorx.CodeInternal, "database error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeAuthError, "invalid credentials")
	}

	// Generate JWT
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + s.accessExpire
	claims["iat"] = now
	claims["userId"] = user.ID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", errorx.NewCodeError(errorx.CodeInternal, "failed to generate token")
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (string, bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
