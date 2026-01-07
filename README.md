# 🦅 Claw Template (Antigravity x Claude Code)

Antigravity（設計・Frontend）と Claude Code（Backend）が連携するための爆速開発環境テンプレートです。

## 🚀 使い方

### 1. 新しいプロジェクトを開始する

このディレクトリ一式を新しいプロジェクトフォルダとしてコピーしてください。

### 2. セットアップの確認

環境をリセットまたは初期化したい場合は、以下のコマンドを実行します。

```bash
node setup_claw.js
```

### 3. 開発フロー

詳細は `claw.md` を参照してください。

1. **キックオフ**: Antigravityに作りたいものを伝えます。
2. **仕様策定**: 壁打ちで仕様を決めるか、既存の設計書を渡します。
3. **開発開始**:
   - **Antigravity**: 設計、フロントエンド実装、全体の監督。
   - **Claude Code**: バックエンドAPI、ロジックの実装。

## 📂 構成

- `claw.md`: 連携ルールとワークフロー定義
- `setup_claw.js`: 環境セットアップスクリプト
- `antigravity.json` / `claude.json`: MCP連携設定
