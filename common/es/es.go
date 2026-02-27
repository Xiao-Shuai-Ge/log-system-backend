package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Config Elasticsearch 配置
type Config struct {
	Addresses          []string `json:",optional" yaml:",optional"`
	Username           string   `json:",optional" yaml:",optional"`
	Password           string   `json:",optional" yaml:",optional"`
	APIKey             string   `json:",optional" yaml:",optional"`
	CloudID            string   `json:",optional" yaml:",optional"`
	ServiceToken       string   `json:",optional" yaml:",optional"`
	InsecureSkipVerify bool     `json:",optional" yaml:",optional"`
	MaxRetries         int      `json:",optional" yaml:",optional"`
}

// NewClient 创建一个新的 Elasticsearch 客户端
func NewClient(cfg Config) (*elasticsearch.Client, error) {
	addresses := cfg.Addresses
	if len(addresses) == 0 {
		addresses = []string{"http://localhost:9200"}
	}

	esCfg := elasticsearch.Config{
		Addresses:    addresses,
		Username:     cfg.Username,
		Password:     cfg.Password,
		APIKey:       cfg.APIKey,
		CloudID:      cfg.CloudID,
		ServiceToken: cfg.ServiceToken,
	}

	if cfg.InsecureSkipVerify {
		esCfg.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if cfg.MaxRetries > 0 {
		esCfg.MaxRetries = cfg.MaxRetries
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("创建 Elasticsearch 客户端失败: %w", err)
	}

	// 验证连接
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("连接 Elasticsearch 失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch 返回错误: %s", res.String())
	}

	return client, nil
}

// IndexJSON 将文档索引到指定的 index 中
func IndexJSON(ctx context.Context, client *elasticsearch.Client, index string, doc any, opts ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	body, err := jsonBody(doc)
	if err != nil {
		return nil, err
	}
	return client.Index(index, body, append([]func(*esapi.IndexRequest){
		client.Index.WithContext(ctx),
	}, opts...)...)
}

// SearchJSON 在指定的 index 中执行搜索
func SearchJSON(ctx context.Context, client *elasticsearch.Client, index string, query any, opts ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	body, err := jsonBody(query)
	if err != nil {
		return nil, err
	}
	return client.Search(append([]func(*esapi.SearchRequest){
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(body),
	}, opts...)...)
}

// ReadBody 读取并关闭响应体，如果响应包含错误则返回错误
func ReadBody(resp *esapi.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("响应为空")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("Elasticsearch 错误: 状态=%s 响应体=%s", resp.Status(), string(b))
	}
	return b, nil
}

// ReadJSON 读取响应并将其反序列化为 JSON
func ReadJSON(resp *esapi.Response, out any) error {
	b, err := ReadBody(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// jsonBody 将对象转换为 JSON 的 io.Reader
func jsonBody(v any) (io.Reader, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}
