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

// GetGameHistoryRequest 查询玩家游戏历史记录（平台透传 history GameHistoryList）。
type GetGameHistoryRequest struct {
	AppId        string `json:"appId"`
	GameBrand    string `json:"gameBrand"`
	GameId       string `json:"gameId"`
	PlayerId     string `json:"playerId"`
	Status       string `json:"status,omitempty"`       // 记录状态，空为全部
	StartTime    int64  `json:"startTime"`              // 开始时间戳（毫秒）
	EndTime      int64  `json:"endTime"`                // 结束时间戳（毫秒）
	PageIndex    int32  `json:"pageIndex,omitempty"`    // 分页索引
	PageSize     int32  `json:"pageSize,omitempty"`     // 每页数量
	NeedTotalNum bool   `json:"needTotalNum,omitempty"` // 是否返回总条数
	IsAsc        bool   `json:"isAsc,omitempty"`        // 是否升序（默认倒序）
	NoGameData   bool   `json:"noGameData,omitempty"`   // true 时不返回游戏数据
}

type GetGameHistoryReply struct {
	List      []*GetGameHistoryInfo `json:"list,omitempty"`
	PageIndex int32                 `json:"pageIndex,omitempty"`
	PageSize  int32                 `json:"pageSize,omitempty"`
	TotalNum  int64                 `json:"totalNum,omitempty"`
}

type GetGameHistoryInfo struct {
	Bet              float64 `json:"bet"`
	Win              float64 `json:"win"`
	RoundId          string  `json:"roundId"`
	OrderNo          string  `json:"orderNo"`
	PlayerId         string  `json:"playerId"`
	AppId            string  `json:"appId"`
	GameBrand        string  `json:"gameBrand"`
	GameId           string  `json:"gameId"`
	Sample           string  `json:"sample,omitempty"`
	Note             string  `json:"note,omitempty"`
	Rtp              string  `json:"rtp,omitempty"`
	Status           string  `json:"status,omitempty"`
	WinBalance       float64 `json:"winBalance,omitempty"`
	AfterBetBalance  float64 `json:"afterBetBalance,omitempty"`
	BeforeBetBalance float64 `json:"beforeBetBalance,omitempty"`
	CreateTime       int64   `json:"createTime,omitempty"` // 毫秒
	BetTime          int64   `json:"betTime,omitempty"`    // 毫秒
	WinTime          int64   `json:"winTime,omitempty"`    // 毫秒
	GameData         []byte  `json:"gameData,omitempty"`
	Compress         int32   `json:"compress,omitempty"`
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
	Bet           float64 `json:"bet"`                  // 实际下注额（可含倍率）
	RoundModel    string  `json:"roundModel,omitempty"` // DEFAULT / EXTRA / BUY
	Currency      string  `json:"currency,omitempty"`
	Mock          string  `json:"mock,omitempty"`
	RoundExtModel string  `json:"roundExtModel,omitempty"`
	// OriginBet 基础 bet（未乘倍率的样本基准下注）。
	// 有倍率/买免费等场景时：Bet 为玩家实际扣款额，OriginBet 为选局用的基础 bet；不传则由 rtp 侧按 Bet 处理。
	OriginBet float64 `json:"originBet,omitempty"`
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
	ControlTag     string  `json:"controlTag,omitempty"` // 控制标记，结算 win 时回传
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
	// ControlTag 控制标记，来自 SelectSpin 返回，结算时回传。
	ControlTag string `json:"controlTag,omitempty"`
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
