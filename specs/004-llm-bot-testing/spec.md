# Feature Specification: LLM-Based Bot Testing

**Feature Branch**: `004-llm-bot-testing`
**Created**: 2025-12-17
**Status**: Draft
**Input**: User description: "Build Spec 004: LLM-based automated testing for the Telegram bot (LUI), following the same \"small MVP\" style as `002-calorie-image-estimate/spec.md`."

## User Scenarios & Testing

### User Story 1 - Validate Bot Welcome Message (Priority: P1)

As a developer, I want to automatically verify that the bot's /start command shows the correct welcome message, so that I can catch regressions in the bot's first-impression UX.

**Why this priority**: The /start command is the entry point for all users. If this is broken, users cannot proceed with any other functionality.

**Independent Test**: Can be fully tested by sending `/start` to the bot and validating the response contains welcome text and usage instructions. Delivers confidence that new users see proper onboarding.

**Acceptance Scenarios**:

1. **Given** the bot is running, **When** a user sends `/start`, **Then** the bot responds with a welcome message containing bot introduction and usage instructions
2. **Given** the bot is running, **When** the LLM judge evaluates the /start response, **Then** the judge returns PASS with rationale explaining why the message meets welcome criteria

---

### User Story 2 - Validate Estimate Flow (Priority: P1)

As a developer, I want to automatically verify that uploading an image for /estimate produces a calorie estimate response, so that I can ensure the core bot functionality works end-to-end.

**Why this priority**: The /estimate flow is the primary value proposition of the bot. Without this working, the bot has no purpose.

**Independent Test**: Can be fully tested by triggering /estimate, uploading a test food image, and validating that a structured estimate response is returned (food items, calories, confidence). Delivers confidence that the core AI-powered estimation works.

**Acceptance Scenarios**:

1. **Given** the bot is in awaiting-image state after /estimate, **When** a user uploads a food image, **Then** the bot responds with an estimate message containing food items, calorie count, and confidence level
2. **Given** an estimate response is captured, **When** the LLM judge evaluates it against the rubric, **Then** the judge returns PASS if the response has reasonable structure (not validating calorie accuracy)

---

### User Story 3 - Validate Re-estimate Button (Priority: P2)

As a developer, I want to automatically verify that clicking the "Re-estimate" button prompts for another image without deleting the previous estimate, so that I can ensure users can iterate on estimates while preserving history.

**Why this priority**: Re-estimation is a key UX feature allowing users to get better estimates. Message preservation is critical for user trust and audit trails.

**Independent Test**: Can be fully tested by completing an estimate flow, clicking the "Re-estimate" button, and validating that (a) a new prompt appears and (b) the previous estimate message remains visible. Delivers confidence that iteration works without data loss.

**Acceptance Scenarios**:

1. **Given** an estimate message with "Re-estimate" button is displayed, **When** the user clicks "Re-estimate", **Then** the bot sends a NEW message asking for another image (previous estimate NOT deleted/edited)
2. **Given** the re-estimate interaction is captured, **When** the LLM judge evaluates it, **Then** the judge returns PASS if evidence shows both messages coexist

---

### User Story 4 - Validate Cancel Button (Priority: P3)

As a developer, I want to automatically verify that clicking the "Cancel" button cancels the estimation flow without deleting previous messages, so that I can ensure users can gracefully exit while preserving conversation history.

**Why this priority**: Cancel is a secondary UX feature for flow control. Less critical than core estimation, but important for completeness.

**Independent Test**: Can be fully tested by completing an estimate flow, clicking "Cancel", and validating that (a) a cancellation confirmation is sent and (b) the previous estimate message remains visible. Delivers confidence that cancellation works without data loss.

**Acceptance Scenarios**:

1. **Given** an estimate message with "Cancel" button is displayed, **When** the user clicks "Cancel", **Then** the bot sends a NEW cancellation message (previous estimate NOT deleted/edited)
2. **Given** the cancellation interaction is captured, **When** the LLM judge evaluates it, **Then** the judge returns PASS if evidence shows both messages coexist and cancellation is confirmed

---

### Edge Cases

- What happens when the bot is unreachable during test execution? (Test runner should capture connection errors verbatim and mark scenario as FAIL with the original error)
- What happens when the LLM judge API fails? (Test runner should retry once, then fail the scenario with the original API error preserved)
- What happens when callback data contains unexpected values? (Test runner should capture actual callback data verbatim and let the LLM judge evaluate it)
- What happens when the bot sends an error message instead of expected response? (Test runner should capture the exact error message and include it in the test report, marking scenario as FAIL)

## Requirements

### Functional Requirements

- **FR-001**: The test runner MUST execute all 4 test scenarios: /start validation, image upload estimation, re-estimate button click, and cancel button click
- **FR-002**: The test runner MUST capture bot message text and inline button callback data as evidence for each scenario
- **FR-003**: The test runner MUST call an LLM judge (Gemini or equivalent) with a strict rubric for each scenario and capture the PASS/FAIL verdict plus rationale
- **FR-004**: The test report MUST preserve original error messages from bot or LLM API without rewriting or genericizing them
- **FR-005**: All LLM judge prompts used in testing MUST be saved verbatim to `prompts.md` in the repository root
- **FR-006**: The test runner MUST generate a markdown report file containing: scenario name, execution steps, captured evidence, PASS/FAIL status, LLM judge rationale, and timestamp
- **FR-007**: A single command (e.g., `make test-llm` or `go test`) MUST run all test scenarios and generate the complete report
- **FR-008**: The LLM judge rubric MUST evaluate response structure and presence of expected elements (not semantic accuracy of calorie values)
- **FR-009**: The test runner MUST verify that re-estimate and cancel actions do NOT delete or edit previous messages (message preservation requirement)

### Key Entities

- **Test Scenario**: Represents one test case (e.g., "/start validation"). Contains: scenario ID, description, execution steps, expected behavior, captured evidence, PASS/FAIL status, LLM judge rationale, timestamp
- **Bot Interaction**: Represents a single bot message or callback. Contains: message text, callback data (if button click), timestamp, message ID
- **LLM Judge Prompt**: Represents the prompt sent to the LLM for evaluation. Contains: prompt template, scenario-specific context, rubric criteria, response format instructions
- **Test Report**: Represents the full test run output. Contains: list of Test Scenarios, overall summary (pass/fail counts), test run timestamp, environment info (bot version, etc.)

## Success Criteria

### Measurable Outcomes

- **SC-001**: Running the test command produces a complete markdown report file containing all 4 test scenarios with PASS/FAIL verdicts
- **SC-002**: The test report includes captured bot message text and callback data for each scenario as evidence
- **SC-003**: All LLM judge prompts used during testing are saved verbatim to `prompts.md` before test execution
- **SC-004**: Original error messages from bot or LLM API appear unmodified in the test report (no generic rewrites like "API error occurred")
- **SC-005**: Developers can run tests locally using a single command without additional setup beyond standard bot environment variables
- **SC-006**: Test execution completes within 5 minutes for all 4 scenarios (excluding LLM API latency)

## Scope

### In Scope

- Automated testing of 4 bot interaction scenarios using LLM as judge
- Capture of bot message text and callback data as test evidence
- Generation of markdown test report with PASS/FAIL per scenario
- Storage of LLM judge prompts in `prompts.md`
- Single-command test execution
- Preservation of original error messages end-to-end

### Out of Scope

- Mini App GUI automated testing (only bot LUI is tested)
- Accuracy validation of calorie estimation results (only structure validation)
- Continuous integration or deployment pipeline changes
- Performance or load testing of the bot
- Multi-language or internationalization testing
- Security or penetration testing

## Assumptions

- The bot is running and accessible during test execution (test environment setup is a prerequisite)
- Test scenarios will use a hardcoded test food image stored in the repository
- LLM judge will use Gemini API with the same `GEMINI_API_KEY` used by the bot
- Test scenarios execute sequentially (not in parallel) to avoid state conflicts
- Bot behavior matches current implementation (tests validate existing behavior, not ideal behavior)
- The test runner will use the existing Telegram bot API client libraries available in the codebase
- Developers have access to a test Telegram account/bot for running tests

## Dependencies

- Existing Telegram bot implementation (specs 001, 002, 003)
- Gemini API access (same as used for calorie estimation)
- Telegram Bot API client library already in use
- Go testing framework (standard library `testing` package)
- Markdown file generation capability (standard library or simple template)

## Constraints

- Test scenarios must not modify production bot state or data
- Test execution should be isolated and repeatable (use test-specific Telegram bot if needed)
- LLM judge prompts must be deterministic and version-controlled (saved to `prompts.md`)
- Test report format must be markdown for easy version control and diffing
- No external test frameworks beyond Go standard library (keep dependencies minimal)

## Non-Functional Requirements

- **Reliability**: Tests should be deterministic and produce consistent PASS/FAIL results when bot behavior is unchanged
- **Maintainability**: LLM judge prompts should be human-readable and easy to update when bot behavior intentionally changes
- **Debuggability**: Test reports must include enough evidence (captured messages, timestamps) to debug failures without re-running tests
- **Developer Experience**: Single command execution, clear PASS/FAIL output, and readable reports

## Open Questions

None - specification is complete based on provided requirements.
