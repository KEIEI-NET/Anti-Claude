const fs = require('fs');
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

const FORBIDDEN_PATTERNS = [
    { regex: /\(Original Diagram Here\)/, message: '物理構成図がプレースホルダーのままです。元の図を転記してください。' },
    { regex: /\.\.\./, message: '省略記号 "..." が残っています。詳細を記述してください。' },
    { regex: /TODO/, message: '"TODO" が残っています。' }
];

function validate() {
    if (!fs.existsSync(TARGET_FILE)) {
        console.error('❌ Error: docs/design.md が見つかりません。');
        process.exit(1);
    }

    const content = fs.readFileSync(TARGET_FILE, 'utf-8');
    let errors = [];

    // 1. Check Required Sections
    console.log('🔍 Checking for required sections...');
    REQUIRED_SECTIONS.forEach(section => {
        if (!content.includes(section)) {
            errors.push(`MISSING SECTION: "${section}" が見つかりません。`);
        }
    });

    // 2. Check Forbidden Patterns (Laziness Check)
    console.log('🔍 Checking for incompleteness...');
    FORBIDDEN_PATTERNS.forEach(pattern => {
        if (pattern.regex.test(content)) {
            errors.push(`LAZINESS DETECTED: ${pattern.message}`);
        }
    });

    // 3. Check Diagrams
    const codeBlocks = content.match(/```[a-z]*\s*[\r\n]+[\s\S]*?```/g) || [];
    const mermaidBlocks = codeBlocks.filter(block => block.includes('mermaid'));
    const textBlocks = codeBlocks.filter(block => 
        !block.includes('mermaid') && 
        (block.includes('+') || block.includes('|') || block.includes('┌') || block.includes('└'))
    );

    if (mermaidBlocks.length < 1) errors.push('MISSING DIAGRAM: Mermaidによる論理構成図がありません。');
    if (textBlocks.length < 1) errors.push('MISSING DIAGRAM: 物理構成図（ASCIIアート等）が見つかりません。');

    // Result
    if (errors.length > 0) {
        console.error('\n❌ VALIDATION FAILED: ドキュメントの品質基準を満たしていません。');
        errors.forEach(e => console.error(` - ${e}`));
        console.error('\n👉 AIへの指示: 上記のエラーを修正してから再提出してください。');
        process.exit(1);
    } else {
        console.log('\n✅ VALIDATION PASSED: ドキュメント品質は良好です。次のフェーズに進めます。');
    }
}

validate();