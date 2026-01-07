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
                // Resolve anyway to prevent stopping the flow unless critical
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
// 📄 Templates
// ---------------------------------------------------------
const DESIGN_TEMPLATE = `# [Project Name] 詳細設計書 (System Design Document)

## 1. プロジェクト概要 (Overview)
... (Standard Content) ...

## 2. システムアーキテクチャ
...
`;

const PROG_SPEC_TEMPLATE = `# [Component Name] プログラム仕様書
...
`;

// ---------------------------------------------------------
// �️ Tool: Document Normalizer (Generate this script)
// ---------------------------------------------------------
const NORMALIZE_SCRIPT = `const fs = require('fs');
const path = require('path');
const iconv = require('iconv-lite');
const jschardet = require('jschardet');

const TARGET_EXTS = ['.md', '.txt', '.csv', '.json', '.js', '.ts'];
const IGNORE_DIRS = ['node_modules', '.git', 'dist', 'build'];

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
  const detected = jschardet.detect(buffer);
  
  if (!detected || !detected.encoding) return;
  
  const encoding = detected.encoding;
  // If already UTF-8 (and high confidence), skip
  if (encoding.toLowerCase() === 'utf-8' && detected.confidence > 0.9) return;

  console.log(\`Converting \${filepath} (From: \${encoding})...\`);
  
  // Backup
  fs.writeFileSync(filepath + '.bak', buffer);
  
  // Convert
  try {
    const str = iconv.decode(buffer, encoding);
    const utf8Buffer = iconv.encode(str, 'utf8');
    fs.writeFileSync(filepath, utf8Buffer);
    console.log(\`  ✅ Converted to UTF-8. Backup saved as .bak\`);
  } catch (e) {
    console.error(\`  ❌ Conversion failed: \${e.message}\`);
  }
};

console.log("🔍 Scanning for non-UTF-8 files...");
const allFiles = walkSync('.');
allFiles.forEach(f => convertFile(f));
console.log("✨ Normalization complete.");
`;

// ---------------------------------------------------------
// Mode Definitions
// ---------------------------------------------------------
const MODES = {
    '1': {
        name: '🚀 Speed Vibe Mode (Prototyping)',
        description: 'スピード優先モード / Build fast based on loose instructions.',
        workflow: `### Phase 0: Quick Start ⚡
1. **Input**: User gives a rough idea ("Vibe").
2. **Execution**: Antigravity generates scaffolding immediately.
3. **Iterate**: Claude Code implements tasks directly from chat.`
    },
    '2': {
        name: '🛡️ Deep Dive Mode (Production Grade)',
        description: '詳細設計モード / Enterprise Specs for System & Programs.',
        workflow: `### Phase 0: Detailed Architecture 🏛️
1. **Normalization**: Run \`node tools/normalize_docs.js\` to fix encoding of imported docs.
2. **System Spec**: Antigravity creates \`docs/design.md\` using \`design_template.md\`.
3. **Program Specs**: Create \`docs/specs/[Name].md\` using \`program_spec_template.md\`.
4. **Approval**: User MUST approve specs before coding starts.

### Phase 1: Structured Implementation
- **Frontend**: Antigravity implements strict component design.
- **Backend**: Claude Code implements API strictly following the Program Specs.`
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

## 1. コンテキスト同期
- Antigravity と Claude Code は同一のプロジェクトルートを共有します。
- 両ツールは MCP (Model Context Protocol) を使用して共有状態にアクセスします。

## 2. 役割と責任

### 🧠 Antigravity (設計 & フロントエンド)
- **設計 & ドキュメント**:
  - **テンプレートシステム**: \`.claw/templates/\` 内のテンプレートとツールを活用すること。
  - **自動正規化 (Auto-Normalize)**:
    - 外部ファイルを取り込んだり、プロジェクトを開始する際は、**一番最初に** \`node tools/normalize_docs.js\` を実行し、エンコードをUTF-8に統一すること。
    - これを怠ると文字化けのリスクがあるため、最優先事項とする。
  - **文書構成**:
    1. **システム詳細設計書**: \`docs/design.md\`
    2. **プログラム仕様書**: \`docs/specs/xxx.md\`
  - **逆同期 (Reverse Sync)**: 実装変更時はドキュメントを即時更新すること。
- **フロントエンド開発**:
  - UI設計および実装を担当。
- **監督**:
  - バックエンドコードの厳格なレビュー。

### ⚡ Claude Code (バックエンド専門)
- **バックエンド開発**:
  - プログラム仕様書に基づいて実装を行う。
  - 制約事項: **${modeChoice === '1' ? '速度優先' : '仕様書の完全再現'}**

## 3. ワークフロー

### Phase 0: Initialization & Import 📥
**Antigravity MUST execute the following sequence first:**
1.  **Normalization**: Run \`node tools/normalize_docs.js\` to fix encodings.
    - *If new files are added externally during development, Run this tool again.*
2.  **Kickoff**: Confirm requirements with User.

${mode.workflow.replace('### Phase 0: Detailed Architecture 🏛️\n1. **Normalization**: Run `node tools/normalize_docs.js` to fix encoding of imported docs.\n', '')}

## 4. ステータス
- **MCP Status**: Active
- **Template System**: Enabled
`;

    const files = {
        'claw.md': generateClawMd(selectedMode),
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

    console.log("\n📝 [Step 1/3] Generating configuration & templates...");
    try {
        const dirs = ['.claw/templates', 'docs/specs', 'tools'];
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
                if (filepath.includes('template') || filepath.includes('tools')) {
                    console.log(`  ℹ️  ${filepath} exists. Keeping user customization.`);
                }
            }
        }
    } catch (error) {
        console.error(`❌ Error: ${error.message}`);
        process.exit(1);
    }

    // 3. Dependencies (Sequential Execution)
    console.log("\n📦 [Step 2/3] Checking Dependencies...");
    try {
        if (!fs.existsSync('package.json')) {
            console.log("  Running npm init...");
            await runCommand('npm init -y');
        }

        console.log("  Installing packages (MCP SDK, Zod, iconv-lite, jschardet)...");
        // Added iconv-lite and jschardet for encoding support
        await runCommand('npm install @modelcontextprotocol/sdk zod iconv-lite jschardet --save');

        console.log("  ✅ Dependencies ready.");
        console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment Ready! ✨");
        console.log(`Current Mode: ${selectedMode.name}`);
        console.log(`Tools: Run 'node tools/normalize_docs.js' to fix file encodings.`);
        process.exit(0);
    } catch (e) {
        console.error("Setup Failed:", e);
        process.exit(1);
    }
})();
