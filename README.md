# 🦅 Claw Template (Antigravity x Claude Code)

Antigravity（設計・Frontend）と Claude Code（Backend）が連携するための爆速開発環境テンプレートです。

## 🚀 使い方

### 1. 新しいプロジェクトを開始する

このディレクトリ一式を新しいプロジェクトフォルダとしてコピーしてください。

### 2. セットアップ (Claw CLI)

以下のコマンドを実行すると、**対話型セットアップ**が起動します。

```bash
node setup_claw.js
```

セットアップ時に、開発スタイルに合わせてモードを選択できます：

* **🚀 1. Speed Vibe Mode (スピード優先)**
  * **用途**: プロトタイプ、PoC、ハッカソン。
  * **特徴**: 設計書作り込みを省略。チャットで「Vibe（雰囲気・やりたいこと）」を伝えて爆速実装。
* **🛡️ 2. Deep Dive Mode (詳細設計)**
  * **用途**: 本番環境、複雑なアプリ、チーム開発。
  * **特徴**: `design.md`（ユースケース、API仕様、UI遷移）を作成し、承認を得てから実装。

### 3. 開発フロー

選択したモードに応じて `claw.md` が自動生成されます。そのルールに従ってください。

1. **Kickoff**: Antigravityに作りたいものを伝えます。
2. **Design**: (Deep Dive Modeのみ) 詳細な設計書を作成・承認します。
3. **Claw Cycle**:
   * **Antigravity**: 設計、Frontend、全体の監督。
   * **Claude Code**: Backend実装、API構築。
   * **Reverse Sync**: 実装コードの変更を仕様書へ逆反映。

## 📂 構成

- `claw.md`: 連携ルールとワークフロー定義 (自動生成)
* `setup_claw.js`: 対話型セットアップCLI
* `antigravity.json` / `claude.json`: MCP連携設定
