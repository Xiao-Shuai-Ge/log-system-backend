package service

import (
	"context"
	"fmt"

	"log-system-backend/common/ctxutils"
	"log-system-backend/common/errorx"
	"log-system-backend/common/rpc/auth"
	"log-system-backend/common/rpc/logingester"
	"log-system-backend/common/rpc/logquery"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type LogApiService interface {
	WriteLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error
	WriteAppLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error
	SearchLog(ctx context.Context, source, keyword string, metadata map[string]string, page, pageSize int64) (*logquery.SearchLogResp, error)
}

type logApiService struct {
	ingesterRpc logingester.LogIngester
	queryRpc    logquery.LogQuery
	authRpc     auth.Auth
}

// NewLogApiService 创建日志 API 服务实例
func NewLogApiService(ingesterRpc logingester.LogIngester, queryRpc logquery.LogQuery, authRpc auth.Auth) LogApiService {
	return &logApiService{
		ingesterRpc: ingesterRpc,
		queryRpc:    queryRpc,
		authRpc:     authRpc,
	}
}

// VerifyAccess 验证当前用户是否具有对指定应用的访问权限
func (s *logApiService) VerifyAccess(ctx context.Context, appCode string) error {
	userId, err := ctxutils.GetUserIdFromCtx(ctx)
	if err != nil {
		return err
	}

	accessResp, err := s.authRpc.VerifyAppAccess(ctx, &auth.VerifyAppAccessRequest{
		UserId:  userId,
		AppCode: appCode,
	})
	if err != nil {
		return err
	}
	if !accessResp.HasAccess {
		return errorx.NewCodeError(errorx.CodeForbidden, fmt.Sprintf("无权访问应用: %s", appCode))
	}
	return nil
}

// WriteLog 处理来自用户上下文的日志写入请求（需验证用户权限）
func (s *logApiService) WriteLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
	if err := s.VerifyAccess(ctx, source); err != nil {
		return err
	}
	return s.writeToIngester(ctx, source, level, content, metadata)
}

// WriteAppLog 处理来自应用上下文（AppSecret 验证）的日志写入请求
func (s *logApiService) WriteAppLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
	// 跳过 VerifyAccess，因为它是针对用户上下文的。
	// 应用上下文的验证由中间件完成。
	return s.writeToIngester(ctx, source, level, content, metadata)
}

// writeToIngester 内部方法，将日志数据发送到日志采集 RPC 服务
func (s *logApiService) writeToIngester(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
	data := make(map[string]interface{})

	if metadata != nil {
		for k, v := range metadata {
			data[k] = v
		}
	}

	// 核心字段最后赋值，确保其优先级最高，不会被 metadata 中的同名键覆盖
	data["source"] = source
	data["level"] = level
	data["content"] = content

	logData, err := structpb.NewStruct(data)
	if err != nil {
		return errorx.NewCodeError(errorx.CodeInternal, "创建日志数据失败")
	}

	_, err = s.ingesterRpc.WriteLog(ctx, &logingester.WriteLogReq{
		Data: logData,
	})
	if err != nil {
		return s.convertGrpcError(err)
	}
	return nil
}

// SearchLog 搜索指定应用的日志（需验证用户权限）
func (s *logApiService) SearchLog(ctx context.Context, source, keyword string, metadata map[string]string, page, pageSize int64) (*logquery.SearchLogResp, error) {
	if err := s.VerifyAccess(ctx, source); err != nil {
		return nil, err
	}
	rpcResp, err := s.queryRpc.SearchLog(ctx, &logquery.SearchLogReq{
		Source:   source,
		Keyword:  keyword,
		Metadata: metadata,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, s.convertGrpcError(err)
	}

	return rpcResp, nil
}

// convertGrpcError 将 gRPC 错误转换为自定义的 CodeError
func (s *logApiService) convertGrpcError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		return errorx.NewCodeError(int(st.Code()), st.Message())
	}
	return errorx.NewCodeError(errorx.CodeInternal, err.Error())
}
