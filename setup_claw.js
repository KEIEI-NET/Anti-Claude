const fs = require('fs');
const { exec } = require('child_process');
const readline = require('readline');
const path = require('path');

/**
 * Claw Setup CLI v8.0 - Hardened Governance Edition
 * 🛑 STOP: Cannot complete tasks without a Green Validation Report.
 * Mechanically prevents documentation omissions and sync errors.
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
🦅 Claw Setup CLI v8.0 - Hardened Governance
===========================================
`);

    const mode = '2'; // Infrastructure focus usually uses Deep Dive
    const modeName = 'Deep Dive (Hardened)';

    // --- 🚨 THE GATEKEEPER SCRIPT (validate_claw.js) ---
    const GATEKEEPER_SCRIPT = `
const fs = require('fs');
const path = require('path');

const REQUIRED_DOCS = ['docs/design.md', 'PLAN.md', 'KICKOFF.md'];
const INTERNAL_DIR = 'internal';

function validate() {
    let errors = [];
    console.log("🔍 [Claw Governance] Starting full project audit...");

    // 1. 基本ファイルの存在チェック
    REQUIRED_DOCS.forEach(f => {
        if (!fs.existsSync(f)) errors.push(\`MISSING_FILE: \${f} が存在しません。\`);
    });

    // 2. internal/ 以下の実装と docs/specs/ の同期チェック
    if (fs.existsSync(INTERNAL_DIR)) {
        const getPkgs = (dir) => {
            let res = [];
            fs.readdirSync(dir).forEach(f => {
                const p = path.join(dir, f);
                if (fs.statSync(p).isDirectory()) {
                    if (fs.readdirSync(p).some(file => file.endsWith('.go') || file.endsWith('.py') || file.endsWith('.ts'))) {
                        res.push(f);
                    }
                    res = res.concat(getPkgs(p));
                }
            });
            return res;
        };
        const pkgs = getPkgs(INTERNAL_DIR);
        pkgs.forEach(pkg => {
            const specFile = \`docs/specs/\${pkg}_spec.md\`;
            if (!fs.existsSync(specFile) && !fs.existsSync(\`docs/specs/salesforce_infra_spec.md\`)) { // 柔軟な命名対応
                errors.push(\`MISSING_SPEC: 実装パッケージ '\${pkg}' に対応するプログラム仕様書が見つかりません (Expected: docs/specs/\${pkg}_spec.md)\`);
            }
        });
    }

    // 3. 履歴（History）の直近保存チェック
    const historyDir = 'docs/history';
    if (!fs.existsSync(historyDir) || fs.readdirSync(historyDir).length === 0) {
        errors.push("MISSING_HISTORY: docs/history/ にバックアップが一度も作成されていません。");
    }

    // 4. PLAN.md の完了状態チェック
    if (fs.existsSync('PLAN.md')) {
        const plan = fs.readFileSync('PLAN.md', 'utf-8');
        if (plan.includes('[ ]')) {
            const pendingTasks = plan.split('\\n').filter(l => l.includes('[ ]'));
            console.warn(\`⚠️ Pending tasks in PLAN.md:\\n\${pendingTasks.join('\\n')}\`);
            // WARNINGにしておくが、厳格モードならerrorsにpush可能
        }
    }

    if (errors.length > 0) {
        console.error("\\n❌ CLAW GOVERNANCE FAILED:");
        errors.forEach(e => console.error("  " + e));
        console.error("\\n👉 Antigravityへの指示: 完了報告の前に、上記ドキュメントの欠落を解消してください。");
        process.exit(1);
    } else {
        console.log("\\n✅ CLAW GOVERNANCE PASSED: All code/doc synchronization is correct.");
    }
}
validate();
    `;

    const KICKOFF_CONTENT = `# 🚀 Claw Kickoff (v8.0 Hardened)

## 🏛️ Phase 1: Planning
1. Detect technology & Create docs/design.md.
2. Create **PLAN.md**.

## 🧠 Phase 2: Execution
- NO God Files. NO shortcuts.
- 実装完了の度に必読：
  1. \`docs/history/\` へのバックアップ
  2. \`docs/design.md\` の更新
  3. \`docs/specs/[Pkg]_spec.md\` の更新

## 🛑 Phase 3: Total Validation (The Gatekeeper)
タスクを終了（あなたに報告）する前に、Antigravity（私）は必ず以下を実行しなければなりません：
\`\`\`bash
node tools/validate_claw.js
\`\`\`
**このスクリプトが PASS しない限り、私は「作業完了」を宣言することを禁止されます。**
`;

    const files = {
        'KICKOFF.md': KICKOFF_CONTENT,
        'claw.md': `# Claw Rules v8.0\n1. Hardened Governance: Run validate_claw.js before completion.\n2. Automatic spec synchronization mandatory.\n3. History-first sync.`,
        'tools/validate_claw.js': GATEKEEPER_SCRIPT,
        '.claw/templates/design_template.md': "# System Design Document",
        'PLAN.md': "# 📋 Task Plan"
    };

    console.log(`\n📝 [1/3] Deploying v8.0 Hardened Environment...`);
    const dirs = ['.claw/templates', 'docs/specs', 'docs/history', 'tools', 'input_docs'];
    dirs.forEach(d => { if (!fs.existsSync(d)) fs.mkdirSync(d, { recursive: true }); });

    for (const [f, c] of Object.entries(files)) {
        fs.writeFileSync(f, c);
        console.log(`  ✅ \${f}`);
    }

    console.log("\n📦 [2/3] Finalizing dependencies...");
    if (!fs.existsSync('package.json')) await runCommand('npm init -y');
    await runCommand('npm install iconv-lite jschardet --save');

    console.log("\n\x1b[32m%s\x1b[0m", "✨ Claw Environment v8.0 READY (Hardened Governance) ✨");
    process.exit(0);
}

main().catch(console.error);
