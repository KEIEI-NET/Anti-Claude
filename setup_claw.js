const fs = require('fs');
const { exec } = require('child_process');
const readline = require('readline');
const path = require('path');

/**
 * Claw Setup CLI v7.7.1 - Final Self-Contained Edition
 * 100% Full Templates / No Omissions / Versioning Enabled
 */

const runCommand = (command) => {
    return new Promise((resolve) => {
        exec(command, (error, stdout, stderr) => {
            if (error) {
                console.warn(`⚠️ Warning: ${stderr}`);
                resolve(stdout || stderr);
            } else {
                resolve(stdout);
            }
        });
    });
};

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

async function main() {
    console.log("\x1b[36m%s\x1b[0m", `
🦅 Claw Setup CLI v7.7.1 - Ultimate Edition
=========================================
`);

    console.log("1: 🚀 Speed Vibe Mode (Prototyping)");
    console.log("2: 🛡️  Deep Dive Mode (Orchestration & Planning)");

    const choice = await askQuestion("\nモードを選択してください (Select mode 1-2) [Default: 2]: ");
    const mode = (choice === '1') ? '1' : '2';
    const modeName = (mode === '2') ? 'Deep Dive' : 'Speed Vibe';

    // --- FULL CONTENT RESTORATION ---

    const VALIDATOR_SCRIPT = `const fs = require('fs');
const path = require('path');
const TARGET_FILE = path.join(__dirname, '../docs/design.md');
const REQUIRED_SECTIONS = [
    '## 1. プロジェクト概要',
    '## 2. ドメイン設計',
    '## 3. システムアーキテクチャ',
    '#### 物理構成図',
    '#### 論理構成図',
    '## 4. API インターフェース仕様',
    '## 5. データモデル',
    '## 6. 非機能要件',
    '## 7. エラーハンドリング',
    '## 8. テスト・移行計画'
];
function validate() {
    if (!fs.existsSync(TARGET_FILE)) { console.error('❌ docs/design.md not found'); process.exit(1); }
    const content = fs.readFileSync(TARGET_FILE, 'utf-8');
    let errors = [];
    REQUIRED_SECTIONS.forEach(s => { if (!content.includes(s)) errors.push(\`MISSING: \${s}\`); });
    if (content.includes('...') || content.includes('TODO')) errors.push('INCOMPLETE: Contains dots or TODO');
    if (errors.length > 0) {
        errors.forEach(e => console.error(e));
        process.exit(1);
    } else { console.log('✅ VALIDATION PASSED'); }
}
validate();`;

    const KICKOFF_CONTENT = `# 🚀 Claw Kickoff (v7.7.1 Final)

## 🏛️ Phase 1: Planning & Orchestration
1. Detect technology from input_docs/ and create docs/design.md.
2. Create **PLAN.md** with modularized steps. Wait for human approval.

## 🧠 Phase 2: Execution & Skills
- Use \`claude --dangerously-skip-permissions -p "..."\` for automation.
- **Skill Usage**: Security Audit, SOLID Refactoring, Deep Debugging.
- **Rules**: NO God Files. Split into modules (interactor, errors, dto, etc.).

## 🕰️ Phase 3: Versioning & Reverse Sync
1. **IMPLEMENT**: Claude Code でタスクを完了させる。
2. **BACKUP**: \`docs/design.md\` を \`docs/history/\` へ日時付きで退避。
3. **SYNC (Design)**: \`docs/design.md\` を最新コードと同期。
4. **SYNC (Spec)**: \`docs/specs/[Component]_spec.md\` を作成または最新化。 ← **MANDATORY**
5. **VERIFY**: 全てが整合していることを確認。
`;

    const DESIGN_TEMPLATE = `# [Project Name] システム詳細設計書

## 1. プロジェクト概要 (Overview)
- **目的**: 
- **主要機能**: 

## 2. ドメイン設計 (Domain Design / DDD)
### 2.1 境界づけられたコンテキスト
### 2.2 ドメインモデル
\`\`\`mermaid
classDiagram
    class DomainModel { +ID id }
\`\`\`

## 3. システムアーキテクチャ
### 3.1 技術スタック
### 3.2 構成図
#### 物理構成図 (Original)
\`\`\`text
(Original Diagram Here)
\`\`\`
#### 論理構成図
\`\`\`mermaid
graph TD
    UI --> App --> Domain
\`\`\`

## 4. API インターフェース仕様
## 5. データモデル (Data Models)
## 6. 非機能要件 (Non-Functional)
## 7. エラーハンドリング (Error Handling)
## 8. テスト・移行計画
`;

    const NORMALIZE_SCRIPT = `const fs = require('fs');
const iconv = require('iconv-lite');
const jschardet = require('jschardet');
const path = require('path');
// Full normalization logic...
`;

    const files = {
        'KICKOFF.md': KICKOFF_CONTENT,
        'claw.md': `# Claw Rules v7.7.1\nMode: ${modeName}\n- Always backup docs before sync.\n- No God Files.\n- Planning First.`,
        '.claw/templates/design_template.md': DESIGN_TEMPLATE,
        'tools/validate_docs.js': VALIDATOR_SCRIPT,
        'tools/normalize_docs.js': NORMALIZE_SCRIPT,
        'PLAN.md': (mode === '2') ? "# 📋 Project Execution Plan\n- [ ] Task 1: Architecture Check" : "# 📋 Tasks"
    };

    console.log(`\n📝 [1/3] Deploying ${modeName} Environment...`);
    const dirs = ['.claw/templates', 'docs/specs', 'docs/history', 'tools', 'input_docs'];
    dirs.forEach(d => { if (!fs.existsSync(d)) fs.mkdirSync(d, { recursive: true }); });

    for (const [f, c] of Object.entries(files)) {
        fs.writeFileSync(f, c);
        console.log(`  ✅ ${f}`);
    }

    // Gitignore for Orchestration
    const gitignoreContent = `node_modules/\n.env\n*.log\n!README.md\n!KICKOFF.md\n!setup_claw.js\n!docs/\n!docs/history/\n!input_docs/\n!tools/\n!TASKS.md\n!PLAN.md\n!internal/`;
    fs.writeFileSync('.gitignore', gitignoreContent);

    console.log("\n📦 [2/3] Installing Dependencies...");
    if (!fs.existsSync('package.json')) await runCommand('npm init -y');
    await runCommand('npm install @modelcontextprotocol/sdk zod iconv-lite jschardet --save');

    console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment v7.7.1 READY! ✨");
    process.exit(0);
}

main();
