# 📋 開発実行計画書: Nippou Infrastructure (Salesforce Implementation)

## 1. 現状の課題 (Current Issues)

- `internal/domain/nippou/nippou.go` にインターフェース定義はあるが、実体がない。
- Salesforce API との疎通、認証、エラーハンドリング（リトライ等）を抽象化して実装する必要がある。

## 2. 目標構成 (Target Architecture)

`internal/infrastructure/salesforce/` フォルダを以下の役割で構成する。

- `client.go`: Salesforce REST API 共通の低レベル HTTP クライアント。
- `models.go`: `Nippou__c` などのカスタムオブジェクトに対応する JSON 変換用構造体定義。
- `nippou_repository.go`: `nippou.Repository` インターフェースの具象クラス実装。
- `client_test.go`: HTTP モックを使用した疎通テスト。

## 3. 実装TODOリスト (Task Breakdown)

### 🏗️ Phase 1: 共通クライアントの基盤構築

- [ ] **Task 1: Salesforce クライアントの実装**
  - [ ] `client.go` の作成（Auth Header, BaseURL, HTTP Client）。
  - [ ] 認可トークンのインジェクション機構。
- [ ] **Task 2: エラーマッピングの定義**
  - [ ] API エラー（401, 403, 429等）をドメイン/ユースケースエラーに変換する関数の作成。

### 🚀 Phase 2: リポジトリ実装

- [ ] **Task 3: Nippou__c マッピング**
  - [ ] `models.go` に Salesforce 側のカスタムフィールドと Go 構造体のマッピングを定義。
- [ ] **Task 4: Save メソッドの実装**
  - [ ] `nippou_repository.go` にて、ドメインエンティティを JSON に変換し、Salesforce へ POST するロジックを実装。
- [ ] **Task 5: Find メソッドの実装**
  - [ ] SOQL を使ったデータの取得（FindByID, FindByDate）の実装。

### 🛡️ Phase 3: 品質保証 & 同期

- [ ] **Task 6: 統合テスト (Mocked)**
  - [ ] Salesforce API のレスポンスをモックし、リポジトリ層のテストを完備する。
- [ ] **Task 7: 逆同期 (Reverse Sync)**
  - [ ] 最終的なインフラ構成を `docs/design.md` に反映。

## 4. 完了定義 (Definition of Done)

- [ ] `nippou.Repository` インターフェースが Salesforce 版で完全に実装されていること。
- [ ] SOQL などのクエリロジックが安全に実装されていること。
- [ ] 履歴バックアップ作成後、設計書が更新されていること。
