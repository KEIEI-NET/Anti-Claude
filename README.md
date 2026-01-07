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

環境構築完了後、**`KICKOFF.md`** というファイルが生成されます。
このファイルを開き、**今の状況に合ったコマンド（プロンプト）をコピーして、チャットに貼り付けてください。**

* 例: 外部の仕様書を取り込むなら **「ケース 1」**
* 例: ゼロから設計するなら **「ケース 2」**

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
└── tools/
    └── normalize_docs.js        # 🛠️ 文字コード正規化ツール (Shift-JIS対策)
```

---

## 📝 機能詳細

### 📐 詳細設計テンプレート (Clean Arch / DDD)

生成されるテンプレートは、**Clean Architecture** のレイヤー構造（Domain, Application, Adapter, Infrastructure）および **DDD** の概念（Context Map, Aggregate）に対応しています。
Antigravityはこれに従い、フレームワークに依存しない純粋なドメインモデルを優先して設計します。

### 📝 プログラム仕様書 (Detailed Specs)

`program_spec_template.md` は、外注開発にも対応可能なレベルの詳細さを持ちます。

* クラス図 (Mermaid)
* インターフェース定義 (入出力DTO)
* 非機能要件、エラーハンドリング、テスト戦略

### 🛠️ 自動正規化ツール (Auto-Normalizer)

外部から持ち込まれたファイル（古いShift-JISのテキスト等）は、Antigravityによって自動的に検知・修正されます。

* `node tools/normalize_docs.js`
* 変換前のファイルは `.bak` としてバックアップされます。

---

## 🛠️ 2つの開発モード比較

| モード | 名称 | アーキテクチャ | 特徴 | ドキュメント |
| :--- | :--- | :--- | :--- | :--- |
| **Mode 1** | **🚀 Speed Vibe Mode** | 自由 (MVC等) | **速度最優先**。<br>動くものを最速で作る。 | 最小限。<br>走り書きレベル。 |
| **Mode 2** | **🛡️ Deep Dive Mode** | **Clean Architecture**<br>**DDD** | **堅牢性・保守性優先**。<br>SOLID原則遵守。<br>ドメインロジックの隔離。 | **厳格**。<br>詳細テンプレート準拠。 |

---

## 📈 開発ワークフロー (Mode 2)

### Phase 0: Domain Modeling

1. **Import**: `tools/normalize_docs.js` で既存資料を正規化。
2. **Analysis**: Antigravityが **DDD** に基づきドメインモデルを定義します。
3. **Spec**: `docs/design.md` (全体) と `docs/specs/[Component].md` (詳細) を作成。

### Phase 1: Implementation (Solid Principles)

1. **Domain Layer**: Claude Codeが依存関係を持たない純粋なドメインロジックを実装。
2. **Application Layer**: ユースケースの実装。
3. **Infrastructure**: DB接続やWebフレームワークの実装（DIによる注入）。

---
(c) 2026 KEIEI-NET
