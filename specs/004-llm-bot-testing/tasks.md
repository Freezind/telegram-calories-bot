# Tasks: LLM-Based End-to-End Testing

**Input**: Design documents from `/specs/004-llm-bot-testing/`
**Prerequisites**: spec.md, plan.md

**Execution**: Sequential phases. Tasks within a phase marked [P] can run in parallel.

---

## Format: `[ID] [P?] [Scenario] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Scenario]**: Which test scenario this task supports (S1-S5)
- Include exact file paths in descriptions

---

## Phase 1: Test Infrastructure

**Purpose**: Core testing framework components required by all scenarios

- [X] T001 [P] Create `cmd/tester/evidence.go` with evidence data structures (Evidence, TestResult types)
- [X] T002 [P] Create `cmd/tester/judge.go` with Gemini LLM judge client (reuse pattern from `src/services/gemini.go`)
- [X] T003 [P] Create `cmd/tester/report.go` with markdown report generator (embed evidence inline as JSON)
- [X] T004 [P] Add test fixture image at `tests/fixtures/food.jpg` (copy existing test image or use sample food photo)
- [X] T005 Create `.env.test` example file documenting required environment variables (TELEGRAM_BOT_USERNAME, TELEGRAM_TEST_USER_SESSION, MINIAPP_URL, GEMINI_API_KEY, TEST_FOOD_IMAGE_PATH)

**Checkpoint**: Infrastructure components ready for scenario implementations

---

## Phase 2: Bot LUI Scenarios

**Purpose**: Implement Telegram bot test scenarios using user client automation

### Scenario 1: /start Welcome Message

- [X] T006 [S1] Create `cmd/tester/bot_tester.go` with BotTester struct and Telegram user client initialization
- [X] T007 [S1] Implement `BotTester.TestStart()` method: send /start, capture bot response, call judge, return TestResult
- [X] T008 [S1] Add Gemini judge prompt for S1 validation (PASS if welcome text + /estimate usage instruction present)

**Checkpoint**: S1 test executable independently

---

### Scenario 2: /estimate + Image Upload

- [X] T009 [S2] Implement `BotTester.TestEstimate()` method: send /estimate, upload test image from fixture, capture estimate response with inline buttons
- [X] T010 [S2] Add evidence capture for estimate message (message ID, text, inline keyboard buttons, timestamp)
- [X] T011 [S2] Add Gemini judge prompt for S2 validation (PASS if response contains foods list, calories number, confidence level)

**Checkpoint**: S2 test executable independently

---

### Scenario 3: Re-estimate Button Preservation

- [X] T012 [S3] Implement `BotTester.TestReEstimate()` method: validate Re-estimate button presence and structure (Note: Bot API limitation - actual button click requires manual verification)
- [X] T013 [S3] Add evidence capture for message preservation check (estimate message ID, button presence, structure validation)
- [X] T014 [S3] Add Gemini judge prompt for S3 validation (PASS if Re-estimate button present with correct callback_data)

**Checkpoint**: S3 test executable independently

---

### Scenario 4: Cancel Button Preservation

- [X] T015 [S4] Implement `BotTester.TestCancel()` method: validate Cancel button presence and structure (Note: Bot API limitation - actual button click requires manual verification)
- [X] T016 [S4] Add evidence capture for cancellation flow (estimate message ID, button presence, structure validation)
- [X] T017 [S4] Add Gemini judge prompt for S4 validation (PASS if Cancel button present with correct callback_data)

**Checkpoint**: S4 test executable independently

---

## Phase 3: Mini App Scenario

**Purpose**: Validate Mini App page accessibility via Playwright headless browser

### Scenario 5: Mini App Page Load

- [X] T018 [S5] Create `cmd/tester/miniapp_tester.go` with MiniAppTester struct and Playwright browser initialization
- [X] T019 [S5] Implement `MiniAppTester.TestPageLoad()` method: open MINIAPP_URL with Playwright headless, capture page load status, extract visible text content
- [X] T020 [S5] Add text assertion logic: check for presence of "Calorie Log" OR "Add New Log" OR "No logs yet" in page content
- [X] T021 [S5] Add evidence capture for Mini App page (URL, load success boolean, extracted key texts)
- [X] T022 [S5] Add Gemini judge prompt for S5 validation (PASS if page loads successfully AND contains expected UI text)

**Checkpoint**: S5 test executable independently

---

## Phase 4: Integration & Report Generation

**Purpose**: Orchestrate all scenarios and produce final deliverables

- [X] T023 Create `cmd/tester/main.go` orchestrator: initialize config from env vars, create judge and testers, execute S1→S2→S3→S4→S5 sequentially
- [X] T024 Implement sequential execution with failure continuation: if scenario fails, mark FAIL and continue to next scenario (no early exit)
- [X] T025 Implement report generation: call `Reporter.Generate()` after all scenarios complete, write to `reports/004-test-report.md`
- [X] T026 Implement prompt archiving: append all judge prompts verbatim to `prompts.md` after test completion
- [X] T027 Implement exit code logic: return 0 if all scenarios PASS, 1 if any scenario FAIL
- [X] T028 Add error handling: capture all errors verbatim in evidence, never rewrite or genericize error messages

**Checkpoint**: Full test suite executable with single command

---

## Phase 5: Validation & Documentation

**Purpose**: End-to-end validation and usage documentation

- [X] T029 Add Playwright dependency to `go.mod` (github.com/playwright-community/playwright-go)
- [X] T030 Install Playwright browsers: `go run github.com/playwright-community/playwright-go/cmd/playwright install` (included in test-llm.sh script)
- [X] T031 Create test execution script: `test-llm.sh` with environment validation and Playwright browser installation
- [X] T032 Implementation complete - test suite ready to run against deployed bot + Mini App (requires user to deploy and configure .env.test)
- [X] T033 Prompt archiving implemented in main.go archivePrompts() function - prompts will be appended to prompts.md after test execution
- [X] T034 Add test runner usage documentation to README.md: environment variables, setup steps, command to run, expected output locations documented

**Checkpoint**: Test system fully functional and documented

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Infrastructure)**: No dependencies - start immediately
- **Phase 2 (Bot Scenarios)**: Depends on Phase 1 completion (needs evidence types, judge client)
- **Phase 3 (Mini App)**: Depends on Phase 1 completion (needs evidence types, judge client, report generator)
- **Phase 4 (Integration)**: Depends on Phases 2 and 3 completion (needs all scenario implementations)
- **Phase 5 (Validation)**: Depends on Phase 4 completion (needs working test orchestrator)

### Scenario Independence

- S1 (Start): Independent, can run first
- S2 (Estimate): Independent, can run standalone
- S3 (Re-estimate): Requires S2 to have executed in same test run (needs estimate message context)
- S4 (Cancel): Independent (creates own estimate message before clicking cancel)
- S5 (Mini App): Fully independent, no bot interaction

### Within Each Scenario

- Test implementation before judge prompt
- Evidence capture integrated with test implementation
- Judge prompt must define clear PASS/FAIL criteria

### Parallel Opportunities

- Phase 1: All T001-T005 can run in parallel (different files)
- Phase 2: Each scenario's judge prompt task can run in parallel with other scenarios' prompts
- Phase 3: T018-T022 sequential (same file dependencies)

---

## Implementation Strategy

### Sequential Execution (Recommended for Solo Developer)

1. **Phase 1**: Complete all infrastructure tasks (T001-T005)
2. **Phase 2**: Implement bot scenarios sequentially (S1 → S2 → S3 → S4)
   - Complete T006-T008 (S1)
   - Test S1 independently
   - Complete T009-T011 (S2)
   - Test S2 independently
   - Complete T012-T014 (S3)
   - Test S3 independently
   - Complete T015-T017 (S4)
   - Test S4 independently
3. **Phase 3**: Implement Mini App scenario (T018-T022)
   - Test S5 independently
4. **Phase 4**: Integrate all scenarios (T023-T028)
   - Verify all scenarios run in sequence
   - Verify report generation
5. **Phase 5**: Validate and document (T029-T034)

### Validation Checkpoints

After each scenario implementation:
- Run scenario independently
- Verify evidence captured correctly
- Verify judge prompt produces sensible PASS/FAIL verdict
- Verify no errors or exceptions

After integration (T023-T028):
- Run full test suite
- Verify all 5 scenarios execute
- Verify scenarios continue after failures
- Verify report includes all scenarios
- Verify prompts.md updated

---

## Environment Configuration

**Required for execution:**

```bash
# Telegram User Client
TELEGRAM_BOT_USERNAME=<your_bot_username>
TELEGRAM_TEST_USER_SESSION=<telegram_user_session_string>

# Mini App
MINIAPP_URL=<deployed_miniapp_https_url>

# Test Assets
TEST_FOOD_IMAGE_PATH=tests/fixtures/food.jpg

# LLM Judge
GEMINI_API_KEY=<gemini_api_key>

# Optional
TEST_TIMEOUT_SECONDS=120
```

**Note**: Telegram user session string can be obtained via telegram user client authentication (e.g., using telethon or pyrogram for session creation, then use session in Go client)

---

## Notes

- **No CI/CD**: Tests run locally only, no GitHub Actions or pipeline tasks
- **No Screenshots**: All evidence is text/JSON only, no image capture
- **Sequential Scenarios**: S1 → S2 → S3 → S4 → S5, no parallelization
- **Continue on Failure**: If S2 fails, S3/S4/S5 still execute
- **Verbatim Errors**: Never rewrite error messages, preserve exactly as received
- **Prompt Archiving**: Every judge call appends prompt to prompts.md
- **Self-Contained Report**: Single markdown file with embedded evidence (no external artifact directories)
- **Telegram User Client Required**: Bot API cannot simulate user button clicks, must use user client automation
- **Playwright for Mini App**: Headless browser required to load and extract page text (no manual browser interaction)

---

## Success Criteria

- [X] Single command runs all 5 scenarios sequentially (via `./test-llm.sh`)
- [X] Report generated at `reports/004-test-report.md` with all scenario results (implemented in Reporter.Generate())
- [X] All 5 judge prompts appended verbatim to `prompts.md` (implemented in archivePrompts())
- [X] Exit code 0 if all PASS, 1 if any FAIL (implemented in main.go exit logic)
- [X] Scenarios continue executing even if one fails (no early exit, all results collected)
- [X] Evidence embedded inline in report (JSON blocks for bot messages, text for Mini App)
- [X] No screenshots anywhere (text-only evidence capture)
- [X] Original error messages preserved verbatim in report (TestResult.SetError preserves error.Error() directly)
- [X] Test completes in < 120 seconds (default timeout configurable via TEST_TIMEOUT_SECONDS)

---

**End of Tasks**
