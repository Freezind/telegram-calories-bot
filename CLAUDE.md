# telegram-calories-bot Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-15

## Active Technologies
- Golang 1.21+ + github.com/tucnak/telebot (v3), Google Gemini SDK for Go (002-calorie-image-estimate)
- N/A (stateless operation, in-memory session state only) (002-calorie-image-estimate)
- Go 1.21+ (backend), JavaScript/TypeScript (React 18+ frontend) (003-miniapp-mvp)
- In-memory only (sync.Map or mutex-protected map); LogStorage interface for future D1 migration (003-miniapp-mvp)

- Golang 1.21+ + github.com/tucnak/telebot (v3) (001-hello-bot-mvp)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Golang 1.21+

## Code Style

Golang 1.21+: Follow standard conventions

## Recent Changes
- 003-miniapp-mvp: Added Go 1.21+ (backend), JavaScript/TypeScript (React 18+ frontend)
- 002-calorie-image-estimate: Added Golang 1.21+ + github.com/tucnak/telebot (v3), Google Gemini SDK for Go

- 001-hello-bot-mvp: Added Golang 1.21+ + github.com/tucnak/telebot (v3)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
