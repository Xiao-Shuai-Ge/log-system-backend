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
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// 验证连接
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned an error: %s", res.String())
	}

	return client, nil
}

func IndexJSON(ctx context.Context, client *elasticsearch.Client, index string, doc any, opts ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	body, err := jsonBody(doc)
	if err != nil {
		return nil, err
	}
	return client.Index(index, body, append([]func(*esapi.IndexRequest){
		client.Index.WithContext(ctx),
	}, opts...)...)
}

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

func ReadBody(resp *esapi.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil response")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("es error: status=%s body=%s", resp.Status(), string(b))
	}
	return b, nil
}

func ReadJSON(resp *esapi.Response, out any) error {
	b, err := ReadBody(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func jsonBody(v any) (io.Reader, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}
