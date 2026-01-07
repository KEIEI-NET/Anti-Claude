# Task: Implement Nippou Domain Layer

**Target**: Claude Code
**Priority**: High
**Deadline**: ASAP

## 概要

`internal/domain/nippou` パッケージの実装をお願いします。
これはClean Architectureのドメイン層（最重要部分）なので、**外部ライブラリへの依存（UUID生成など）は極力避けるか、Interface越しにしてください**（今回は簡易的にUUIDライブラリ使用可とします）。

## 参照資料

- 設計: `docs/design.md` (Quality Guard Passed)
- 詳細仕様: `docs/specs/nippou_spec.md`

## 実装項目

1. `internal/domain/nippou/nippou.go`: エンティティと構造体定義
2. `internal/domain/nippou/nippou_test.go`: 基本的な単体テスト

## 品質基準 (Definition of Done)

- [ ] `go test ./internal/domain/nippou/...` がパスすること。
- [ ] `docs/specs/nippou_spec.md` のデータ構造と一致していること。
- [x] Lintエラーがないこと。

## Status

- **2026-01-08**: Implemented by Claude Code.
- **2026-01-08**: Reviewed by Antigravity. **PASSED**.
