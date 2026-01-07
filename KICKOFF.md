# 🚀 Claw Kickoff Prompts

環境セットアップ後、以下の手順でプロジェクトを開始してください。

## 📥 既存の仕様書がある場合 (Deep Import Flow)
**手順**:
1. プロジェクトルートにある `input_docs/` フォルダに、既存の資料を全て入れてください。
2. 以下のコマンドをチャットに貼り付けてください。

```text
@Antigravity
【プロジェクト開始: 既存仕様の完全インポート】

## 🚫 禁止事項 & 監視体制
- 私は **Quality Guard (tools/validate_docs.js)** を起動してあなたの成果物を監視します。
- 「省略」「要約」「図の欠落」があると、バリデーターがエラーを吐き、**作業完了と認められません**。
- 一発合格を目指して、一言一句漏らさず `docs/design.md` を作成してください。

## 📋 実行タスク
1. **正規化**: `node tools/normalize_docs.js` を実行。
2. **全量読込**: `input_docs/` を一字一句読み込む。
3. **品質検証**: Clean Architecture, DDD違反がないかチェック。
4. **統合と生成**: `.claw/templates/design_template.md` を使い、ダウングレードなしで生成。
   - **重要**: 最後に必ず `node tools/validate_docs.js` を実行し、合格すること。
```

## 🆕 新規開発の場合 (New Design Flow)
**手順**: 以下のコマンドをチャットに貼り付けてください。

```text
@Antigravity
【プロジェクト開始: 新規設計】
1. 詳細設計モード(Deep Dive)で進めます。
2. ヒアリング後、.claw/templates/design_template.md に基づいて docs/design.md を作成してください。
3. 最後に `node tools/validate_docs.js` を実行し、漏れがないか確認してください。
```

## ⚙️ プログラム仕様書の作成 (Implementation Prep)
```text
@Antigravity
【フェーズ移行: プログラム詳細設計】
1. docs/design.md のバリデーション(`node tools/validate_docs.js`)が通っていることを確認してください。
2. その後、docs/design.md に基づき、優先度の高いコンポーネントから順にプログラム仕様書を作成してください。
```
