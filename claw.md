# Claw - Antigravity & Claude Code 連携ルール

## 開発モード: 🛡️ Deep Dive Mode (Production Grade)

## 1. コンテキスト同期
- Antigravity と Claude Code は同一のプロジェクトルートを共有します。
- 両ツールは MCP (Model Context Protocol) を使用して共有状態にアクセスします。

## 2. 役割と責任

### 🧠 Antigravity (設計 & フロントエンド)
- **設計 & ドキュメント**:
  - ドキュメント管理権限: **厳格 (Strict/Single Source of Truth)**
  - 詳細な設計書(design.md)を作成し、承認を得てから開発へ進みます。
  - **逆同期 (Reverse Sync)**: 必須 (コードの変更を仕様書へ反映)
- **フロントエンド開発**:
  - ユーザーインターフェース（UI）の設計と実装を行います。
- **監督**:
  - Claude Code が生成したバックエンドコードをレビューします。

### ⚡ Claude Code (バックエンド専門)
- **バックエンド開発**:
  - Antigravity の仕様に基づいて実装を行います。
  - 制約事項: **design.md の仕様を厳守すること。**

## 3. ワークフロー
### Phase 0: Detailed Architecture 🏛️
1. **Requirement Analysis**: Antigravity interviews User to define scope. (詳細ヒアリング)
2. **Specification**: Antigravity creates detailed `design.md`. (詳細設計書の作成: ER図, API, UIフロー)
3. **Approval**: User MUST approve `design.md` before any coding starts. (ユーザー承認後に着手)

### Phase 1: Structured Implementation
- **Frontend**: Antigravity implements strict component design.
- **Backend**: Claude Code implements API strictly following the Spec.

## 4. ステータス
- **MCP Status**: Active
- **Sync Status**: Verified
