# プログラム仕様書: Claw Governance Validator

## 1. 概要
AI の作業品質を機械的に検証するゲートキーパー。

## 2. 監視ルール
- **Sync**: Package名に対応する spec ファイルが docs/specs/ に存在すること。
- **Coverage**: テストカバレッジ 80% 以上。
- **Plan**: PLAN.md に未完了項目 [ ] がないこと。
- **Arch**: ドメイン層が上位レイヤー（infrastructure等）に依存していないこと。
