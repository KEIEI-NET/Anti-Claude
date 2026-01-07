# 🚀 Claw Kickoff Prompts (v7.0)

## 🏛️ Phase 1: Deep Design & Orchestration
**Antigravity**が全体の設計と、**Claude Code**への並列タスク指示を行います。

```text
@Antigravity
【プロジェクト開始: オーケストレーション】
1. input_docs/ から全ての仕様を読み込み docs/design.md を生成してください。
2. その後、コンポーネントごとに TASKS.md を分割し、Claude Code (CLI) を使って並行実装を開始してください。
3. 長時間のタスク（20分以上）も制限せず、生存監視を行ってください。
```

## 🧠 Phase 2: Claude Skills Invocation (Deep Work)
特定のスキルを呼び出す際のコマンドセットです。

### 🛡️ セキュリティ監査 (Security Audit)
```text
@Antigravity
Claude Code に Security Review Skills を適用させ、現状のコードの脆弱性を洗い出し、修正させてください。
CMD: claude -p "Perform a deep security audit on internal/ and implement necessary fixes."
```

### 💎 品質リファクタリング (Quality Refactor)
```text
@Antigravity
Claude Code に Refactoring Skills を適用させ、SOLID原則に基づいたコードの洗練を行わせてください。
CMD: claude -p "Refactor the code in internal/domain/ for perfect SOLID compliance and clean architecture."
```

### 🧪 テスト駆動開発 (TDD & Quality Check)
```text
@Antigravity
Claude Code に Testing Skills を適用させ、カバレッジ100%を目指してテストを拡充させてください。
CMD: claude -p "Increase test coverage to 100% for all components in internal/ and fix any discovered bugs."
```

## 🕰️ 長時間タスクのルール
- Claude Code が思考中の場合、Antigravity は途中で強制終了してはいけません。
- ファイルの更新を 30秒ごとに監視し、動きがある限り継続させます。
- 終了後、Antigravity が精査し統合します。
