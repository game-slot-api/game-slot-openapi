package openapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDecodeTokenAndBet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/openapi/player/decode-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("AppSecret") != "secret" {
			t.Fatalf("missing AppSecret")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"appId":          "app1",
				"playerId":       "p1",
				"gameBrand":      "brand",
				"gameId":         "g1",
				"expire":         123,
				"currency":       "USD",
				"currencySymbol": "$",
				"balance":        100.5,
			},
			"requestId": "rid",
		})
	})
	mux.HandleFunc("/openapi/player/balance", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{
				"currency":       "USD",
				"currencySymbol": "$",
				"balance":        88.0,
			},
		})
	})
	mux.HandleFunc("/openapi/game/bet", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"msg":  "success",
			"data": map[string]any{"currency": "USD", "balance": 99.5},
		})
	})
	mux.HandleFunc("/openapi/rtp/select-spin", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 1001,
			"msg":  "internal server error",
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	cli := New(Config{BaseURL: srv.URL, AppSecret: "secret", Timeout: time.Second})

	token, err := cli.DecodeToken(t.Context(), &DecodeSSOKeyRequest{Token: "tok"})
	if err != nil {
		t.Fatal(err)
	}
	if token.PlayerId != "p1" || token.AppId != "app1" || token.CurrencySymbol != "$" {
		t.Fatalf("unexpected token reply: %+v", token)
	}

	bal, err := cli.GetBalance(t.Context(), &GetBalanceRequest{AppId: "app1", PlayerId: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	if bal.Balance != 88.0 || bal.CurrencySymbol != "$" {
		t.Fatalf("unexpected balance reply: %+v", bal)
	}

	bet, err := cli.Bet(t.Context(), &BetRequest{
		AppId: "app1", PlayerId: "p1", RoundId: "r1", Bet: 1, GameBrand: "brand", GameId: "g1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if bet.Balance != 99.5 {
		t.Fatalf("unexpected balance: %v", bet.Balance)
	}

	_, err = cli.SelectSpin(t.Context(), &SelectSpinRequest{AppId: "app1", PlayerId: "p1", GameBrand: "brand", GameId: "g1", Bet: 1})
	apiErr, ok := AsAPIError(err)
	if !ok || apiErr.Code != 1001 {
		t.Fatalf("expected APIError 1001, got %v", err)
	}
}
