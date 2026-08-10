# game-slot-openapi

平台 `platform_proxy_server` 的 `/openapi/*` HTTP 客户端，供游戏逆向服务复用。

## 安装

```bash
go get github.com/game-slot-api/game-slot-openapi@latest
```

## 使用

```go
import "github.com/game-slot-api/game-slot-openapi"

client := openapi.New(openapi.Config{
    BaseURL:   "https://platform.example.com",
    AppSecret: "your-secret",
})
```

鉴权请求头：`AppSecret`。统一响应 `code != 0` 时返回 `*openapi.APIError`。

## 接口

- `DecodeToken`（返回含 `currencySymbol`）
- `GetBalance`
- `GetGameBetConfig`
- `GetPlayerRtp`
- `SelectSpin`（`DecodeSpinData` 解码 base64）
- `Bet` / `Win` / `Refund`
