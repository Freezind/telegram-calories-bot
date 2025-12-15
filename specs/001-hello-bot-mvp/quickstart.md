# Quickstart: Hello Bot MVP

**Feature**: 001-hello-bot-mvp
**Last Updated**: 2025-12-15

## Prerequisites

- Go 1.21 or higher installed
- Telegram account
- Internet connection

## Step 1: Get a Telegram Bot Token

1. Open Telegram and search for `@BotFather`
2. Start a chat and send `/newbot`
3. Follow prompts to choose a name and username for your bot
4. BotFather will provide a bot token that looks like: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`
5. **Keep this token secret** - it gives full access to your bot

## Step 2: Set Environment Variable

Export the bot token as an environment variable:

**Linux/macOS:**
```bash
export TELEGRAM_BOT_TOKEN="your-token-here"
```

**Windows (PowerShell):**
```powershell
$env:TELEGRAM_BOT_TOKEN="your-token-here"
```

**Windows (Command Prompt):**
```cmd
set TELEGRAM_BOT_TOKEN=your-token-here
```

## Step 3: Initialize Go Module

From the repository root:

```bash
go mod init telegram-calories-bot
go get github.com/tucnak/telebot/v3
```

This creates `go.mod` and downloads the Telebot dependency.

## Step 4: Run the Bot

From the repository root:

```bash
go run cmd/bot/main.go
```

You should see output like:
```
2025/12/15 12:00:00 Bot started successfully
```

## Step 5: Test the Bot

1. Open Telegram and search for your bot by username
2. Send `/start` - bot should reply with "Hello! 👋"
3. Send `hello` (lowercase) - bot should reply with "Hello! 👋"
4. Try other messages - bot should ignore them (no response)

Check your terminal - you should see logs like:
```
2025/12/15 12:01:23 [UPDATE] UserID=123456789 Text=/start
2025/12/15 12:01:30 [UPDATE] UserID=123456789 Text=hello
```

## Step 6: Run Tests

```bash
go test ./tests/unit/... -v
```

Expected output:
```
=== RUN   TestGreet
--- PASS: TestGreet (0.00s)
PASS
ok      telegram-calories-bot/tests/unit    0.001s
```

## Stopping the Bot

Press `Ctrl+C` in the terminal running the bot. It will shut down gracefully.

## Troubleshooting

### "Failed to create bot: 401 Unauthorized"

- Check that TELEGRAM_BOT_TOKEN environment variable is set correctly
- Verify the token with BotFather (send `/token` to @BotFather)
- Ensure no extra spaces or quotes in the token

### "command not found: go"

- Install Go 1.21+ from https://go.dev/dl/
- Verify installation: `go version`

### "package github.com/tucnak/telebot/v3 is not in GOROOT"

- Run `go mod tidy` to download dependencies
- Check internet connection

### Bot doesn't respond to messages

- Verify bot is running (check terminal for "Bot started successfully")
- Check logs for UPDATE entries - if none, Telegram isn't sending updates
- Try stopping and restarting the bot
- Ensure no other instance of the bot is running

### Tests fail with "Hello! 👋" mismatch

- Verify the emoji is correctly copied (waving hand emoji: U+1F44B)
- Check file encoding is UTF-8

## What's Next?

This MVP validates basic bot functionality. Future features will add:
- Calorie calculation via LLM (image processing)
- Mini App interface for historical data
- Database persistence
- Deployment automation

For now, the bot demonstrates:
- ✅ Successful connection to Telegram
- ✅ Command handling (/start)
- ✅ Text message handling ("hello")
- ✅ Exact response matching ("Hello! 👋")
- ✅ Logging with user IDs
- ✅ Automated testing

## File Structure Reference

```
telegram-calories-bot/
├── cmd/bot/main.go              # Start here - bot entry point
├── internal/
│   ├── handlers/greeting.go     # Telegram message handlers
│   └── services/greeter.go      # Business logic
├── tests/unit/greeter_test.go   # Unit tests
├── go.mod                        # Dependencies
└── .env.example                  # Token documentation
```

## Environment Variable Alternative

Instead of exporting TELEGRAM_BOT_TOKEN each time, you can create a `.env` file:

```bash
# .env (DO NOT COMMIT THIS FILE)
TELEGRAM_BOT_TOKEN=your-token-here
```

Then load it before running:

```bash
source .env  # Linux/macOS
go run cmd/bot/main.go
```

**Important**: Add `.env` to `.gitignore` to avoid committing secrets.
