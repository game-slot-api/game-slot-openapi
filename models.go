package openapi

import (
	"encoding/base64"
	"fmt"
)

// DecodeSSOKeyRequest 解析进房 Token。
type DecodeSSOKeyRequest struct {
	Token string `json:"token"`
}

type DecodeSSOKeyReply struct {
	AppId          string  `json:"appId,omitempty"`
	PlayerId       string  `json:"playerId,omitempty"`
	GameBrand      string  `json:"gameBrand,omitempty"`
	GameId         string  `json:"gameId,omitempty"`
	Expire         int64   `json:"expire,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	CurrencySymbol string  `json:"currencySymbol,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
}

type GetBalanceRequest struct {
	AppId    string `json:"appId"`
	PlayerId string `json:"playerId"`
}

type GetBalanceReply struct {
	Currency       string  `json:"currency,omitempty"`
	CurrencySymbol string  `json:"currencySymbol,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
}

type GetGameBetConfigRequest struct {
	AppId     string `json:"appId"`
	GameBrand string `json:"gameBrand"`
	GameId    string `json:"gameId"`
}

type GetGameBetConfigReply struct {
	BetConfig  *GameBetConfig `json:"betConfig,omitempty"`
	ConfigJson string         `json:"configJson,omitempty"`
}

type GameBetConfig struct {
	MinBet      float64   `json:"minBet"`
	MaxBet      float64   `json:"maxBet"`
	DefaultBet  float64   `json:"defaultBet"`
	MaxWinLimit float64   `json:"maxWinLimit"`
	BetOptions  []float64 `json:"betOptions,omitempty"`
}

type GetPlayerRtpRequest struct {
	AppId     string `json:"appId"`
	PlayerId  string `json:"playerId"`
	GameBrand string `json:"gameBrand"`
	GameId    string `json:"gameId"`
}

type GetPlayerRtpReply struct {
	Rtp string `json:"rtp,omitempty"`
}

type SelectSpinRequest struct {
	AppId         string  `json:"appId"`
	PlayerId      string  `json:"playerId"`
	GameBrand     string  `json:"gameBrand"`
	GameId        string  `json:"gameId"`
	Bet           float64 `json:"bet"`
	RoundModel    string  `json:"roundModel,omitempty"` // DEFAULT / EXTRA / BUY
	Currency      string  `json:"currency,omitempty"`
	Mock          string  `json:"mock,omitempty"`
	RoundExtModel string  `json:"roundExtModel,omitempty"`
}

type SelectSpinReply struct {
	SpinId         int64   `json:"spinId,omitempty"`
	GameBrand      string  `json:"gameBrand,omitempty"`
	GameId         string  `json:"gameId,omitempty"`
	Bet            float64 `json:"bet,omitempty"`
	OriginBet      float64 `json:"originBet,omitempty"`
	OriginTotalWin float64 `json:"originTotalWin,omitempty"`
	Data           string  `json:"data,omitempty"` // 解压后 base64
	RoundModel     string  `json:"roundModel,omitempty"`
	TableName      string  `json:"tableName,omitempty"`
	Rtp            string  `json:"rtp,omitempty"`
	Rate           float64 `json:"rate,omitempty"`
	Compress       int32   `json:"compress,omitempty"`
}

// DecodeSpinData 将 SelectSpin 返回的 base64 Data 解码为原始字节。
func (r *SelectSpinReply) DecodeSpinData() ([]byte, error) {
	if r == nil || r.Data == "" {
		return nil, fmt.Errorf("empty spin data")
	}
	b, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return nil, fmt.Errorf("decode spin data: %w", err)
	}
	return b, nil
}

type BetRequest struct {
	AppId            string  `json:"appId"`
	PlayerId         string  `json:"playerId"`
	RoundId          string  `json:"roundId"`
	PreRoundId       string  `json:"preRoundId,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Bet              float64 `json:"bet"`
	GameBrand        string  `json:"gameBrand"`
	GameId           string  `json:"gameId"`
	Rtp              string  `json:"rtp,omitempty"`
	RoundModel       string  `json:"roundModel,omitempty"`
	GameData         []byte  `json:"gameData,omitempty"`
	TraceId          string  `json:"traceId,omitempty"`
	TransactionId    string  `json:"transactionId,omitempty"`
	IsFree           bool    `json:"isFree,omitempty"`
	PreTransactionId string  `json:"preTransactionId,omitempty"`
}

type BetReply struct {
	Currency string  `json:"currency,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
}

type WinRequest struct {
	AppId            string  `json:"appId"`
	PlayerId         string  `json:"playerId"`
	RoundId          string  `json:"roundId"`
	PreRoundId       string  `json:"preRoundId,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Bet              float64 `json:"bet,omitempty"`
	GameBrand        string  `json:"gameBrand"`
	GameId           string  `json:"gameId"`
	Win              float64 `json:"win"`
	BetTransactionId string  `json:"betTransactionId"`
	GameData         []byte  `json:"gameData,omitempty"`
}

type WinReply struct {
	Currency    string  `json:"currency,omitempty"`
	Balance     float64 `json:"balance"`
	HashBalance bool    `json:"hashBalance,omitempty"`
}

type RefundRequest struct {
	AppId            string  `json:"appId"`
	PlayerId         string  `json:"playerId"`
	RoundId          string  `json:"roundId"`
	PreRoundId       string  `json:"preRoundId,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Bet              float64 `json:"bet,omitempty"`
	GameBrand        string  `json:"gameBrand"`
	GameId           string  `json:"gameId"`
	BetTransactionId string  `json:"betTransactionId"`
}

type RefundReply struct {
	Status  string  `json:"status,omitempty"` // SETTLED / CANCELED / NOTFOUND / REFUNDING
	Balance float64 `json:"balance,omitempty"`
}
