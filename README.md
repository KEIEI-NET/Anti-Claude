# 🦅 Claw Template (Antigravity x Claude Code)

**"Brain & Brawn"** — 最強のAIペアプログラミング環境を、たった1ファイルのスクリプトから。

Claw Templateは、**Antigravity（設計・Frontend）** と **Claude Code（Backend）** をMCPで連携させ、爆速かつ高品質な開発を実現するための環境構築キットです。

---

## 🚀 使い方 (How to Use)

### 1. 準備

このリポジトリ全体をダウンロードする必要はありません。
**`setup_claw.js`** だけを、新しいプロジェクトフォルダにコピーしてください。

### 2. セットアップ実行

以下のコマンドを実行します。これだけで全ての環境が自動生成されます。

```bash
node setup_claw.js
```

### 3. モード選択

実行すると対話メニューが表示されます。用途に合わせて選択してください。

```text
開発モードを選択してください:
[1] 🚀 Speed Vibe Mode (Prototyping)  ... 速度最優先。ざっくり指示で即実装。
[2] 🛡️ Deep Dive Mode (Production Grade) ... 品質優先。詳細設計と承認プロセスを経て実装。
```

---

## 📂 生成される環境 (Directory Structure)

スクリプトを実行すると、以下のディレクトリ構成が自動的に展開されます。

```text
.
├── claw.md                      # 🦅 連携ルール定義書 (AIの行動指針)
├── setup_claw.js                # ⚙️ セットアップスクリプト (本ファイル)
├── antigravity.json             # 🔌 MCP設定 (Antigravity用)
├── claude.json                  # 🔌 MCP設定 (Claude Code用)
├── package.json                 # 📦 依存パッケージ定義
├── node_modules/                # 📦 インストールされたライブラリ
├── .claw/
│   └── templates/
│       ├── design_template.md        # 📄 システム詳細設計書テンプレート (全体設計)
│       └── program_spec_template.md  # 📄 プログラム仕様書テンプレート (詳細実装)
├── docs/
│   └── specs/                   # 📁 設計書・仕様書の格納先
└── tools/
    └── normalize_docs.js        # 🛠️ 文字コード正規化ツール (Shift-JIS対策)
```

---

## 📝 機能詳細

### 📐 仕様書テンプレート (Enterprise Templates)

`Mode 2 (Deep Dive)` を選択すると、`.claw/templates/` 配下に**外注開発にも対応できる品質のMarkdownテンプレート**が生成されます。
Antigravityは、必ずこのフォーマット（ER図、API定義、クラス図など）に従って設計書を作成します。

* **カスタマイズ**: 生成された `.md` ファイルを編集することで、独自のフォーマットルールをAIに適用させることができます。

### 🛠️ 自動正規化ツール (Auto-Normalizer)

セットアップ時に `tools/normalize_docs.js` が生成されます。
外部から持ち込まれた仕様書や古いソースコード（Shift-JISなど）を読み込む際に使用します。

* **自動実行ルール**: プロジェクト開始時やファイルインポート時、Antigravityはこのツールを実行するように義務付けられています。
* **機能**: エンコードを判別して **UTF-8** に統一し、元のファイルを `.bak` としてバックアップします。

---

## 🛠️ 2つの開発モード比較

| モード | 名称 | 用途 | 特徴 | ドキュメント |
| :--- | :--- | :--- | :--- | :--- |
| **Mode 1** | **🚀 Speed Vibe Mode** | - プロトタイプ<br>- ハッカソン<br>- PoC | **速度最優先**。<br>ざっくりした指示（Vibe）から即実装。 | 最小限。<br>走り書きレベルでOK。 |
| **Mode 2** | **🛡️ Deep Dive Mode** | - 本番開発<br>- チーム開発<br>- 複雑なアプリ | **品質・堅牢性優先**。<br>詳細設計を行い、承認後に実装。 | **厳格**。<br>指定テンプレートに準拠。 |

---

## 📈 開発ワークフロー

### Phase 0: Kickoff & Import

1. **Normalization**: 外部ファイルのインポートがある場合、まず `node tools/normalize_docs.js` で環境を浄化します。
2. **Definition**: Antigravityと対話して要件を固めます。

### Phase 1: Design (Mode 2のみ)

* Antigravityがテンプレートを用いて `docs/design.md` を作成します。
* **あなたの承認**があるまで、実装フェーズには進みません。

### Phase 2: Implementation (Claw Cycle)

1. **Components**: Antigravityが機能ごとに `docs/specs/[名称].md` (プログラム仕様書) を書きます。
2. **Coding (Backend)**: Claude Codeが仕様書通りに実装します。
3. **Frontend & Review**: Antigravityがフロントエンド実装とコードレビューを行います。

---
(c) 2026 KEIEI-NET
