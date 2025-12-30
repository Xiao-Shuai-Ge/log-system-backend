package service

import (
	"context"

	"log-system-backend/common/errorx"
	"log-system-backend/common/rpc/logingester"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type LogApiService interface {
	WriteLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error
}

type logApiService struct {
	ingesterRpc logingester.LogIngester
}

func NewLogApiService(ingesterRpc logingester.LogIngester) LogApiService {
	return &logApiService{
		ingesterRpc: ingesterRpc,
	}
}

func (s *logApiService) WriteLog(ctx context.Context, source, level, content string, metadata map[string]interface{}) error {
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
		// 将 gRPC 错误转换为业务错误
		st, ok := status.FromError(err)
		if ok {
			return errorx.NewCodeError(int(st.Code()), st.Message())
		}
		return errorx.NewCodeError(errorx.CodeInternal, err.Error())
	}
	return nil
}
