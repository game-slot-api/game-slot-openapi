package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config OpenAPI HTTP 客户端配置。
type Config struct {
	BaseURL   string
	AppSecret string
	Timeout   time.Duration
	HTTP      *http.Client
}

// Client 调用平台 /openapi/* 接口。
type Client struct {
	baseURL   string
	appSecret string
	http      *http.Client
}

func New(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	httpClient := cfg.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}
	return &Client{
		baseURL:   trimTrailingSlash(cfg.BaseURL),
		appSecret: cfg.AppSecret,
		http:      httpClient,
	}
}

func (c *Client) DecodeToken(ctx context.Context, req *DecodeSSOKeyRequest) (*DecodeSSOKeyReply, error) {
	var out DecodeSSOKeyReply
	if err := c.post(ctx, "/openapi/player/decode-token", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceReply, error) {
	var out GetBalanceReply
	if err := c.post(ctx, "/openapi/player/balance", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetGameBetConfig(ctx context.Context, req *GetGameBetConfigRequest) (*GetGameBetConfigReply, error) {
	var out GetGameBetConfigReply
	if err := c.post(ctx, "/openapi/history/get-game-bet-config", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetGameHistory(ctx context.Context, req *GetGameHistoryRequest) (*GetGameHistoryReply, error) {
	var out GetGameHistoryReply
	if err := c.post(ctx, "/openapi/history/get-game-history", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetPlayerRtp(ctx context.Context, req *GetPlayerRtpRequest) (*GetPlayerRtpReply, error) {
	var out GetPlayerRtpReply
	if err := c.post(ctx, "/openapi/rtp/get-player-rtp", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SelectSpin(ctx context.Context, req *SelectSpinRequest) (*SelectSpinReply, error) {
	var out SelectSpinReply
	if err := c.post(ctx, "/openapi/rtp/select-spin", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Bet(ctx context.Context, req *BetRequest) (*BetReply, error) {
	var out BetReply
	if err := c.post(ctx, "/openapi/game/bet", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Win(ctx context.Context, req *WinRequest) (*WinReply, error) {
	var out WinReply
	if err := c.post(ctx, "/openapi/game/win", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Refund(ctx context.Context, req *RefundRequest) (*RefundReply, error) {
	var out RefundReply
	if err := c.post(ctx, "/openapi/game/refund", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type envelope struct {
	Code      int             `json:"code"`
	Data      json.RawMessage `json:"data"`
	Msg       string          `json:"msg"`
	RequestId string          `json:"requestId"`
}

func (c *Client) post(ctx context.Context, path string, in any, out any) error {
	body, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + path
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("AppSecret", c.appSecret)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{
			HTTPStatus: resp.StatusCode,
			Message:    fmt.Sprintf("http %d: %s", resp.StatusCode, string(raw)),
			RawBody:    string(raw),
		}
	}

	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("unmarshal envelope: %w body=%s", err, string(raw))
	}
	if env.Code != 0 {
		return &APIError{
			Code:       env.Code,
			Message:    env.Msg,
			RequestID:  env.RequestId,
			HTTPStatus: resp.StatusCode,
			RawBody:    string(raw),
		}
	}
	if out == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		return fmt.Errorf("unmarshal data: %w body=%s", err, string(env.Data))
	}
	return nil
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
