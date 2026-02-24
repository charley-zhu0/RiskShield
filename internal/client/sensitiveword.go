package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charley/riskshield/internal/logger"
	"go.uber.org/zap"
)

type SensitiveWordClient interface {
	Match(ctx context.Context, content string) (bool, []string, error)
}

type sensitiveWordClient struct {
	baseURL string
	client  *http.Client
}

type sensitiveWordResponse struct {
	Query    string   `json:"query"`
	HitWords []string `json:"hit_words"`
}

func NewSensitiveWordClient(baseURL string, timeout time.Duration) SensitiveWordClient {
	return &sensitiveWordClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *sensitiveWordClient) Match(ctx context.Context, content string) (bool, []string, error) {
	url := fmt.Sprintf("%s/mod_sensitiveword/Match?tenantId=1", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(content))
	if err != nil {
		logger.Error("创建敏感词请求失败", zap.Error(err))
		return false, nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error("敏感词服务调用失败", zap.String("url", url), zap.Error(err))
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("敏感词服务返回错误", zap.Int("status", resp.StatusCode))
		return false, nil, fmt.Errorf("敏感词服务返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取敏感词响应失败", zap.Error(err))
		return false, nil, err
	}

	var result sensitiveWordResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("解析敏感词响应失败", zap.Error(err))
		return false, nil, err
	}

	return len(result.HitWords) > 0, result.HitWords, nil
}
