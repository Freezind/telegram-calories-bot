# Telegram Calories Bot

A Telegram bot with Mini App for tracking calorie intake through images and manual entry.

## Features

- **Hello Bot MVP** (001): Basic bot greeting functionality
- **Calorie Image Estimate** (002): AI-powered calorie estimation from food images using Google Gemini
- **Mini App MVP** (003): Web-based Mini App for viewing and managing calorie logs

## Project Structure

```
telegram-calories-bot/
├── cmd/
│   ├── bot/           # Main Telegram bot (features 001, 002)
│   └── miniapp/       # Mini App backend server (feature 003)
├── internal/
│   ├── auth/          # Authentication middleware
│   ├── handlers/      # API handlers
│   ├── middleware/    # HTTP middleware
│   └── storage/       # Data storage layer
├── web/               # Mini App frontend (React + TypeScript)
│   ├── src/
│   │   ├── components/
│   │   └── api/
│   └── vite.config.ts
└── specs/             # Feature specifications
```

## Setup

### Prerequisites

- Go 1.21+
- Node.js 20.19+ (for Mini App frontend)
- Telegram Bot Token (from [@BotFather](https://t.me/BotFather))
- Google Gemini API Key (from [ai.google.dev](https://ai.google.dev))

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/telegram-calories-bot.git
   cd telegram-calories-bot
   ```

2. Copy `.env.example` to `.env` and fill in your credentials:
   ```bash
   cp .env.example .env
   # Edit .env with your TELEGRAM_BOT_TOKEN and GEMINI_API_KEY
   ```

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Install frontend dependencies (for Mini App):
   ```bash
   cd web
   npm install
   cd ..
   ```

### Running the Bot (Features 001, 002)

```bash
go run cmd/bot/main.go
```

### Running the Mini App (Feature 003)

**Terminal 1 - Backend Server**:
```bash
go run cmd/miniapp/main.go
```

**Terminal 2 - Frontend Dev Server**:
```bash
cd web
npm run dev
```

The Mini App will be available at `http://localhost:5173` and will proxy API requests to `http://localhost:8080`.

### Accessing the Mini App

The Mini App is designed to be accessed through Telegram only:

1. Open your bot in Telegram
2. Send a command that triggers the Mini App
3. Click on the Mini App button to launch the interface

**Note**: The Mini App authenticates users via `X-Telegram-Init-Data` header. Direct browser access will not work without proper Telegram authentication.

## Development

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt`
- **TypeScript**: ESLint configured with React best practices

### Testing

- **Go**: Run tests with `go test ./...`
- **Mini App**: Manual testing only (demo MVP scope)
- **LLM-Based E2E Testing** (Feature 004): Automated testing using Gemini as LLM judge

#### Running LLM-Based Integration Tests

The LLM-based testing suite validates bot and Mini App behavior using fully automated integration tests with Gemini as an automated judge.

**Key Features:**
- **Fully automated** - no manual user interaction required
- **No real Telegram API calls** - bot handlers invoked directly with mock Context
- **Integration-level testing** - tests real handler logic with mocks/fakes
- **Message preservation validation** - tracks message deletions via mock
- **CRUD operation testing** - httptest for Mini App HTTP handlers

**Prerequisites:**
1. Test food image at `tests/fixtures/food.jpg` (see `tests/fixtures/README.md`)
2. Environment variables configured in `.env.test`
3. **NO** bot deployment or Telegram API access needed (tests run locally)

**Setup:**
```bash
# Copy example config
cp .env.test.example .env.test

# Edit .env.test and fill in:
# - MINIAPP_URL (can be any URL, only used for logging)
# - GEMINI_API_KEY (Gemini API key for LLM judge)
# - TEST_FOOD_IMAGE_PATH (optional, defaults to tests/fixtures/food.jpg)
```

**Run tests:**
```bash
./test-llm.sh
```

Or manually:
```bash
source .env.test
go run cmd/tester/*.go
```

**Test output:**
- **Test report**: `reports/004-test-report.md` (self-contained markdown with embedded evidence)
- **LLM prompts**: Appended to `prompts.md` (all judge prompts verbatim)
- **Exit code**: 0 if all scenarios PASS, 1 if any FAIL

**Test scenarios:**
1. S1: /start command returns welcome message (integration test)
2. S2: /estimate + image upload returns structured estimate (integration test)
3. S3: Re-estimate button does NOT delete previous message (preservation test)
4. S4: Cancel button does NOT delete estimate message (preservation test)
5. S5: Mini App CRUD operations work correctly (httptest integration)

**Testing Approach:**
- **S1-S4**: Direct handler invocation with MockBot capturing sent messages and tracking deletions
- **S5**: httptest.NewServer with real HTTP handlers and fake auth middleware
- **All scenarios**: LLM judge evaluates captured evidence for correctness

## Architecture

### Authentication

User identity is derived **EXCLUSIVELY** from `X-Telegram-Init-Data` header containing Telegram WebApp initData. Backend does NOT accept userID via query parameters or request body.

### Storage

- **Demo MVP**: In-memory storage with `sync.RWMutex`
- **Future**: Designed for migration to Cloudflare D1

### API

RESTful API with CRUD operations:
- `GET /api/logs` - List user's logs
- `POST /api/logs` - Create new log
- `PATCH /api/logs/:id` - Update log
- `DELETE /api/logs/:id` - Delete log

## License

MIT
