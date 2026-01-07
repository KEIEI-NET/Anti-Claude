# Salesforce MCP Server プログラム仕様書

## 1. 概要 (Overview)

### 1.1 目的

本システムは、Model Context Protocol（MCP）サーバーとして動作し、OAuth 2.0認証を介してSalesforce REST APIに接続し、オブジェクトスキーマ（describe）情報を取得・提供することを目的とする。
v5.0.0より、業務特化型機能として「日報 (Nippou)」の作成・管理機能を強化し、GPS位置情報管理および音声入力インターフェースへの対応を追加する。

### 1.2 スコープ

- **ターゲット**: Go 1.21+
- **対象機能**: Salesfoce REST API連携、MCPプロトコル実装 (2024-11-05版)、日報管理機能。
- **バージョン**: v5.0.0

## 2. アーキテクチャ設計 (Architecture Design)

### 2.1 構成図

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
│  │ - tools     │  │ - PKCE      │  │ - REST API (HTTP再利用) │  │
│  │ - resources │  │ - token mgr │  │ - describe (ETag対応)   │  │
│  │ - prompts   │  │ - refresh   │  │ - circuit breaker       │  │
│  │ - logging   │  │ - rotation  │  │ - composite API         │  │
│  │ - nippou(New)│  │             │  │                         │  │
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
│  │ - /revoke           │  │ - /composite (バッチ処理)       │   │
│  └─────────────────────┘  └─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 依存関係 / 技術スタック

```go
require (
    golang.org/x/oauth2 v0.15.0           // OAuth 2.0 クライアント
    github.com/google/uuid v1.5.0         // UUID生成
    golang.org/x/crypto v0.16.0           // Argon2 KDF
    golang.org/x/text v0.14.0             // Unicode正規化
    go.opentelemetry.io/otel v1.21.0      // 分散トレーシング
    github.com/prometheus/client_golang v1.17.0  // メトリクス
)
```

## 3. インターフェース定義 (Interface Definition)

### 3.1 入力 (Inputs)

#### 環境変数設定

| 環境変数名 | 対応する設定 | 説明 |
|------------|--------------|------|
| `SF_CLIENT_ID` | salesforce.client_id | Consumer Key |
| `SF_CLIENT_SECRET` | salesforce.client_secret | Consumer Secret |
| `SF_REDIRECT_URI` | salesforce.redirect_uri | コールバックURL |
| `SF_API_VERSION` | salesforce.api_version | APIバージョン |
| `SF_IS_SANDBOX` | salesforce.is_sandbox | Sandbox環境フラグ |
| `GOOGLE_MAPS_API_KEY` | - | Geocoding API用 (New) |

### 3.2 出力 (Outputs) / エンドポイント

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/startup` | スタートアッププローブ |
| GET | `/ready` | レディネスプローブ |
| GET | `/health` | ライブネスプローブ |
| GET | `/metrics` | Prometheusメトリクス |
| GET | `/auth/start` | 認証開始 |
| GET | `/auth/callback` | OAuthコールバック |
| GET | `/auth/status` | 認証状態確認 |
| POST | `/auth/logout` | ログアウト |

### 3.3 MCPツール定義

#### 3.3.1 基本ツール

* `auth_start`: Salesforce OAuth認証開始
- `auth_status`: 認証状態確認
- `list_sobjects`: オブジェクト一覧取得
- `describe_sobject`: スキーマ情報取得
- `describe_sobject_fields`: フィールド情報取得

#### 3.3.2 日報拡張ツール (New)

* **`nippou_create`**:
  - 説明: 新しい日報を作成する。GPS情報や音声入力フラグをサポート。
  - 引数: `date`(string), `content`(string), `gps_lat`(number), `gps_lng`(number), `voice_mode`(boolean)
- **`nippou_voice_settings`**:
  - 説明: 日報入力時のデフォルト音声設定（モデル選択など）を取得・設定する。
  - 引数: `action`, `default_mode`
- **`admin_geocode_accounts`**:
  - 説明: 住所座標変換バッチ処理。
  - 引数: `target`, `account_id`

## 4. データモデル (Data Models)

### 4.1 共通設定 (Config)

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

### 4.2 認証トークン (TokenData)

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

### 4.3 日報データ (NippouData)

```go
type NippouData struct {
    Date         string      `json:"date"`          // YYYY-MM-DD
    Content      string      `json:"content"`       // 日報本文
    Location     *Location   `json:"location"`      // GPS情報 (Optional)
    VoiceConfig  *VoiceConfig `json:"voice_config"` // 音声設定 (Optional)
    Tags         []string    `json:"tags"`          // 関連タグ
}

type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Address   string  `json:"address,omitempty"` // 逆ジオコーディング結果など
}
```

## 5. 機能詳細 (Functional Details)

### 5.1 認証フロー (OAuth 2.0)

RFC 7636/9126準拠のPKCE認証を実装する。

#### PKCE生成

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

### 5.2 Salesforce API クライアント

#### バージョン管理

最小バージョン: `v55.0` 以上。

#### Describeキャッシュ

```go
func (c *DescribeCache) Get(objectName string) (*CacheEntry, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    entry, ok := c.entries[objectName]
    if !ok { return nil, false }
    if time.Since(entry.CachedAt) > c.ttl { return entry, false } // 期限切れだがETag用に返す
    return entry, true
}
```

#### Composite API (Batch)

最大25サブクエストまでのバッチ処理をサポート。

### 5.3 日報拡張機能 (Nippou Logic)

#### 5.3.1 GPS機能と顧客特定

モバイル利用および将来的な独自SFAへの移行を見据え、GPS座標は**「訪問先（関連顧客）」の特定**に活用する。

- **近隣顧客検索ロジック**:
    Salesforce SOQLの `DISTANCE` 関数を利用し、現在地から半径1km以内の顧客を検索する。

    ```sql
    SELECT Id, Name, DEVICE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') dist
    FROM Account
    WHERE DISTANCE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') < 1
    ORDER BY DISTANCE(BillingAddress, GEOLOCATION(:lat, :lng), 'km') ASC
    LIMIT 5
    ```

- **Geocoding (住所座標変換)**:
    `admin_geocode_accounts` ツールにより、Google Maps APIを利用して `Account` の住所から座標を補完・更新する。

#### 5.3.2 音声機能とデータ保全

将来的な移行（全データダウンロード）を考慮し、Salesforce標準の `ContentVersion` を利用する。

1. **入力**: クライアントから音声データ（Base64）とテキストを受け取る。
2. **保存**:
    - テキスト -> `Nippou__c.Content__c`
    - 音声バイナリ -> `ContentVersion`
3. **紐付け**: `ContentDocumentLink` を作成して両者をリンクする。

## 6. 非機能要件 (Non-Functional Requirements)

### 6.1 文字エンコーディング

* BOM検出・変換
- サロゲートペア検証
- Unicode正規化 (NFC)
- **Windowsコンソール対応**: `SetConsoleOutputCP` (CP_UTF8) の自動実行。

### 6.2 セキュリティ

* **暗号化**: Argon2 KDF + AES-256-GCM によるトークン暗号化。
- **HTTPヘッダー**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options。
- **プロキシ**: X-Forwarded-For の検証と信頼。

### 6.3 運用・監視

* **メトリクス**: Prometheus形式での公開 (`/metrics`)。Salesforce APIコール数、キャッシュヒット率等。
- **ログ**: 構造化JSONロギング。

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

Salesforce APIへの過負荷や連続エラー時に遮断を行う。

```go
var ErrCircuitOpen = &AppError{ Code: 5001, Message: "サービスが一時的に利用できません", Retryable: true }
```

## 8. テスト計画 (Test Plan)

### 8.1 テスト戦略

* **単体テスト**: 各パッケージ (`oauth`, `salesforce`, `nippou`) に対して `go test` を実施。
- **Lint**: `golangci-lint` を必須とする。

### 8.2 ビルド

MakefileおよびDockerfileを用意し、マルチプラットフォームビルドに対応する。

```bash
make build-all
```

## 9. 付録 (Appendix)

### 9.1 Salesforce Connected App 設定

* Callback URL: `http://localhost:8080/auth/callback`
- Scopes: `api`, `refresh_token`, `offline_access`
