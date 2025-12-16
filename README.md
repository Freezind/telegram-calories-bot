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
