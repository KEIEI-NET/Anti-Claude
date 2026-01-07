# Claw - Antigravity & Claude Code 連携ルール

## 開発モード: 🛡️ Deep Dive Mode (Clean Arch & DDD)

## 1. アーキテクチャ原則
- **Clean Architecture** & **DDD** & **SOLID原則** を遵守。

## 2. 役割と責任
- **Antigravity**: Architect, Domain Expert, Frontend.
- **Claude Code**: Backend Implementation (SOLID compliant).

## 3. まずはじめに (Getting Started)
**KICKOFF.md を参照し、適切なコマンドをAntigravityに送信してください。**

## 4. ワークフロー
### Phase 0: Domain Analysis & Design 🏛️
1. **Kickoff**: User sends "Import" or "New Design" command (See KICKOFF.md).
2. **Normalization**: If importing, Antigravity runs `node tools/normalize_docs.js`.
3. **Domain Modeling**: Antigravity analyzes requirements using **DDD**.
4. **Specification**: 
   - Create `docs/design.md`.
   - Create `docs/specs/[Component].md`.
5. **Approval**: User MUST approve models & specs.

### Phase 1: Implementation (SOLID Principles)
- **Domain Layer**: Implement Pure Domain Logic.
- **Application Layer**: Implement Use Cases.
- **Interface/Infra**: Adapters & DB.

## 5. ステータス
- **MCP Status**: Active
- **Template System**: Enabled (Clean Arch/DDD)
