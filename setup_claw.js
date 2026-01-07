const fs = require('fs');
const { exec } = require('child_process');
const readline = require('readline');
const path = require('path');

// ---------------------------------------------------------
// 🔧 Utility: Promisified Exec
// ---------------------------------------------------------
const runCommand = (command) => {
    return new Promise((resolve, reject) => {
        exec(command, (error, stdout, stderr) => {
            if (error) {
                console.warn(`⚠️  Warning in command: ${command}\n${stderr}`);
                resolve(stdout || stderr);
            } else {
                resolve(stdout);
            }
        });
    });
};

// CLI for user interaction
const askQuestion = (query) => {
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    return new Promise(resolve => rl.question(query, answer => {
        rl.close();
        resolve(answer);
    }));
};

console.log("\x1b[36m%s\x1b[0m", `
🦅 Claw Setup CLI - Antigravity x Claude Code
==============================================
`);

// ---------------------------------------------------------
// 📄 Kickoff Prompts (New Feature)
// ---------------------------------------------------------
const KICKOFF_CONTENT = `# 🚀 Claw Kickoff Prompts

環境セットアップ後、以下の手順でプロジェクトを開始してください。

## 📥 既存の仕様書がある場合 (Import Flow)
**手順**:
1. プロジェクトルートにある \`input_docs/\` フォルダに、既存の資料（Markdown, Text, Source Code等）を全て入れてください。
2. 以下のコマンドをチャットに貼り付けてください。

\`\`\`text
@Antigravity
【プロジェクト開始: 既存仕様の全量インポート】
1. まず \`node tools/normalize_docs.js\` を実行し、\`input_docs/\` 内の全ファイルの文字コードをUTF-8に統一してください。
2. 次に、\`input_docs/\` フォルダ内にある **全てのファイルの内容を読み込んで** ください。
3. 読み込んだ内容を「前提知識」として学習し、不明点があれば質問してください。なければ \`docs/design.md\` の構成案を作成してください。
\`\`\`

## 🆕 新規開発の場合 (New Design Flow)
**手順**: 以下のコマンドをチャットに貼り付けてください。

\`\`\`text
@Antigravity
【プロジェクト開始: 新規設計】
1. 詳細設計モード(Deep Dive)で進めます。
2. まずは私の作りたいアプリの「要件ヒアリング」を開始してください。
3. ヒアリング後、.claw/templates/design_template.md に基づいて docs/design.md を作成してください。
\`\`\`

## ⚙️ プログラム仕様書の作成 (Implementation Prep)
**手順**: 設計完了後、実装に入る前に実行します。

\`\`\`text
@Antigravity
【フェーズ移行: プログラム詳細設計】
1. docs/design.md の内容に基づき、優先度の高いコンポーネントから順にプログラム仕様書を作成してください。
2. テンプレートは .claw/templates/program_spec_template.md を厳守すること。
3. まずは [コンポーネント名] の仕様書作成をお願いします。
\`\`\`

## ⚡ とにかく動くものを作りたい (Speed Vibe Mode)
**手順**: 

\`\`\`text
@Antigravity
【プロジェクト開始: スピード優先】
1. アプリの概要は「[ここにアイデアを入力]」です。
2. 面倒な設計書はスキップして、すぐに動くプロトタイプの実装を開始してください。
3. 必要なファイル構成を提案し、Claude Codeに実装指示を出してください。
\`\`\`
`;

// ---------------------------------------------------------
// 📄 System Design Template
// ---------------------------------------------------------
const DESIGN_TEMPLATE = `# [Project Name] システム詳細設計書 (System Design Document)

## 1. プロジェクト概要 (Overview)
- **目的**: 
- **対象ユーザー**: 
- **主要機能**: 

## 2. ドメイン設計 (Domain Design / DDD)
### 2.1 境界づけられたコンテキスト (Bounded Contexts)
- **Context A**: ...
- **Context B**: ...

### 2.2 ドメインモデル (Domain Models)
\`\`\`mermaid
classDiagram
    class User {
        +UserId id
        +UserName name
        +email changeEmail()
    }
\`\`\`

## 3. システムアーキテクチャ (Clean Architecture)
### 3.1 レイヤー構成
- **Presentation Layer**: UI, API Handlers
- **Application Layer**: UseCases
- **Domain Layer**: Entities, ValueObjects, Domain Services (No dependencies)
- **Infrastructure Layer**: DB, External APIs

### 3.2 システム構成図 (Mermaid)
\`\`\`mermaid
graph TD
    Client --> Presenter
    Presenter --> UseCase
    UseCase --> Domain
    UseCase --> RepositoryInterface
    Infrastructure --> RepositoryInterface
\`\`\`

## 4. API インターフェース仕様
...

## 5. UI/UX 設計
...

## 6. 非機能要件
...
`;

// ---------------------------------------------------------
// 📄 Program Specification Template
// ---------------------------------------------------------
const PROG_SPEC_TEMPLATE = `# [Component Name] プログラム仕様書詳細版

## 目次

1. [概要 (Overview)](#1-概要-overview)
2. [アーキテクチャ設計 (Architecture Design)](#2-アーキテクチャ設計-architecture-design)
3. [環境・依存関係 (Environment & Dependencies)](#3-環境依存関係-environment--dependencies)
4. [インターフェース定義 (Interface Definition)](#4-インターフェース定義-interface-definition)
5. [データモデル (Data Models)](#5-データモデル-data-models)
6. [機能詳細 (Functional Details)](#6-機能詳細-functional-details)
7. [非機能要件 (Non-Functional Requirements)](#7-非機能要件-non-functional-requirements)
8. [セキュリティ設計 (Security Design)](#8-セキュリティ設計-security-design)
9. [エラーハンドリング (Error Handling)](#9-エラーハンドリング-error-handling)
10. [テスト・品質保証 (Test & QA)](#10-テスト品質保証-test--qa)
11. [運用・監視 (Operations & Monitoring)](#11-運用監視-operations--monitoring)
12. [付録 (Appendix)](#12-付録-appendix)

---

## 1. 概要 (Overview)

### 1.1 目的 (Purpose)

[このコンポーネントが達成すべき目的を具体的かつ定量的に記述する]

### 1.2 スコープ (Scope)

- **対象 (In-Scope)**: [実装する機能、サポートする環境]
- **対象外 (Out-of-Scope)**: [今回は実装しない機能、前提としない環境]

## 2. アーキテクチャ設計 (Architecture Design)

### 2.1 システム構成図 (System Architecture)

\`\`\`mermaid
graph TD;
    A-->B;
\`\`\`

[Mermaid記法やASCIIアートで構成図を記述]

### 2.2 モジュール構成 (Module Structure)

[ディレクトリ構成やパッケージ構成の定義]

## 3. 環境・依存関係 (Environment & Dependencies)

### 3.1 開発言語・フレームワーク

| 項目 | バージョン/要件 | 備考 |
|------|-----------------|------|
| 言語 | | |
| フレームワーク | | |

### 3.2 外部ライブラリ (Libraries)

[主要な依存ライブラリ一覧]

## 4. インターフェース定義 (Interface Definition)

### 4.1 APIエンドポイント (API Endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/resource | ... |

### 4.2 入出力仕様 (I/O Specs)

#### 入力 (Input)

- 環境変数
- 引数・パラメータ

#### 出力 (Output)

- 戻り値
- ログ出力形式

## 5. データモデル (Data Models)

### 5.1 データベース設計 (Schema)

[テーブル定義、ER図]

### 5.2 構造体・クラス定義 (Class Definitions)

\`\`\`go
// 主要なデータ構造の定義
type Example struct {
    ID string \`json:"id"\`
}
\`\`\`

## 6. 機能詳細 (Functional Details)

### 6.1 [機能名A]

#### 概要

[機能の説明]

#### ロジック・アルゴリズム

[処理フローの詳細]

### 6.2 [機能名B]

...

## 7. 非機能要件 (Non-Functional Requirements)

### 7.1 パフォーマンス (Performance)

[レスポンスタイム目標、スループット等]

### 7.2 可用性・拡張性 (Availability & Scalability)

[冗長化方針、スケールアウト計画]

## 8. セキュリティ設計 (Security Design)

### 8.1 認証・認可 (AuthN/AuthZ)

[認証方式の詳細]

### 8.2 データ保護 (Data Protection)

[暗号化、マスキング処理]

## 9. エラーハンドリング (Error Handling)

### 9.1 エラーコード体系

| Code | Type | Description |
|------|------|-------------|
| E001 | Auth | ... |

### 9.2 リカバリープラン

[障害時の復旧手順、Circuit Breaker等]

## 10. テスト・品質保証 (Test & QA)

### 10.1 テスト戦略

[単体テスト、結合テスト、E2Eテストの範囲]

### 10.2 CI/CDパイプライン

[自動テスト、Lint、ビルドフロー]

## 11. 運用・監視 (Operations & Monitoring)

### 11.1 ログ設計

[ログレベル、出力フォーマット]

### 11.2 メトリクス・アラート

[監視項目、閾値]

## 12. 付録 (Appendix)

### 12.1 関連ドキュメント

### 12.2 用語集
`;

// ---------------------------------------------------------
// 🛠️ Tool: Document Normalizer
// ---------------------------------------------------------
const NORMALIZE_SCRIPT = `const fs = require('fs');
const path = require('path');
const iconv = require('iconv-lite');
const jschardet = require('jschardet');

const TARGET_EXTS = ['.md', '.txt', '.csv', '.json', '.js', '.ts', '.go', '.py', '.java'];
const IGNORE_DIRS = ['node_modules', '.git', 'dist', 'build', 'obj', 'bin'];

const walkSync = (dir, filelist = []) => {
  const files = fs.readdirSync(dir);
  files.forEach(file => {
    if (IGNORE_DIRS.includes(file)) return;
    const filepath = path.join(dir, file);
    if (fs.statSync(filepath).isDirectory()) {
      filelist = walkSync(filepath, filelist);
    } else {
      if (TARGET_EXTS.includes(path.extname(file))) {
        filelist.push(filepath);
      }
    }
  });
  return filelist;
};

const convertFile = (filepath) => {
  const buffer = fs.readFileSync(filepath);
  // jschardet may fail on very small files, default to utf8 if null
  const detected = jschardet.detect(buffer);
  
  if (!detected || !detected.encoding) return;
  
  const encoding = detected.encoding;
  if (encoding.toLowerCase() === 'utf-8' && detected.confidence > 0.9) return;

  console.log(\`Converting \${filepath} (From: \${encoding})...\`);
  fs.writeFileSync(filepath + '.bak', buffer);
  
  try {
    const str = iconv.decode(buffer, encoding);
    const utf8Buffer = iconv.encode(str, 'utf8');
    fs.writeFileSync(filepath, utf8Buffer);
    console.log(\`  ✅ Converted to UTF-8. Backup saved.\`);
  } catch (e) {
    console.error(\`  ❌ Conversion failed: \${e.message}\`);
  }
};

console.log("🔍 Scanning for non-UTF-8 files...");
try {
    const allFiles = walkSync('.');
    allFiles.forEach(f => convertFile(f));
    console.log("✨ Normalization complete.");
} catch(e) {
    console.error("Error during normalization:", e);
}
`;

// ---------------------------------------------------------
// Mode Definitions
// ---------------------------------------------------------
const MODES = {
    '1': {
        name: '🚀 Speed Vibe Mode (Prototyping)',
        description: 'スピード優先モード / Build fast based on loose instructions.',
        workflow: `### Phase 0: Quick Start ⚡
1. **Kickoff**: User sends "Start" command (See KICKOFF.md).
2. **Execution**: Antigravity generates scaffolding immediately.
3. **Iterate**: Claude Code implements tasks directly from chat.`
    },
    '2': {
        name: '🛡️ Deep Dive Mode (Clean Arch & DDD)',
        description: '詳細設計モード / Clean Architecture, DDD, SOLID Principles.',
        workflow: `### Phase 0: Domain Analysis & Design 🏛️
1. **Kickoff**: User sends "Import" or "New Design" command (See KICKOFF.md).
2. **Normalization**: If importing, Antigravity runs \`node tools/normalize_docs.js\`.
3. **Domain Modeling**: Antigravity analyzes requirements using **DDD**.
4. **Specification**: 
   - Create \`docs/design.md\`.
   - Create \`docs/specs/[Component].md\`.
5. **Approval**: User MUST approve models & specs.

### Phase 1: Implementation (SOLID Principles)
- **Domain Layer**: Implement Pure Domain Logic.
- **Application Layer**: Implement Use Cases.
- **Interface/Infra**: Adapters & DB.`
    }
};

(async () => {
    // 1. Select Mode
    let modeChoice;
    const args = process.argv.slice(2);
    const modeArg = args.find(arg => arg.startsWith('--mode='));

    if (modeArg) {
        modeChoice = modeArg.split('=')[1];
        console.log(`🤖 Auto-detected mode from arguments: ${modeChoice}`);
    } else {
        console.log("開発モードを選択してください (Select Development Mode):");
        console.log(`[1] ${MODES['1'].name} \n    - ${MODES['1'].description}`);
        console.log(`[2] ${MODES['2'].name} \n    - ${MODES['2'].description}`);
        modeChoice = await askQuestion("\n番号を入力してください (Enter 1 or 2) [Default: 2]: ");
    }

    if (!['1', '2'].includes(modeChoice.trim())) modeChoice = '2';

    const selectedMode = MODES[modeChoice];
    console.log(`\n✅ 選択モード (Selected): ${selectedMode.name}`);

    // 2. Generate Files Content
    const generateClawMd = (mode) => `# Claw - Antigravity & Claude Code 連携ルール

## 開発モード: ${mode.name}

## 1. アーキテクチャ原則
- **Clean Architecture** & **DDD** & **SOLID原則** を遵守。

## 2. 役割と責任
- **Antigravity**: Architect, Domain Expert, Frontend.
- **Claude Code**: Backend Implementation (SOLID compliant).

## 3. まずはじめに (Getting Started)
**KICKOFF.md を参照し、適切なコマンドをAntigravityに送信してください。**

## 4. ワークフロー
${mode.workflow}

## 5. ステータス
- **MCP Status**: Active
- **Template System**: Enabled (Clean Arch/DDD)
`;

    const files = {
        'claw.md': generateClawMd(selectedMode),
        'KICKOFF.md': KICKOFF_CONTENT,  // New File!
        'claude.json': JSON.stringify({
            "mcpServers": {
                "filesystem": { "command": "npx", "args": ["-y", "@modelcontextprotocol/server-filesystem", "."] }
            }
        }, null, 2),
        'antigravity.json': JSON.stringify({
            "mcpServers": {
                "filesystem": { "command": "npx", "args": ["-y", "@modelcontextprotocol/server-filesystem", "."] }
            },
            "contextSharing": true,
            "partnerTool": "Claude Code"
        }, null, 2),
        '.claw/templates/design_template.md': DESIGN_TEMPLATE,
        '.claw/templates/program_spec_template.md': PROG_SPEC_TEMPLATE,
        'tools/normalize_docs.js': NORMALIZE_SCRIPT
    };

    console.log("\n📝 [Step 1/4] Generating configuration & templates...");
    try {
        const dirs = ['.claw/templates', 'docs/specs', 'tools', 'input_docs'];
        dirs.forEach(d => {
            const p = path.join(process.cwd(), d);
            if (!fs.existsSync(p)) fs.mkdirSync(p, { recursive: true });
        });

        fs.writeFileSync('claw.md', files['claw.md']);
        console.log(`  ✅ Update claw.md with ${selectedMode.name} rules.`);

        for (const [filepath, content] of Object.entries(files)) {
            if (filepath === 'claw.md') continue;
            const fullPath = path.join(process.cwd(), filepath);
            if (!fs.existsSync(fullPath)) {
                fs.writeFileSync(fullPath, content);
                console.log(`  ✅ Created ${filepath}`);
            } else {
                if (filepath.includes('template') || filepath.includes('tools') || filepath.includes('KICKOFF')) {
                    console.log(`  ℹ️  ${filepath} exists. Keeping user customization.`);
                }
            }
        }
    } catch (error) {
        console.error(`❌ Error: ${error.message}`);
        process.exit(1);
    }

    // 3. Dependencies
    console.log("\n📦 [Step 2/4] Checking Dependencies...");
    try {
        if (!fs.existsSync('package.json')) {
            console.log("  Running npm init...");
            await runCommand('npm init -y');
        }

        console.log("  Installing packages (MCP SDK, Zod, iconv-lite, jschardet)...");
        await runCommand('npm install @modelcontextprotocol/sdk zod iconv-lite jschardet --save');

        console.log("  ✅ Dependencies ready.");
        console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment Ready (v4.1)! ✨");
        console.log(`Current Mode: ${selectedMode.name}`);
        console.log(`🚀 Next Step: Open KICKOFF.md and copy the start command to the chat.`);
        process.exit(0);
    } catch (e) {
        console.error("Setup Failed:", e);
        process.exit(1);
    }
})();
