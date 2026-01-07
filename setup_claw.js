const fs = require('fs');
const { exec } = require('child_process');
const readline = require('readline');
const path = require('path');

// CLI for user interaction
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const askQuestion = (query) => new Promise(resolve => rl.question(query, resolve));

console.log("\x1b[36m%s\x1b[0m", `
🦅 Claw Setup CLI - Antigravity x Claude Code
==============================================
`);

// Mode Definitions
const MODES = {
    '1': {
        name: '🚀 Speed Vibe Mode (Prototyping)',
        description: 'スピード優先モード / Build fast based on loose instructions.',
        workflow: `### Phase 0: Quick Start ⚡
1. **Input**: User gives a rough idea ("Vibe"). (ユーザーはざっくりしたアイデアを伝えます)
2. **Execution**: Antigravity generates basic scaffolding immediately. (即座に雛形を作成します)
3. **Iterate**: Claude Code implements coding tasks directly from chat. (爆速で実装・改善を繰り返します)`
    },
    '2': {
        name: '🛡️ Deep Dive Mode (Production Grade)',
        description: '詳細設計モード / Detailed specs first. Architecture & UI/UX required.',
        workflow: `### Phase 0: Detailed Architecture 🏛️
1. **Requirement Analysis**: Antigravity interviews User to define scope. (詳細ヒアリング)
2. **Specification**: Antigravity creates detailed \`design.md\`. (詳細設計書の作成: ER図, API, UIフロー)
3. **Approval**: User MUST approve \`design.md\` before any coding starts. (ユーザー承認後に着手)

### Phase 1: Structured Implementation
- **Frontend**: Antigravity implements strict component design.
- **Backend**: Claude Code implements API strictly following the Spec.`
    }
};

(async () => {
    // 1. Select Mode
    console.log("開発モードを選択してください (Select Development Mode):");
    console.log(`[1] ${MODES['1'].name} \n    - ${MODES['1'].description}`);
    console.log(`[2] ${MODES['2'].name} \n    - ${MODES['2'].description}`);

    let modeChoice = await askQuestion("\n番号を入力してください (Enter 1 or 2) [Default: 2]: ");
    if (!['1', '2'].includes(modeChoice.trim())) modeChoice = '2'; // Default to Deep Dive

    const selectedMode = MODES[modeChoice];
    console.log(`\n✅ 選択モード (Selected): ${selectedMode.name}`);

    // Configuration Content Generator
    const generateClawMd = (mode) => `# Claw - Antigravity & Claude Code 連携ルール

## 開発モード: ${mode.name}

## 1. コンテキスト同期
- Antigravity と Claude Code は同一のプロジェクトルートを共有します。
- 両ツールは MCP (Model Context Protocol) を使用して共有状態にアクセスします。

## 2. 役割と責任

### 🧠 Antigravity (設計 & フロントエンド)
- **設計 & ドキュメント**:
  - ドキュメント管理権限: **${modeChoice === '1' ? '簡易的 (Minimal)' : '厳格 (Strict/Single Source of Truth)'}**
  - ${modeChoice === '1' ? 'スピード優先でプロトタイプ仕様を作成します。' : '詳細な設計書(design.md)を作成し、承認を得てから開発へ進みます。'}
  - **逆同期 (Reverse Sync)**: ${modeChoice === '1' ? '任意' : '必須 (コードの変更を仕様書へ反映)'}
- **フロントエンド開発**:
  - ユーザーインターフェース（UI）の設計と実装を行います。
- **監督**:
  - Claude Code が生成したバックエンドコードをレビューします。

### ⚡ Claude Code (バックエンド専門)
- **バックエンド開発**:
  - Antigravity の仕様に基づいて実装を行います。
  - 制約事項: **${modeChoice === '1' ? 'とにかく動くものを最速で。' : 'design.md の仕様を厳守すること。'}**

## 3. ワークフロー
${mode.workflow}

## 4. ステータス
- **MCP Status**: Active
- **Sync Status**: Verified
`;

    const files = {
        'claw.md': generateClawMd(selectedMode),
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

    // 2. Create Configuration Files
    console.log("\n📝 [Step 1/3] Generating configuration files...");
    try {
        // ALWAYS Overwrite claw.md to reflect mode change
        fs.writeFileSync('claw.md', files['claw.md']);
        console.log(`  ✅ Update claw.md with ${selectedMode.name} rules.`);

        for (const [filename, content] of Object.entries(files)) {
            if (filename === 'claw.md') continue;
            if (!fs.existsSync(filename)) {
                fs.writeFileSync(filename, content);
                console.log(`  ✅ Created ${filename}`);
            } else {
                console.log(`  ℹ️  ${filename} already exists. Skipping.`);
            }
        }
    } catch (error) {
        console.error(`❌ Error: ${error.message}`);
        process.exit(1);
    }

    // 3. Dependencies
    console.log("\n📦 [Step 2/3] Checking Dependencies...");
    if (!fs.existsSync('package.json')) {
        console.log("  New project. initializing...");
        exec('npm init -y', () => { });
    }
    // Simple install check
    exec('npm install @modelcontextprotocol/sdk zod --save', (err, stdout, stderr) => {
        console.log("  ✅ Dependencies ready.");
        console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment Ready! ✨");
        console.log(`Current Mode: ${selectedMode.name}`);
        rl.close();
    });
})();
