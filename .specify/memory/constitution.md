<!--
Sync Impact Report:
- Version change: [none] → 1.0.0
- Modified principles: N/A (initial creation)
- Added sections: All sections (initial constitution)
- Removed sections: N/A
- Templates requiring updates:
  ✅ plan-template.md - Constitution Check section aligned
  ✅ spec-template.md - Requirements structure aligned
  ✅ tasks-template.md - Test-first workflow aligned
- Follow-up TODOs: None
-->

# Telegram Calories Bot Constitution

## Core Principles

### I. Quality-Controlled Vibe Coding

**Vibe coding** (natural language AI-assisted development) MUST be controlled through:
- Clear, testable acceptance criteria for every feature
- Mandatory code review gates before merge
- Automated linting and formatting (gofmt, golangci-lint)
- Explicit prompt engineering with success/failure examples
- Test coverage requirements (minimum 70% for core services)

**Rationale**: The project deliverables explicitly require "effective ways to control Claude Code's code quality." Vibe coding accelerates development but needs guardrails to prevent technical debt and ensure maintainability.

### II. Test-First with LLM Validation

**All features MUST follow this workflow**:
1. Write acceptance tests based on user scenarios (Given-When-Then format)
2. Ensure tests FAIL (red state)
3. Implement feature until tests PASS (green state)
4. Use LLM-based testing tools to validate bot interaction flows
5. Generate test reports as deliverables

**Rationale**: Requirements mandate "bot testing via LLM and related tools" and test reports as deliverables. Test-first ensures features meet requirements before implementation, and LLM tools validate conversational UX.

### III. Dual Interface Architecture

**The bot MUST maintain two distinct, decoupled interfaces**:
- **LUI (Language User Interface)**: Slash commands + inline buttons for conversational interaction
- **Mini App**: GUI-based CRUD operations for historical data management

**Both interfaces MUST**:
- Share the same data layer and business logic
- Operate independently (LUI works without Mini App and vice versa)
- Provide consistent calorie calculation results
- Support the same authentication/authorization context

**Rationale**: Requirements specify two user interaction modes. Decoupling ensures each can evolve independently while sharing core functionality.

### IV. Fixed Technology Stack

**Technology choices are NON-NEGOTIABLE**:
- **Language**: Golang (Go 1.21+)
- **Bot Framework**: github.com/tucnak/telebot (v3)
- **No alternative frameworks or languages permitted**

**Additional dependencies MUST**:
- Be justified in plan.md if adding new third-party libraries
- Prefer standard library solutions when possible
- Use Go modules for dependency management

**Rationale**: Requirements explicitly mandate Golang and the Telebot framework. Standardization reduces learning curve and ensures consistency.

### V. Deliverable-Driven Development

**Every feature implementation MUST produce**:
1. **Working Bot**: Deployable Telegram bot with accessible bot name
2. **Clean Code**: Well-structured, reviewed, and documented Go code
3. **Vibe Coding Prompts**: Complete prompt history used to generate the code
4. **LLM Test Prompts**: Prompts used for automated bot testing
5. **Test Reports**: Automated test results with pass/fail status

**Rationale**: Requirements list five specific deliverables. Development is incomplete until all deliverables are ready for handoff.

## Technical Standards

### Code Organization

**Project structure MUST follow**:
```
telegram-calories-bot/
├── cmd/bot/              # Main application entry point
├── internal/
│   ├── handlers/         # Telegram handlers (LUI)
│   ├── miniapp/          # Mini app logic (GUI)
│   ├── services/         # Business logic (calorie calculation, etc.)
│   ├── models/           # Data models
│   └── storage/          # Data persistence layer
├── tests/
│   ├── integration/      # LLM-based bot interaction tests
│   ├── unit/             # Standard unit tests
│   └── contract/         # API contract tests
├── prompts/              # Vibe coding prompts archive
└── docs/                 # Documentation
```

**Separation of concerns**:
- Handlers translate Telegram events → service calls
- Services contain pure business logic (testable without Telegram)
- Storage abstracts persistence (enables testing with mocks)
- Mini app and LUI handlers never share state directly

### Error Handling & Logging

**Error handling MUST**:
- Return errors explicitly (no panic in production code except initialization)
- Wrap errors with context using `fmt.Errorf("%w", err)` or errors package
- Log errors with structured logging (e.g., zerolog, zap)

**Logging MUST include**:
- User ID (for debugging user-specific issues)
- Operation type (e.g., "calorie_calculation", "data_query")
- Timestamp and severity level
- Never log sensitive data (images should be logged as metadata only)

### Image Processing & Calorie Calculation

**Calorie calculation via LLM MUST**:
- Use a defined prompt template (stored in prompts/)
- Include error handling for LLM API failures
- Cache results to avoid redundant API calls for same image hash
- Return confidence score or uncertainty indicators when applicable

**Image handling MUST**:
- Accept common formats (JPEG, PNG, WebP)
- Enforce size limits (e.g., max 10MB per image)
- Store images securely with access controls
- Provide image URLs in Mini App data views

## Development Workflow

### Pre-Implementation Checklist

Before starting any feature:
1. ✅ User scenarios written in spec.md (Given-When-Then format)
2. ✅ Acceptance criteria defined and testable
3. ✅ Constitution compliance verified
4. ✅ Technical approach documented in plan.md
5. ✅ Test stubs created and failing

### Implementation Cycle

For each user story:
1. **Red**: Write integration test that exercises full user journey (must fail)
2. **Green**: Implement minimum code to make test pass
3. **Refactor**: Clean up code while keeping tests green
4. **Review**: Run gofmt, golangci-lint, and manual code review
5. **Document**: Update prompts/ with AI prompts used
6. **Report**: Generate test report with results

### Code Review Gates

**No code merges without**:
- ✅ All tests passing (unit + integration)
- ✅ Linter passing with zero warnings
- ✅ Code coverage ≥70% for services and handlers
- ✅ Manual review by human or detailed AI review with approval
- ✅ Prompts archived in prompts/ directory

### Testing Requirements

**Integration tests MUST**:
- Use LLM-based testing tools to simulate user conversations
- Cover all slash commands and inline button flows
- Validate Mini App CRUD operations end-to-end
- Test error cases (invalid input, API failures, etc.)

**Unit tests MUST**:
- Test business logic in isolation (services layer)
- Mock external dependencies (LLM API, storage, Telegram API)
- Run fast (<100ms per test) to enable rapid iteration

## Governance

### Amendment Process

**Constitution changes require**:
1. Written proposal with rationale in GitHub issue or discussion
2. Impact analysis on existing code and workflows
3. Approval from project maintainer or team consensus
4. Version bump following semantic versioning
5. Update to all dependent templates (plan, spec, tasks)

### Version Semantics

- **MAJOR** (X.0.0): Remove or redefine core principles (e.g., drop test-first requirement)
- **MINOR** (x.Y.0): Add new principles or expand sections (e.g., add security requirements)
- **PATCH** (x.y.Z): Clarifications, typo fixes, non-semantic improvements

### Compliance Verification

**Every pull request MUST**:
- Include a constitution compliance checklist in PR description
- Verify no violations of core principles (I-V)
- Document any complexity additions in plan.md with justification
- Pass automated checks (tests, linting, coverage)

**Constitution supersedes all other practices**. If conflicts arise, constitution principles take precedence.

**Version**: 1.0.0 | **Ratified**: 2025-12-15 | **Last Amended**: 2025-12-15