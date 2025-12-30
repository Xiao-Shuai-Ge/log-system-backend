package service

import (
	"context"

	"log-system-backend/common/rpc/logingester"

	"google.golang.org/protobuf/types/known/structpb"
)

type LogApiService interface {
	WriteLog(ctx context.Context, source, level, content string) error
}

type logApiService struct {
	ingesterRpc logingester.LogIngester
}

func NewLogApiService(ingesterRpc logingester.LogIngester) LogApiService {
	return &logApiService{
		ingesterRpc: ingesterRpc,
	}
}

func (s *logApiService) WriteLog(ctx context.Context, source, level, content string) error {
	data := map[string]interface{}{
		"source":  source,
		"level":   level,
		"content": content,
	}

	logData, err := structpb.NewStruct(data)
	if err != nil {
		return err
	}

	_, err = s.ingesterRpc.WriteLog(ctx, &logingester.WriteLogReq{
		Data: logData,
	})
	return err
}
