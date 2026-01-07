
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

function validate() {
    let errors = [];
    console.log("🛡️ [Claw v8.3 Governance] Running Absolute Audit...");

    // 1. Dependency/Isolation Check (Domain Purity)
    const domainDir = 'internal/domain';
    if (fs.existsSync(domainDir)) {
        const domains = fs.readdirSync(domainDir, {recursive: true}).filter(f => f.endsWith('.go'));
        domains.forEach(f => {
            const c = fs.readFileSync(path.join(domainDir, f), 'utf-8');
            if (c.includes('internal/usecase') || c.includes('internal/infrastructure')) 
                errors.push(`ARCH_VIOLATION: Domain layer '${f}' imports upper layers!`);
        });
    }

    // 2. Package Spec Sync Check
    const internalDir = 'internal';
    if (fs.existsSync(internalDir)) {
        const getPkgs = (dir) => {
            let res = [];
            fs.readdirSync(dir).forEach(f => {
                const p = path.join(dir, f);
                if (fs.statSync(p).isDirectory()) {
                    if (fs.readdirSync(p).some(file => file.endsWith('.go') || file.endsWith('.py') || file.endsWith('.ts'))) res.push(f);
                    res = res.concat(getPkgs(p));
                }
            });
            return res;
        };
        const pkgs = getPkgs(internalDir);
        pkgs.forEach(pkg => {
            const specExists = fs.readdirSync('docs/specs').some(s => s.toLowerCase().includes(pkg.toLowerCase()) || s === 'salesforce_infra_spec.md');
            if (!specExists) errors.push(`MISSING_SPEC: Package '${pkg}' lacks documentation in docs/specs/.`);
        });
    }

    // 3. Self-Spec Mandatory Check
    ['setup_claw_spec.md', 'validate_claw_spec.md'].forEach(s => {
        if (!fs.existsSync(path.join('docs/specs', s))) errors.push(`MISSING_CORE_SPEC: ${s} is mandatory.`);
    });

    // 4. PLAN.md Status
    if (fs.existsSync('PLAN.md')) {
        const plan = fs.readFileSync('PLAN.md', 'utf-8');
        if (plan.includes('[ ]')) errors.push("UNFINISHED_TASKS: PLAN.md contains unfinished check-boxes [ ].");
    }

    if (errors.length > 0) {
        console.error("\n❌ GOVERNANCE FAILED:");
        errors.forEach(e => console.error("  - " + e));
        process.exit(1);
    } else {
        console.log("\n✅ ALL GREEN: Systems are in perfect sync.");
    }
}
validate();
