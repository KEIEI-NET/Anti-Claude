# 🦅 Claw Template (Antigravity x Claude Code)

> **最強のAIペアプログラミング環境を、たった1ファイルのスクリプトから。**

Claw Templateは、**Antigravity (Architect)** と **Claude Code (Backend)** を連携させ、堅牢な **Clean Architecture / DDD** システムを爆速で構築するための環境構築キットです。
v5.3では、既存の仕様書を一切劣化させずに取り込み、最新アーキテクチャへと昇華させる「**Deep Import**」機能を搭載しました。

## 🚀 クイックスタート (3 Steps)

このリポジトリをcloneする必要はありません（インストーラー以外のファイルはむしろ邪魔です）。

### Step 1: 準備

空のプロジェクトフォルダを作成し、[setup_claw.js](setup_claw.js) だけをコピーして配置してください。

### Step 2: インストール

ターミナルで以下を実行します。必要なファイルと依存ライブラリが全自動で生成されます。

```bash
node setup_claw.js
```

実行後、**開発モード**を選択します：

- `[1]` 🚀 **Speed Vibe Mode**: プロトタイプ向け。とにかく速く動くものを作る。
- `[2]` 🛡️ **Deep Dive Mode** (推奨): 本格開発向け。Clean Architecture & DDDを強制。

### Step 3: プロジェクト開始 (KICKOFF)

環境構築が終わると、**`KICKOFF.md`** というファイルが生成されます。
このファイルを開き、**今のあなたの状況に合ったプロンプト** をコピーして、Antigravity (AI) に送信してください。

| あなたの状況 | 推奨アクション |
| :--- | :--- |
| **既存の仕様書がある** | 資料を `input_docs/` に入れて、インポート用プロンプトを送信 |
| **アイデアだけある** | 「新規設計」用プロンプトを送信し、要件定義から開始 |
| **実装に入りたい** | 設計書 (`docs/design.md`) 完成後、「プログラム仕様書作成」用プロンプトを送信 |

---

## ✨ v5.3 の新機能: "Deep Import"

自動生成される `input_docs/` フォルダに、既存の仕様書（Markdown, Text, Source Code等）を放り込んでください。
Antigravityは、単にそれを読むだけでなく、以下のポリシーに従って **アップグレード（Upgrade）** します。

### ✅ No Downgrade Policy (劣化厳禁)

AI特有の「勝手な要約」を許しません。

- **図解の完全包含**: 元の物理構成図（ASCIIアート等）をそのまま残しつつ、新たにClean Architectureの論理図を追加します。
- **非機能要件の全移植**: セキュリティ要件、エンコーディング仕様、ログ形式などを一字一句漏らさず移植します。
- **コードブロック維持**: 構造体定義やSQLクエリ等の技術詳細をそのまま保持します。

### ✅ Architecture Upgrade (品質向上)

古い設計（MVC等）で書かれた仕様書であっても、**Clean Architecture / DDD** の観点で再構成し、堅牢な設計書 (`docs/design.md`) へと生まれ変わらせます。

---

## 📂 生成される環境

セットアップ完了後のディレクトリ構成です。

```text
.
├── setup_claw.js                # ⚙️ インストーラー (これ以外は削除してOK)
├── KICKOFF.md                   # 🚀 AIへの指示書 (プロジェクトのナビゲーター)
├── claw.md                      # 🦅 AI行動指針 (Clean Arch/DDDのルール定義)
├── antigravity.json             # 🔌 MCP設定ファイル
├── claude.json                  # 🔌 MCP設定ファイル
├── input_docs/                  # 📥 ここに資料を入れると全自動インポート
├── docs/                        # 📘 設計書・仕様書の出力先
└── .claw/templates/
    ├── design_template.md       # 📄 システム詳細設計書 (物理図・論理図併記対応)
    └── program_spec_template.md # 📄 プログラム仕様書 (詳細実装レベル)
```

---

## 🛠️ モード比較表

| 特徴 | 🚀 Mode 1: Speed Vibe | 🛡️ Mode 2: Deep Dive (推奨) |
| :--- | :--- | :--- |
| **目的** | ハッカソン、PoC、使い捨て | 長期運用、業務システム、受託開発 |
| **アーキテクチャ** | 自由 (MVCなど) | **Clean Architecture + DDD** |
| **設計書の粒度** | メモレベル | **詳細仕様書 (7章以上)** |
| **AIの役割** | とにかくコードを書く | 設計者(Anti) と 実装者(Claude) の分業 |
| **インポート挙動** | 参考程度に読む | **品質を高めて完全移植する** |
