# 📋 開発実行計画書: Nippou UseCase 分割 & 品質向上

## 1. 現状の課題 (Current Issues)

- `internal/usecase/nippou/create.go` が 350行を超え、エラー定義・DTO・ロジックが混在している（God File）。
- 今後のエンドポイント拡張において、エラーコードや構造体の再利用が困難。

## 2. 目標構成 (Target Architecture)

`internal/usecase/nippou/` フォルダ内を以下の役割で分割する。

- `errors.go`: ユースケース層の共通エラー定義
- `dto.go`: 入出力リクエスト・レスポンス構造体
- `interactor.go`: ビジネスロジック実装（本体）
- `interactor_test.go`: ユニットテスト

## 3. 実装TODOリスト (Task Breakdown)

- [x] **Task 1: 定義層の分離**
  - [x] `errors.go` を作成し、`UseCaseError` 体系を移動。
  - [x] `dto.go` を作成し、`CreateInput`, `CreateOutput` 構造体を移動。
- [x] **Task 2: ロジックの再構築**
  - [x] `create.go` を `interactor.go` にリネームし、ロジックのみを保持。
  - [x] インポートパスの整合性を確認。
- [x] **Task 3: テストの移行と拡充**
  - [x] 既存の 44件のテストを `interactor_test.go` に整理。
  - [x] モックの定義が適切に分離されているか確認。
- [x] **Task 4: セキュリティ & 品質監査**
  - [x] Claude Code の Security Skills を発動し、分割後のコードを再監査。
- [x] **Task 5: 逆同期 (Reverse Sync)**
  - [x] 最終的な構成を `docs/design.md` に反映。

## 4. 完了定義 (Definition of Done)

- [ ] `go test ./internal/usecase/nippou/...` が 100% 合格すること。
- [ ] ファイルが3つ以上に適切に分割されていること。
- [ ] `docs/design.md` と実装が完全に同期していること。
