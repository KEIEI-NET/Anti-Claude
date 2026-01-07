# プログラム仕様書: Salesforce Infrastructure (Nippou)

## 1. 概要

Salesforce REST API を介して日報データの永続化と取得を行うコンポーネント。ドメイン層の `nippou.Repository` インターフェースを具象化する。

## 2. 構成ファイル (File Splits)

- **client.go**: Salesforce 共通 HTTP クライアント。リトライ、OAuthヘッダー、エラーハンドリングを統括。
- **models.go**: Salesforce カスタムオブジェクト (`Nippou__c`) のデータマッピング。
- **nippou_repository.go**: 永続化ロジック本体。

## 3. 詳細仕様: client.go

### 3.1 認証管理

- 送出されるリクエスト全てに対し、Bearerトークンを自動付与。
- 401 Unauthorized 受信時にリトークンを試行、成功した場合は元のリクエストを1回リトライ（実装済み）。

### 3.2 エラー処理

- 構造化された Salesforce API エラーを `UseCaseError` に変換。
- タイムアウトはコンテキストにより管理。

## 4. 詳細仕様: models.go / nippou_repository.go

### 4.1 データマッピング (`Nippou__c`)

- `Date__c`: 日付 (YYYY-MM-DD)
- `Content__c`: 日報本文 (Long TextArea)
- `Geolocation__latitude__s`, `Geolocation__longitude__s`: 位置情報
- `Address__c`: 逆ジオコーディング結果

### 4.2 特記事項: SOQL クエリ

- `FindByID`: 指定された ID に基づく単一取得。
- `FindByDate`: 範囲指定による一覧取得。

## 5. テスト仕様

- `client_test.go` にて、`httptest` を用いたモックレスポンスによる以下の確認を実施：
  - 正常な保存フロー。
  - API 制限 (429) 発生時のエラー伝播。
  - 不正なトークン (401) 発生時の挙動。
