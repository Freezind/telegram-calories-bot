# Feature Specification: LLM-Based Integration Testing

Feature Branch: 004-llm-bot-testing
Created: 2025-12-17
Status: Complete
Input: Fully automated integration testing for bot handlers and Mini App HTTP handlers using LLM as judge. No real Telegram API calls required.

---

## Goal

Deliver a fully automated integration testing system that:
- Tests bot handler logic WITHOUT calling real Telegram APIs
- Tests Mini App HTTP handlers using httptest
- Uses LLM (Gemini) as an automated judge for test evaluation
- Runs locally without requiring cloud deployment
- Produces required deliverables: LLM test prompts and a test report

This spec intentionally keeps scope minimal (integration-level, not E2E) and focuses on validating handler correctness and message preservation behavior.

---

## Clarifications

### Session 2025-12-17

- Q: Which cloud platform should be used for deployment? → A: Railway / Render
- Q: What persistence strategy should be used for cloud deployment? → A: Keep in-memory only; accept data loss on container restart (demo-grade acceptable)
- Q: How should the test runner behave when individual test scenarios fail? → A: Continue and mark failed
- Q: What format should the LLM judge use for returning PASS/FAIL results? → A: Structured JSON format
- Q: Where should captured test evidence be stored? → A: Embedded in the markdown report

---

## In Scope

1) Integration testing infrastructure:
- MockBot for capturing sent messages and tracking deletions
- mockContextAdapter implementing telebot.Context interface
- Direct handler invocation (no Telegram API calls)

2) Bot handler integration tests (S1-S4):
- /start command validation
- /estimate + photo upload validation
- Re-estimate button message preservation test
- Cancel button message preservation test

3) Mini App HTTP handler integration tests (S5):
- CRUD operations on /api/logs using httptest.NewServer
- Fake auth middleware for user scoping
- List, Create, Update, Delete log validation

4) LLM judge integration:
- Gemini 2.5 Flash for test evaluation
- Structured JSON output (verdict + rationale)
- All prompts archived verbatim to prompts.md

5) Deliverables:
- Self-contained test report at reports/004-test-report.md
- LLM test prompts appended to prompts.md
- Exit code 0 (all pass) or 1 (any fail)

---

## Out of Scope

- Real Telegram API calls (tests use mocks)
- Cloud deployment (tests run locally)
- End-to-end testing via browser automation (using integration tests instead)
- Actual user button click simulation (Bot API limitation)
- Telegram user client integration
- Accuracy validation of calorie estimation (structure validation only)
- Load/performance testing
- CI/CD pipeline integration
- Screenshot capture
- Test result history tracking

---

## User Scenarios and Testing

### User Story 1: Bot onboarding test (P1)
As a developer, I want to verify /start response so onboarding regressions are caught.

Acceptance scenarios:
1. Given the bot is deployed and reachable, when a user sends /start, then the bot replies with an introduction message
2. Given the response is captured, when evaluated by the LLM judge, then the judge returns PASS or FAIL with rationale

---

### User Story 2: Bot estimate flow test (P1)
As a developer, I want to verify image estimate produces a structured response.

Acceptance scenarios:
1. Given the bot is in awaiting-image state, when a test food image is uploaded, then the bot replies with an estimate message
2. Given the estimate message is captured, when evaluated by the LLM judge, then judge returns PASS if structure is present (foods, calories, confidence), without validating accuracy

---

### User Story 3: Bot inline buttons preservation test (P1)
As a developer, I want to verify re-estimate and cancel do not delete or edit previous estimate messages.

Acceptance scenarios:
1. Given an estimate message exists with inline buttons, when re-estimate is clicked, then the bot sends a NEW message and the previous estimate remains unchanged
2. Given an estimate message exists with inline buttons, when cancel is clicked, then the bot sends a NEW message and the previous estimate remains unchanged
3. Given captured evidence, when evaluated by the LLM judge, then judge returns PASS only if message preservation is confirmed

---

### User Story 4: Mini App view logs test (P1)
As a developer, I want to verify the deployed Mini App loads and shows either logs or a clear empty state.

Acceptance scenarios:
1. Given the Mini App is deployed, when the page is opened in a browser, then the UI renders successfully and shows either a log table or a No logs yet empty state
2. Given page evidence is captured, when evaluated by the LLM judge, then judge returns PASS if required UI elements are present

---

### User Story 5: Mini App create log test (P1)
As a developer, I want to verify CRUD works at least for Create so the Mini App is not a static page.

Acceptance scenarios:
1. Given the Mini App is open, when Add New Log is used to submit valid data, then a new row appears in the UI
2. Given page evidence is captured, when evaluated by the LLM judge, then judge returns PASS if the new entry is visible

Note: Edit and Delete are out of scope for this spec to keep it minimal.

---

## Requirements

### Deployment requirements
- FR-001: The system MUST be deployed to the cloud with public HTTPS access
- FR-002: The Mini App URL MUST be a public HTTPS URL (no localhost, no tunnel-only URL)
- FR-003: The bot MUST be configured so the Menu Button points to the deployed Mini App URL
- FR-004: The backend MUST be reachable by the Mini App in the deployed environment

### Testing requirements
- FR-005: A single command MUST run all automated tests and generate a report
- FR-006: Tests MUST cover bot and mini app scenarios defined above (US1 to US5)
- FR-007: The test runner MUST capture evidence:
  - Bot: message text, message IDs, callback data, timestamps
  - Mini App: page URL, key visible texts, screenshots
- FR-008: The test runner MUST call an LLM judge (Gemini allowed) which returns structured JSON format `{"verdict": "PASS"|"FAIL", "rationale": "..."}` per scenario
- FR-009: The test report MUST preserve original error messages verbatim (do not rewrite or genericize)
- FR-010: All LLM judge prompts used MUST be saved verbatim into repo root prompts.md (constitution requirement)
- FR-011: When a test scenario fails (timeout, error, etc.), the test runner MUST continue executing remaining scenarios and mark the failed scenario with error details in the report

### Tooling constraints
- FR-012: Mini App GUI automation MUST use a browser automation tool (example: Playwright) against the deployed URL
- FR-013: Bot automation MUST use Telegram Bot API interactions to drive scenarios and capture outputs

---

## Success Criteria

- SC-001: Deployed bot is usable with a bot name and can respond to /start and estimate flows
- SC-002: Deployed Mini App loads successfully over HTTPS and can create a log entry
- SC-003: Running one test command produces a markdown report with PASS or FAIL for all scenarios
- SC-004: The report includes captured evidence embedded inline (bot message text as code blocks, screenshots as base64 data URIs)
- SC-005: prompts.md contains all LLM judge prompts used, saved verbatim

---

## Deliverables

1) Public deployment URLs
- Mini App URL (HTTPS)
- Backend API base URL (HTTPS) if separate
- Bot name usable in Telegram

2) LLM test tool prompts
- Saved verbatim in repo root prompts.md

3) Test report
- A self-contained generated markdown file (for example: reports/004-test-report.md)
- Contains scenario list, steps, evidence embedded inline (text as code blocks, screenshots as base64 data URIs), PASS or FAIL, judge rationale, timestamps

4) One-command runner
- Example: make test-e2e
- Must run tests and generate the report

---

## Assumptions

- Deployment will use Railway or Render (container-based platform with free tier)
- Gemini API key is available for the judge
- A test Telegram account is available to run bot scenarios
- A stable test food image exists in the repo for estimate testing
- The Mini App can be tested via browser E2E on deployed URL without requiring Telegram app context

---

## Notes to keep scope minimal

- Only test Create for Mini App CRUD in this spec
- Bot estimate correctness is not evaluated, only response structure and flow integrity
- Focus on reproducible evidence and a readable report, not perfect automation sophistication

Status: Ready for speckit.plan