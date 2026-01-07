# Salesforce MCP Server (v5.0.0) システム詳細設計書 (System Design Document)

## 1. プロジェクト概要 (Overview)

- **目的**: Model Context Protocol (MCP) サーバーとして動作し、Salesforce REST APIへの接続および拡張業務機能（日報、GPS、音声）を提供する。
- **対象ユーザー**: Claude Desktop等のMCPクライアント利用者、Salesforceを使用するフィールドセールス。
- **主要機能**: OAuth 2.0認証、オブジェクトスキーマ取得(Describe)、**日報作成（GPS/音声対応）**。

## 2. ドメイン設計 (Domain Design / DDD)

### 2.1 境界づけられたコンテキスト (Bounded Contexts)

- **Auth Context**: OAuth 2.0認証、トークン管理、PKCE。
- **Salesforce Core Context**: スキーマ取得 (Describe)、汎用REST操作。
- **Nippou Context (New)**: 日報管理、位置情報特定、音声データ連携。

### 2.2 ドメインモデル (Domain Models)

```mermaid
classDiagram
    class Nippou {
        +NippouId id
        +Date date
        +Content text
        +Location location
        +VoiceData voice
        +create()
        +attachVoice()
    }
    class Location {
        +Latitude lat
        +Longitude lng
        +Address address
        +distanceTo(other)
    }
    class VoiceConfig {
        +ModelName model
        +AutoSpeaking bool
    }
    Nippou *-- Location
    Nippou *-- VoiceConfig
```

## 3. システムアーキテクチャ (Clean Architecture)

### 3.1 レイヤー構成

- **Domain Layer (Pure Go)**: `Nippou`, `TokenData` エンティティ。外部依存なし。
- **Application Layer (Use Cases)**: `CreateNippouUseCase`, `AuthFlowUseCase`.
- **Interface Layer (Adapters)**: `MCPHandler`.
- **Infrastructure Layer**: `SalesforceClient`, `GoogleMapsClient`.

### 3.2 構成図 (Architecture Diagrams)

#### 物理構成図 (Physical Architecture - Original)
>
> **Note**: 元仕様書の物理構成・プロトコル詳細を維持します。

```
┌─────────────────────────────────────────────────────────────────┐
│                        MCP クライアント                          │
│                    (Claude Desktop等)                           │
└─────────────────────┬───────────────────────────────────────────┘
                      │ JSON-RPC 2.0 over stdio (UTF-8)
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                     MCP サーバー (Go)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ MCP Handler │  │ OAuth Mgr   │  │ Salesforce Client       │  │
│  │ - tools     │  │ - resources │  │ - REST API (HTTP再利用) │  │
│  │ - prompts   │  │ - PKCE      │  │ - describe (ETag対応)   │  │
│  │ - logging   │  │ - token mgr │  │ - circuit breaker       │  │
│  │ - nippou    │  │ - rotation  │  │ - composite API         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                      │ HTTPS (TLS 1.2+)
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Salesforce Platform                         │
│  ┌─────────────────────┐  ┌─────────────────────────────────┐   │
│  │ OAuth 2.0 Endpoints │  │ REST API                        │   │
│  │ - /authorize        │  │ - /sobjects (キャッシュ対応)     │   │
│  │ - /token            │  │ - /describe (ETag対応)          │   │
│  └─────────────────────┘  └─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

#### 論理構成図 (Logical Architecture - Clean Arch)
>
> **Note**: 上記コンポーネントをClean Architectureで再整理した図です。

```mermaid
graph TD
    Client[MCP Client] --> MCPHandler
    MCPHandler --> AuthUseCase
    MCPHandler --> NippouUseCase
    
    subgraph Domain Layer
        NippouUseCase --> NippouRepository[I.NippouRepository]
        NippouUseCase --> CustomerRepository[I.CustomerRepository]
    end
    
    subgraph Infrastructure Layer
        NippouRepository --> SalesforceImpl
        CustomerRepository --> SalesforceImpl
        SalesforceImpl --> SalesforceAPI
    end
```

## 4. API インターフェース仕様 (MCP & HTTP)

### 4.1 MCP Tools & HTTP Endpoints

| Tool Name / Path | Description | Arguments / Note |
|---|---|---|
| `auth_start` | Salesforce OAuth認証開始 | - |
| `list_sobjects` | オブジェクト一覧取得 | - |
| `nippou_create` | 日報作成 (GPS/音声) | `date`, `content`, `gps_lat`, `gps_lng`, `voice_mode` |
| `nippou_settings` | 音声設定管理 | `action` |
| `admin_geocode` | 住所座標変換バッチ | `target` |
| `/auth/callback` | OAuth Callback | HTTP GET |
| `/metrics` | Prometheus Metrics | HTTP GET |

## 5. データモデル (Data Models - Implementation Detail)

### 5.1 設定・データ構造体 (Go Structs)
>
> **Strict Retention**: 元仕様書の定義をそのまま維持します。

#### Config Structure

```go
type Config struct {
    Salesforce     SalesforceConfig     `json:"salesforce"`
    Server         ServerConfig         `json:"server"`
    TLS            TLSConfig            `json:"tls"`
    TokenStore     TokenStoreConfig     `json:"token_store"`
    Cache          CacheConfig          `json:"cache"`
    CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`
    Metrics        MetricsConfig        `json:"metrics"`
    Tracing        TracingConfig        `json:"tracing"`
    Logging        LoggingConfig        `json:"logging"`
}
```

#### TokenData Structure

```go
type TokenData struct {
    AccessToken              string `json:"access_token"`
    RefreshToken             string `json:"refresh_token"`
    TokenType                string `json:"token_type"`
    InstanceURL              string `json:"instance_url"`
    IssuedAt                 int64  `json:"issued_at"`
    ExpiresIn                int    `json:"expires_in"`
    Scope                    string `json:"scope"`
    UserID                   string `json:"user_id,omitempty"`
    OrgID                    string `json:"org_id,omitempty"`
}
```

#### NippouData (DTO)

```go
type NippouData struct {
    Date         string      `json:"date"`          // YYYY-MM-DD
    Content      string      `json:"content"`       // 日報本文
    Location     *Location   `json:"location"`      // GPS情報 (Optional)
    VoiceConfig  *VoiceConfig `json:"voice_config"` // 音声設定 (Optional)
    Tags         []string    `json:"tags"`          // 関連タグ
}
```

### 5.2 アルゴリズム・クエリ詳細

#### PKCE Generation

```go
func (g *PKCEGenerator) Generate() (*PKCEParams, error) {
    // S256チャレンジ生成
    verifierBytes := make([]byte, 64)
    rand.Read(verifierBytes)
    verifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
    hash := sha256.Sum256([]byte(verifier))
    challenge := base64.RawURLEncoding.EncodeToString(hash[:])
    return &PKCEParams{
        CodeVerifier: verifier, CodeChallenge: challenge, CodeChallengeMethod: "S256",
    }, nil
}
```

#### Salesforce Query (Neighborhood Search)

```sql
SELECT Id, Name, DEVICE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') dist
FROM Account
WHERE DISTANCE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') < 1
ORDER BY DISTANCE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') ASC
LIMIT 5
```

## 6. 非機能要件 (Non-Functional Requirements)

### 6.1 文字エンコーディング (厳守)

- **BOM検出・変換**: 入力ファイルのBOMを適切に処理する。
- **サロゲートペア検証**: 不正なユニコードシーケンスを検知する。
- **Unicode正規化**: **NFC** (Normalization Form C) を適用する。
- **Windowsコンソール対応**:
  - `SetConsoleOutputCP` (CP_UTF8) の自動実行を実装すること。

### 6.2 セキュリティ

- **暗号化**: Argon2 KDF + AES-256-GCM によるトークン暗号化。
- **HTTPヘッダー**:
  - `HSTS`
  - `CSP` (Content-Security-Policy)
  - `X-Frame-Options`
  - `X-Content-Type-Options`
- **プロキシ**: `X-Forwarded-For` の検証と信頼。

### 6.3 運用・監視

- **メトリクス**: Prometheus形式 (`/metrics`)。Salesforce APIコール数、キャッシュヒット率。
- **ログ**: 構造化JSONロギング (Structured JSON Logging)。

## 7. エラーハンドリング (Error Handling)

### 7.1 エラー構造体

```go
type AppError struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    Retryable bool   `json:"retryable"`
}
```

### 7.2 Circuit Breaker

```go
var ErrCircuitOpen = &AppError{ Code: 5001, Message: "サービスが一時的に利用できません", Retryable: true }
```

## 8. 移行・互換性・テスト

- **Language**: Go 1.21+
- **Lint**: `golangci-lint` (必須)
- **Build**: `make build-all` (Dockerfile / Multi-platform)
- **Sandbox**: `SF_IS_SANDBOX` 環境変数による切り替え。
