package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LLMClient interface {
	Classify(ctx context.Context, content string) (safetyType, attackType string, err error)
}

type llmClient struct {
	baseURL string
	client  *http.Client
}

type llmRequest struct {
	Model    string       `json:"model"`
	Messages []llmMessage `json:"messages"`
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type classifyResult struct {
	SafetyType string `json:"safety_type"`
	AttackType string `json:"attack_type"`
}

func NewLLMClient(baseURL string, timeout time.Duration) LLMClient {
	return &llmClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *llmClient) Classify(ctx context.Context, content string) (string, string, error) {
	prompt := buildPrompt(content)
	reqBody := llmRequest{
		Model: "gpt-4",
		Messages: []llmMessage{
			{Role: "user", Content: prompt},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("大模型服务返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var llmResp llmResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return "", "", err
	}

	if len(llmResp.Choices) == 0 {
		return "", "", fmt.Errorf("大模型返回空结果")
	}

	var result classifyResult
	if err := json.Unmarshal([]byte(llmResp.Choices[0].Message.Content), &result); err != nil {
		return "", "", err
	}

	return result.SafetyType, result.AttackType, nil
}

func buildPrompt(content string) string {
	escaped := strings.ReplaceAll(content, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")

	if len(escaped) > 10000 {
		escaped = escaped[:10000]
	}

	return fmt.Sprintf(`面向大语言模型内容安全的攻击和防护是一个重要研究课题， 攻击者通过五花八门的精心设计的攻击手段，试图绕过大模型的安全防护， 从而引导大模型进行敏感类or危险性的内容回复, 对此，我需要全力分析其攻击特点，以便加强防护。 通常，攻击者的输入prompt由两个层面的内容组成：1）恶意问题本身；2）精心设计的攻击模板或范式，这种攻击模板或范式是一种套路性质的东西，可以被抽象归纳成一套攻击手段分类体系; （请尽量将恶意问题本身和攻击手段两个垂直角度剥离开来分析）

给你一段输入prompt，请输出：1）文本内容的安全分类；2）所采用的攻击手段分类, 分别对应输出dict中的"safety_type", "attack_type" 2个keys；

请按照分类说明，对用户的输入文本进行分类，用以下JSON格式输出:
{
    "safety_type": "文本内容的安全分类，多级标签用":"分割，中间不要有空格",
    "attack_type": "攻击手段分类, 如果不含有攻击手法，记为"无""
}

[用户输入开始]
%s
[用户输入结束]

只需返回JSON结果即可。
/no_think`, escaped)
}
