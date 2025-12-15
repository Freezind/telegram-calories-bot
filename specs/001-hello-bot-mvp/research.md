# Research: Hello Bot MVP

**Feature**: 001-hello-bot-mvp
**Date**: 2025-12-15
**Status**: Complete

## Research Questions

This MVP has minimal unknowns since the spec is explicit about scope and the technology stack is fixed (Golang + Telebot).

### Q1: How to initialize a Telebot bot in Go?

**Decision**: Use `telebot.NewBot()` with bot token and polling settings.

**Rationale**: The Telebot v3 library provides a simple constructor that accepts settings including bot token. Polling mode (vs webhooks) is simpler for MVP as it requires no server setup.

**Example Pattern**:
```go
bot, err := telebot.NewBot(telebot.Settings{
    Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
    Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
})
```

**Alternatives Considered**:
- Webhook mode: Rejected for MVP - requires HTTPS server and public URL
- Manual API calls via net/http: Rejected - Telebot framework handles this

### Q2: How to handle /start command in Telebot?

**Decision**: Use `bot.Handle("/start", handlerFunc)` to register command handler.

**Rationale**: Telebot provides built-in command routing. Commands are prefixed with "/" and automatically parsed.

**Example Pattern**:
```go
bot.Handle("/start", func(c telebot.Context) error {
    return c.Send("Hello! 👋")
})
```

**Alternatives Considered**:
- Manual text parsing: Rejected - Telebot handles command parsing
- Middleware-based routing: Overkill for single command MVP

### Q3: How to handle plain text messages in Telebot?

**Decision**: Use `bot.Handle(telebot.OnText, handlerFunc)` to catch all text messages, then filter for "hello".

**Rationale**: Telebot's `OnText` event fires for all text messages. We can check the message content and respond only to exact match "hello".

**Example Pattern**:
```go
bot.Handle(telebot.OnText, func(c telebot.Context) error {
    if c.Text() == "hello" {
        return c.Send("Hello! 👋")
    }
    return nil // Ignore other messages
})
```

**Alternatives Considered**:
- Regex-based matching: Overkill for exact string match
- Custom middleware: Not needed for simple equality check

### Q4: How to log user ID and update info in Go?

**Decision**: Use Go's standard `log` package with structured format.

**Rationale**: For MVP, standard library logging to stdout is sufficient. Can upgrade to structured logging (zerolog/zap) in future features.

**Example Pattern**:
```go
log.Printf("[UPDATE] UserID=%d Text=%s", c.Sender().ID, c.Text())
```

**Alternatives Considered**:
- Structured logging (zerolog, zap): Deferred to future - overkill for MVP
- File logging: Not needed - stdout sufficient for MVP

### Q5: How to test Telegram bot handlers without actual Telegram API?

**Decision**: Extract greeting logic into a pure service function, test that. Handler is thin wrapper.

**Rationale**: Following constitution's separation of concerns - handlers call services, services contain testable logic. For MVP, testing the greeting service is sufficient.

**Example Pattern**:
```go
// Service (easily testable)
func Greet() string {
    return "Hello! 👋"
}

// Handler (thin wrapper)
func HandleGreeting(c telebot.Context) error {
    return c.Send(Greet())
}

// Test
func TestGreet(t *testing.T) {
    got := Greet()
    want := "Hello! 👋"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

**Alternatives Considered**:
- Mocking Telebot context: Complex for MVP, not worth the effort
- Integration tests with test bot: Deferred to future features

## Technology Decisions

### Bot Token Management

**Decision**: Use environment variable `TELEGRAM_BOT_TOKEN` loaded via `os.Getenv()`.

**Rationale**: Simple, secure (token not in code), and standard practice. Provide `.env.example` for documentation.

**Constraints**: User must set TELEGRAM_BOT_TOKEN before running the bot.

### Dependency Management

**Decision**: Use Go modules (go.mod) with only one external dependency: `github.com/tucnak/telebot/v3`.

**Rationale**: Go modules is the standard Go dependency manager. Telebot v3 is mandated by constitution.

**Constraints**: Requires Go 1.21+ as specified in Technical Context.

### Error Handling

**Decision**: Log errors and exit gracefully on initialization failures. For runtime errors in handlers, log and continue.

**Rationale**: Bot should fail fast if token is invalid or Telegram API unreachable. But individual message handling errors shouldn't crash the bot.

**Example Pattern**:
```go
// Initialization - fail fast
bot, err := telebot.NewBot(settings)
if err != nil {
    log.Fatalf("Failed to create bot: %v", err)
}

// Runtime - log and continue
bot.Handle("/start", func(c telebot.Context) error {
    err := c.Send("Hello! 👋")
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }
    return err
})
```

## Implementation Notes

### What We're NOT Building (Out of Scope)

The following are explicitly excluded per spec:
- Database or any persistence
- Configuration files (only environment variable for token)
- Mini App interface
- LLM integration
- Deployment scripts or Docker
- Additional commands beyond /start
- Inline buttons or keyboards
- Webhook mode
- Rate limiting or retry logic
- Production monitoring

### Minimal Viable Scope

The implementation should contain ONLY:
1. Bot initialization with token from env var
2. /start command handler
3. Text message handler (filters for "hello")
4. Logging (user ID + message text)
5. One unit test for greeting service
6. go.mod with Telebot dependency
7. .env.example documenting TELEGRAM_BOT_TOKEN requirement

Total expected lines of code: ~80-100 lines across 4 files.

## Next Steps

Phase 1 will produce:
- `quickstart.md`: How to get a bot token, set TELEGRAM_BOT_TOKEN, run the bot
- No data-model.md (no entities)
- No contracts/ (no API endpoints)
