const fs = require('fs');
const { exec } = require('child_process');
const path = require('path');

console.log("\x1b[36m%s\x1b[0m", "🚀 Starting Claw (Antigravity x Claude Code) Setup [DEBUG MODE]...");

// Configuration Files Content
const files = {
    'claw.md': `# Claw - Antigravity & Claude Code 連携ルール

## 1. コンテキスト同期
- Antigravity と Claude Code は同一のプロジェクトルートを共有します。
- 両ツールは MCP (Model Context Protocol) を使用して共有状態にアクセスします。

## 2. 役割と責任

### 🧠 Antigravity (設計 & フロントエンド)
- **設計 & ドキュメント**:
  - すべての設計ドキュメント（仕様書、アーキテクチャ、UI/UX）を作成・維持・更新します。
  - プロジェクト要件の「唯一の正解（Single Source of Truth）」として振る舞います。
  - **逆同期 (Reverse Sync)**: Claude Codeの実装内容を定期的に監査し、コード側の変更・改善を仕様書に反映します。
- **フロントエンド開発**:
  - ユーザーインターフェース（UI）の設計と実装を行います。
  - クライアントサイドのロジック、レスポンシブ対応、美観を担当します。
- **監督**:
  - Claude Code が生成したバックエンドコードをレビューし、設計書と一致しているか確認します。

### ⚡ Claude Code (バックエンド専門)
- **バックエンド開発**:
  - Antigravity の仕様に基づいて、サーバーサイドロジック、API、データベーススキーマを実装します。
  - アルゴリズムとデータ処理を最適化します。
- **実行**:
  - バックエンドのボイラープレート構築やレガシーコードのリファクタリングを高速に実行します。

## 3. ワークフロー

### Phase 0: プロジェクト・キックオフ 🚀
**Antigravity はプロジェクト開始時に必ず以下を確認すること:**
1.  **仕様策定のアプローチ**:
    - 🗣️ **壁打ち (Interactive)**: 要件定義から会話形式で一緒に作り上げる。
    - 📄 **既存仕様書あり**: ユーザーが提示するMDファイルの仕様書に基づき開発する。
2.  **技術スタックの選定**:
    - ❓ **未定**: 要件に基づいて Antigravity が提案・選択する。
    - 🎯 **決定済み**: ユーザーの指示に従う。

### Phase 1: アーキテクチャ & 設計
- Antigravity が \`design.md\` を作成/更新します。
- Antigravity が UI/UX システムを設計します。

### Phase 2: 実装 (Clawサイクル)
1. **Antigravity**: フロントエンド構築 & プロジェクト構成セットアップ。
2. **Claude Code**: バックエンドAPI & コアロジック実装。
3. **Antigravity**: フロントエンド統合 & UI仕上げ。

### Phase 3: レビュー
- 実装が \`design.md\` に沿っているか検証します。

## 4. ステータス
- **MCP Status**: Active
- **Sync Status**: Verified
`,
    'claude.json': JSON.stringify({
        "mcpServers": {
            "filesystem": {
                "command": "npx",
                "args": ["-y", "@modelcontextprotocol/server-filesystem", "."]
            }
        }
    }, null, 2),
    'antigravity.json': JSON.stringify({
        "mcpServers": {
            "filesystem": {
                "command": "npx",
                "args": ["-y", "@modelcontextprotocol/server-filesystem", "."]
            }
        },
        "contextSharing": true,
        "partnerTool": "Claude Code"
    }, null, 2)
};

// 1. Create Configuration Files
console.log("📝 [Step 1/3] Checking configuration files...");
try {
    for (const [filename, content] of Object.entries(files)) {
        if (!fs.existsSync(filename)) {
            fs.writeFileSync(filename, content);
            console.log(`  ✅ Created ${filename}`);
        } else {
            console.log(`  ℹ️  ${filename} already exists. Skipping (Preserving existing config).`);
        }
    }
} catch (error) {
    console.error(`❌ [FATAL] Error creating files: ${error.message}`);
    process.exit(1);
}

// 2. Initialize npm and Install Dependencies
console.log("\n📦 [Step 2/3] Managing Dependencies...");

const runCommand = (command) => {
    console.log(`  > Executing: ${command}`);
    return new Promise((resolve, reject) => {
        exec(command, (error, stdout, stderr) => {
            if (error) {
                console.error(`  ❌ Command Failed: ${command}`);
                console.error(`  Error details: ${error.message}`);

                // Keep minimal output unless verbose needed, but for debug request, show all.
                if (stdout) console.log(`  [stdout]: ${stdout.trim()}`);
                if (stderr) console.error(`  [stderr]: ${stderr.trim()}`);
            } else {
                if (stdout && stdout.trim()) console.log(`    ${stdout.trim().split('\n').join('\n    ')}`);
                // npm install info often comes in stderr
                if (stderr && stderr.trim()) console.log(`    (info) ${stderr.trim().split('\n').join('\n    ')}`);
                console.log(`  ✅ Command completed.`);
            }
            resolve();
        });
    });
};

(async () => {
    // Check for package.json
    if (!fs.existsSync('package.json')) {
        console.log("  New project detected. Running npm init...");
        await runCommand('npm init -y');
    } else {
        console.log("  package.json found. Skipping npm init.");
    }

    console.log("  Installing/Updating packages...");
    await runCommand('npm install @modelcontextprotocol/sdk zod --save');

    // 3. Verification
    console.log("\n🔍 [Step 3/3] Verifying Installation...");
    const nodeModulesPath = path.join(process.cwd(), 'node_modules', '@modelcontextprotocol');
    if (fs.existsSync(nodeModulesPath)) {
        console.log("  ✅ MCP SDK found in node_modules.");
    } else {
        console.warn("  ⚠️  Warning: MCP SDK not found in node_modules after install try.");
    }

    console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment Setup Finished! ✨");
    console.log("Debugging complete. System is ready.");
})();
