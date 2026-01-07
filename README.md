# 🦅 Claw Template (Antigravity x Claude Code)

**"Brain & Brawn"** — 最強のAIペアプログラミング環境を、たった1ファイルのスクリプトから。

Claw Templateは、**Antigravity（Architect/Frontend）** と **Claude Code（Backend）** をMCPで連携させ、爆速かつ高品質な開発を実現するための環境構築キットです。
**Clean Architecture** と **DDD (Domain-Driven Design)** を標準採用しており、長期運用に耐えうる堅牢な設計を自動化します。

---

## 🚀 使い方 (How to Use)

### 1. 準備

このリポジトリ全体をダウンロードする必要はありません。
**`setup_claw.js`** だけを、新しいプロジェクトフォルダにコピーしてください。

### 2. セットアップ実行

以下のコマンドを実行します。依存パッケージのインストールからテンプレート配置まで全自動で行われます。

```bash
node setup_claw.js
```

### 3. モード選択

実行すると対話メニューが表示されます。

```text
開発モードを選択してください:
[1] 🚀 Speed Vibe Mode (Prototyping)  ... 速度最優先。ざっくり指示で即実装。
[2] 🛡️ Deep Dive Mode (Clean Arch & DDD) ... 品質優先。詳細設計・承認・SOLID準拠。
```

> **推奨**: 本格的な開発では `[2]` を選択してください。

### 4. プロジェクト開始 (KICKOFF)

環境構築完了後、**`KICKOFF.md`** を開いてください。
用途に合わせた「開始時のプロンプト」が記載されています。これをAntigravityに送信してください。

| ケース | 用途 | 手順 |
| :--- | :--- | :--- |
| **既存資料あり** | 仕様書やコードを取り込む | 1. 資料を `input_docs/` に入れる<br>2. コマンドをコピペして送信 |
| **新規開発** | ゼロから作る | コマンドをコピペして送信するだけ |

---

## 📂 生成される環境 (Directory Structure)

```text
.
├── claw.md                      # 🦅 連携ルール定義書 (AIの行動指針)
├── KICKOFF.md                   # 🚀 開始コマンド集 (コピペ用プロンプト)
├── setup_claw.js                # ⚙️ セットアップスクリプト
├── antigravity.json / claude.json # 🔌 MCP設定
├── .claw/
│   └── templates/
│       ├── design_template.md        # 📄 システム詳細設計書 (Clean Arch/DDD構成)
│       └── program_spec_template.md  # 📄 プログラム仕様書詳細版 (クラス図/I・O定義)
├── docs/
│   └── specs/                   # 📁 設計書・仕様書の格納先
├── input_docs/                  # 📥 外部資料投入口 (ここに置くと一括読込される)
└── tools/
    └── normalize_docs.js        # 🛠️ 文字コード正規化ツール (Shift-JIS対策)
```

---

## 📝 機能詳細

### 📥 既存資料の一括インポート (`input_docs/`)

自動生成される `input_docs/` ディレクトリに、仕様書・メモ・ソースコードなどを放り込んでください。
KICKOFFの手順に従うことで、Antigravityは以下の処理を全自動で行います。

1. **正規化**: 文字コードを判別し、全てUTF-8に変換（バックアップ作成）。
2. **全読込**: フォルダ内の全ファイルをインプットとして学習。
3. **理解**: 内容を踏まえた上で、設計や実装の提案を行います。

### 📐 詳細設計テンプレート (Clean Arch / DDD)

生成されるテンプレートは、**Clean Architecture** のレイヤー構造（Domain, Application, Adapter, Infrastructure）および **DDD** の概念（Context Map, Aggregate）に対応しています。

### 📝 プログラム仕様書 (Detailed Specs)

`program_spec_template.md` は、外注開発にも対応可能なレベルの詳細さを持ちます。

* クラス図 (Mermaid)
* インターフェース定義 (入出力DTO)
* テスト戦略

---

## 🛠️ 2つの開発モード比較

| モード | 名称 | アーキテクチャ | 特徴 | ドキュメント |
| :--- | :--- | :--- | :--- | :--- |
| **Mode 1** | **🚀 Speed Vibe Mode** | 自由 (MVC等) | **速度最優先**。<br>動くものを最速で作る。 | 最小限。<br>走り書きレベル。 |
| **Mode 2** | **🛡️ Deep Dive Mode** | **Clean Architecture**<br>**DDD** | **堅牢性・保守性優先**。<br>SOLID原則遵守。<br>ドメインロジックの隔離。 | **厳格**。<br>詳細テンプレート準拠。 |

---
(c) 2026 KEIEI-NET
