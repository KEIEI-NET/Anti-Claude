# 🦅 Claw Template (Antigravity x Claude Code)

**"Brain & Brawn"** — 最強のAIペアプログラミング環境を、たった1コマンドで。

Claw Templateは、**Antigravity（設計・Frontend・監督）** と **Claude Code（Backend・実装）** をMCP (Model Context Protocol) で連携させ、爆速かつ高品質な開発を実現するためのテンプレートプロジェクトです。

---

## ✨ 特徴

* **⚡ 1-Click Setup**: 対話型CLIツールで、プロジェクトのセットアップが一瞬で完了します。
* **☯️ Dual Roles**: 「設計・思考」と「実装・速度」を分離。2つのAIが最適な役割を分担します。
* **🔄 Reverse Sync**: 実装コードから仕様書への逆反映をルール化。ドキュメントが陳腐化しません。
* **🎯 Selectable Modes**: プロトタイプ作成から本番開発まで、用途に合わせた2つのモードを搭載。

---

## 🚀 クイックスタート

### 1. プロジェクトの作成

このリポジトリをクローンするか、テンプレートとして使用して新しいディレクトリを作成します。

### 2. 環境セットアップ

以下のコマンドを実行します。

```bash
node setup_claw.js
```

### 3. モード選択

CLIが起動し、開発スタイルを聞かれます。

```text
開発モードを選択してください:
[1] 🚀 Speed Vibe Mode (Prototyping)
[2] 🛡️ Deep Dive Mode (Production Grade)
```

用途に合わせて番号（`1` または `2`）を入力してください。
（自動実行したい場合は `node setup_claw.js --mode=2` のように引数を指定可能です）

---

## 🛠️ 2つの開発モード

| モード | 名称 | 用途 | 特徴 | ドキュメント |
| :--- | :--- | :--- | :--- | :--- |
| **Mode 1** | **🚀 Speed Vibe Mode** | - プロトタイプ<br>- ハッカソン<br>- PoC | **速度最優先**。<br>ざっくりした指示（Vibe）から即実装。 | 最小限。<br>走り書きレベルでOK。 |
| **Mode 2** | **🛡️ Deep Dive Mode** | - 本番開発<br>- チーム開発<br>- 複雑なアプリ | **品質・堅牢性優先**。<br>詳細設計を行い、承認後に実装。 | **厳格**。<br>\`design.md\` が正解(SSOT)となる。 |

---

## 📈 開発ワークフロー

セットアップが完了したら、以下の流れで開発を進めます。
`claw.md` にルールが自動生成されているため、Antigravity (AIエージェント) はこのルールに従って動きます。

### Phase 0: Kickoff (キックオフ)

* **あなた**: 「こういうアプリを作りたい」とAntigravityに伝えます。
* **Antigravity**: モードに応じて、ヒアリング手法を変えます。
  * (Mode 1) 「了解、すぐ作ります」
  * (Mode 2) 「要件定義を詰めましょう。技術スタックはどうしますか？」

### Phase 1: Design & Architecture (設計)

* **Antigravity**: 必要に応じて `design.md` (設計書) を作成します。
  * Mode 2では、**あなたの承認**があるまで実装に進みません。
* UI/UXデザインもこのフェーズで決定します。

### Phase 2: Implementation (実装サイクル)

ここから自動化プロセスが加速します。

1. **Front & Oversight (Antigravity)**
    * フロントエンドの実装。
    * プロジェクト全体の構成管理。
2. **Backend & Logic (Claude Code)**
    * Antigravityの指示に基づき、APIやDB処理を高速実装。
3. **Reverse Sync (逆同期)**
    * 実装中に判明した仕様変更を、Antigravityが設計書へ書き戻します。

---

## 📂 フォルダ構成

* **`claw.md`**: 連携ルール定義ファイル（セットアップ時に自動生成/更新）。
* **`setup_claw.js`**: セットアップ用CLIスクリプト。
* **`antigravity.json`**: Antigravity用のMCP設定。
* **`claude.json`**: Claude Code用のMCP設定。
* **`.agent/workflows/`**: AIエージェント用の自動化ワークフロー定義。

---

## ❓ FAQ

**Q. 途中でモードを変更できますか？**
A. はい。再度 `node setup_claw.js` を実行し、モードを選び直してください。`claw.md` のルールが書き換わります（既存のコードは消えません）。

**Q. 英語のドキュメントにしたいのですが。**
A. 現状のテンプレートは日本語ユーザー向けに最適化されていますが、AIに対して「英語で進めて」と指示すれば、英語でのドキュメント作成も可能です。

---
(c) 2026 KEIEI-NET
