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

func NewLogApiService(ingesterRpc logingester.LogIngester, queryRpc logquery.LogQuery, authRpc auth.Auth) LogApiService {
	return &logApiService{
		ingesterRpc: ingesterRpc,
		queryRpc:    queryRpc,
		authRpc:     authRpc,
	}
}

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
		return errorx.NewCodeError(errorx.CodeForbidden, fmt.Sprintf("no access to app: %s", appCode))
	}
	return nil
}

func (s *logApiService) WriteLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
	if err := s.VerifyAccess(ctx, source); err != nil {
		return err
	}
	return s.writeToIngester(ctx, source, level, content, metadata)
}

func (s *logApiService) WriteAppLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
	// Skip VerifyAccess as it is for User context.
	// App context verification is done by middleware.
	return s.writeToIngester(ctx, source, level, content, metadata)
}

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
		return errorx.NewCodeError(errorx.CodeInternal, "failed to create log data")
	}

	_, err = s.ingesterRpc.WriteLog(ctx, &logingester.WriteLogReq{
		Data: logData,
	})
	if err != nil {
		return s.convertGrpcError(err)
	}
	return nil
}

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

func (s *logApiService) convertGrpcError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		return errorx.NewCodeError(int(st.Code()), st.Message())
	}
	return errorx.NewCodeError(errorx.CodeInternal, err.Error())
}
