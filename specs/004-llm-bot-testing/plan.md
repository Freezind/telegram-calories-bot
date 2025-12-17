# Implementation Plan: LLM-Based End-to-End Testing

Feature: 004-llm-bot-testing  
Created: 2025-12-17  
Status: Draft

---

## Overview

Build a minimal, demo-grade LLM-based testing system that validates Telegram bot behavior (LUI + Mini App) using Gemini as an LLM judge and produces a single human-readable markdown test report.

Core principle: simplest possible implementation. Sequential scenarios. No CI. No screenshots.

---

## Execution Constraints

- Tests MUST run locally only.
- No CI/CD integration (no GitHub Actions, no pipelines).
- Scenarios run sequentially. If one fails, continue running the rest and mark FAIL.

---

## Architecture

```
cmd/tester/
├── main.go
├── bot_tester.go
├── miniapp_tester.go
├── judge.go
├── report.go
└── evidence.go
```

Key decisions:
- Language: Go
- LLM judge: Gemini
- Mini App check: Playwright headless, assert key text only
- No screenshots anywhere

---

## Test Scenarios

### Scenario 1: /start welcome message
- Send /start
- Capture bot response text
- PASS if welcome + usage instructions exist

### Scenario 2: /estimate + image upload
- Send /estimate
- Upload test food image
- PASS if response contains foods list, calories number, confidence label

### Scenario 3: Re-estimate button preservation
- Click Re-estimate via user client automation
- Previous estimate message must remain
- New prompt message must appear

### Scenario 4: Cancel button preservation
- Click Cancel via user client automation
- Previous estimate message must remain
- Cancellation confirmation message must appear

### Scenario 5: Mini App page load
- Open MINIAPP_URL with Playwright
- PASS if page loads and contains key text:
  - Calorie Log
  - Add New Log
  - No logs yet

---

## Telegram Button Click Strategy

Inline button clicks must be executed via a Telegram user client.
Bot API cannot synthesize callback clicks.

User client must:
- send messages
- receive bot messages
- click inline buttons
- capture resulting messages and callback data

---

## Evidence Rules

- Text-only evidence
- JSON blocks for bot messages and callbacks
- Raw errors preserved verbatim
- No screenshots

---

## Report Output

Single markdown file including:
- summary table
- per-scenario details
- embedded evidence
- LLM verdict + rationale
- timestamps

---

## Configuration (Environment Variables)

Required:
- GEMINI_API_KEY
- MINIAPP_URL
- TEST_FOOD_IMAGE_PATH
- TELEGRAM_BOT_USERNAME
- TELEGRAM_TEST_USER_SESSION

Optional:
- TEST_TIMEOUT_SECONDS

---

## Success Criteria

- One command runs all tests
- Markdown report generated
- Prompts appended verbatim to prompts.md
- No screenshots
- Raw errors preserved
- Continue on failure

---

End of plan.
