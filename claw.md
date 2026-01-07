# Claw - Antigravity & Claude Code 連携ルール

## 開発モード: 🛡️ Deep Dive Mode (Production Grade)

## 1. コンテキスト同期
- Antigravity と Claude Code は同一のプロジェクトルートを共有します。
- 両ツールは MCP (Model Context Protocol) を使用して共有状態にアクセスします。

## 2. 役割と責任

### 🧠 Antigravity (設計 & フロントエンド)
- **設計 & ドキュメント**:
  - **テンプレートシステム**: `.claw/templates/` 内のテンプレートとツールを活用すること。
  - **インポート対応**: 外部ドキュメント取り込み時は、必ず `tools/normalize_docs.js` でUTF-8化とバックアップを行うこと。
  - **文書構成**:
    1. **システム詳細設計書**: `docs/design.md`
    2. **プログラム仕様書**: `docs/specs/xxx.md`
  - **逆同期 (Reverse Sync)**: 実装変更時はドキュメントを即時更新すること。
- **フロントエンド開発**:
  - UI設計および実装を担当。
- **監督**:
  - バックエンドコードの厳格なレビュー。

### ⚡ Claude Code (バックエンド専門)
- **バックエンド開発**:
  - プログラム仕様書に基づいて実装を行う。
  - 制約事項: **仕様書の完全再現**

## 3. ワークフロー
### Phase 0: Detailed Architecture 🏛️
1. **Normalization**: Run `node tools/normalize_docs.js` to fix encoding of imported docs.
2. **System Spec**: Antigravity creates `docs/design.md` using `design_template.md`.
3. **Program Specs**: Create `docs/specs/[Name].md` using `program_spec_template.md`.
4. **Approval**: User MUST approve specs before coding starts.

### Phase 1: Structured Implementation
- **Frontend**: Antigravity implements strict component design.
- **Backend**: Claude Code implements API strictly following the Program Specs.

## 4. ステータス
- **MCP Status**: Active
- **Template System**: Enabled
