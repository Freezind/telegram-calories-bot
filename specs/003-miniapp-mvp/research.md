# Research Findings: Telegram Mini App MVP

**Feature**: 003-miniapp-mvp
**Date**: 2025-12-16
**Purpose**: Resolve technical unknowns identified in Phase 0 planning

## Summary

This document captures research findings for five key technical decisions required to implement the Mini App MVP. All findings prioritize simplicity and standard practices suitable for local development MVP.

---

## R1: Telegram WebApp initData Structure and Parsing

**Unknown**: How to extract and parse Telegram initData in Go backend for user authentication.

### Decision: URL-encoded query string with custom HTTP header

**Rationale**:
- initData is a URL-encoded query string containing `user` (JSON object), `auth_date`, and `hash`
- Standard practice: Send initData from frontend via custom HTTP header `X-Telegram-Init-Data`
- Parse in Go using `url.ParseQuery()` and `json.Unmarshal()` to extract `user.id`

**Implementation**:

```go
package auth

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
)

type TelegramUser struct {
    ID           int64  `json:"id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name,omitempty"`
    Username     string `json:"username,omitempty"`
    LanguageCode string `json:"language_code,omitempty"`
}

func ParseInitData(r *http.Request) (int64, error) {
    initData := r.Header.Get("X-Telegram-Init-Data")
    if initData == "" {
        return 0, errors.New("missing X-Telegram-Init-Data header")
    }

    values, err := url.ParseQuery(initData)
    if err != nil {
        return 0, errors.New("invalid initData format")
    }

    userJSON := values.Get("user")
    if userJSON == "" {
        return 0, errors.New("missing user field in initData")
    }

    var user TelegramUser
    if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
        return 0, errors.New("invalid user JSON")
    }

    if user.ID == 0 {
        return 0, errors.New("invalid user ID")
    }

    return user.ID, nil
}
```

**Frontend Integration** (React):

```typescript
import WebApp from '@twa-dev/sdk'

fetch('http://localhost:8080/api/logs', {
  headers: {
    'X-Telegram-Init-Data': WebApp.initData,
    'Content-Type': 'application/json'
  }
})
```

**Alternatives Considered**:
- Query parameter: Rejected due to URL length limits and logging exposure
- Request body: Rejected because GET requests don't have bodies
- Cookie: Rejected because unnecessary complexity for MVP

**References**:
- Telegram WebApp documentation: https://core.telegram.org/bots/webapps
- initData format: URL-encoded with fields: `user`, `auth_date`, `hash`, `query_id`, etc.

---

## R2: CORS Configuration for localhost Development

**Unknown**: How to configure CORS for Vite (localhost:5173) → Go (localhost:8080) cross-origin requests.

### Decision: Use github.com/rs/cors library with explicit localhost:5173 origin

**Rationale**:
- `rs/cors` is the de-facto standard Go CORS library (battle-tested, well-maintained)
- Handles preflight OPTIONS requests automatically
- Simpler than manual implementation (fewer bugs)
- Zero-config for common use cases

**Implementation**:

```go
package main

import (
    "log"
    "net/http"
    "github.com/rs/cors"
)

func main() {
    mux := http.NewServeMux()

    // Register API routes
    mux.HandleFunc("/api/logs", handleLogs)

    // CORS configuration for Vite dev server
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "X-Telegram-Init-Data"},
        AllowCredentials: true,
        MaxAge:           86400, // 24 hours
        Debug:            true,  // Enable for development
    })

    handler := c.Handler(mux)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

**Key Configuration Options**:
- `AllowedOrigins`: Must specify exact origin (cannot use `*` with credentials)
- `AllowedHeaders`: Include custom `X-Telegram-Init-Data` header
- `AllowCredentials`: Required if frontend uses `credentials: 'include'`
- `Debug: true`: Logs CORS requests during development (remove in production)

**Common Pitfalls**:
1. **Preflight failures**: Ensure OPTIONS method is allowed
2. **Missing headers**: Add all custom headers to `AllowedHeaders`
3. **Wildcard with credentials**: Cannot use `AllowedOrigins: ["*"]` with `AllowCredentials: true`

**Alternatives Considered**:
- Manual CORS implementation: Rejected due to complexity and error-prone nature
- Vite proxy-only (no CORS): Rejected because production deployment would need CORS anyway

**Dependency**:
```bash
go get github.com/rs/cors
```

---

## R3: Vite Proxy Setup for /api/* Requests

**Unknown**: How to configure Vite to proxy API requests from frontend to Go backend.

### Decision: Use Vite server.proxy with /api path rewrite

**Rationale**:
- Vite's built-in proxy (http-proxy) handles all request forwarding automatically
- Preserves method, headers, and body without configuration
- Simplifies frontend API calls (use relative paths like `/api/logs`)

**Implementation**:

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
```

**How It Works**:
- Frontend request: `fetch('/api/logs')` → `http://localhost:5173/api/logs`
- Vite proxy forwards to: `http://localhost:8080/api/logs`
- Backend receives request as if from same origin

**Configuration Options**:
- `target`: Backend server URL (Go at localhost:8080)
- `changeOrigin: true`: Changes `Host` header to match target
- `secure: false`: Allows self-signed certificates (not needed for localhost)
- `rewrite`: Optional path transformation (not needed if backend uses `/api` prefix)

**Example with Path Rewriting** (if backend doesn't use `/api` prefix):

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '')
  }
}
```

This maps `/api/logs` → `/logs` on backend.

**Alternatives Considered**:
- CORS-only (no proxy): Works but exposes backend port in frontend code
- Nginx proxy: Overkill for local development

---

## R4: sync.Map vs mutex-protected map for In-Memory Storage

**Unknown**: Which concurrency pattern performs better for CRUD operations with per-user log arrays?

### Decision: Use map[int64][]Log with sync.RWMutex

**Rationale**:
- **Better for CRUD**: slice operations (append, delete) require read-modify-write, awkward with sync.Map
- **Type safety**: Direct types vs interface{} casting with sync.Map
- **Atomic operations**: RWMutex allows atomic read-modify-write for slice updates
- **Read-heavy optimization**: RWMutex allows unlimited concurrent readers (perfect for frequent ListLogs calls)
- **Simpler code**: More readable and maintainable than sync.Map for this use case

**Performance Characteristics** (for 10-100 users, frequent reads, occasional writes):
- Read operations: ~100-500ns with RLock (negligible)
- Write operations: ~200-1000ns with Lock
- sync.Map only wins when keys are extremely stable and 99%+ reads

**Implementation**:

```go
package storage

import (
    "sync"
    "github.com/freezind/telegram-calories-bot/internal/models"
)

type MemoryStorage struct {
    mu   sync.RWMutex
    logs map[int64][]models.Log  // userID -> logs array
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        logs: make(map[int64][]models.Log),
    }
}

func (m *MemoryStorage) ListLogs(ctx context.Context, userID int64) ([]models.Log, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    userLogs := m.logs[userID]

    // Deep copy to prevent external mutation
    result := make([]models.Log, len(userLogs))
    copy(result, userLogs)

    // Sort newest-first
    sort.Slice(result, func(i, j int) bool {
        return result[i].Timestamp.After(result[j].Timestamp)
    })

    return result, nil
}

func (m *MemoryStorage) CreateLog(ctx context.Context, log *models.Log) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.logs[log.UserID] = append(m.logs[log.UserID], *log)
    return nil
}
```

**Alternatives Considered**:
- `sync.Map`: Rejected due to complexity with slice values and type unsafety
- Channel-based actor pattern: Overkill for MVP
- No synchronization: Rejected due to race conditions

**Migration Path**: If performance becomes an issue at scale, shard the map (multiple maps with separate mutexes) before considering sync.Map.

---

## R5: React Modal Dialog Best Practices

**Unknown**: How to implement modal dialogs for CRUD forms with accessibility and UX best practices?

### Decision: Use @headlessui/react Dialog component with controlled form state

**Rationale**:
- **Built-in accessibility**: Focus trap, ESC key, ARIA attributes handled automatically
- **Unstyled by design**: Full CSS control for custom design
- **Lightweight**: Maintained by Tailwind Labs, zero dependencies
- **Body scroll lock**: Handled automatically
- **Production-ready**: Used by major companies

**Implementation**:

```jsx
import { Dialog, Transition } from '@headlessui/react';
import { Fragment, useState, useEffect } from 'react';

function LogFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData = null,  // null for create, object for edit
  mode = 'create'
}) {
  const [formData, setFormData] = useState({
    foodItems: [],
    calories: 0,
    confidence: 'medium'
  });

  useEffect(() => {
    if (isOpen) {
      setFormData(initialData || { foodItems: [], calories: 0, confidence: 'medium' });
    }
  }, [isOpen, initialData]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    await onSubmit(formData);
    onClose();
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
        >
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4">
            <Dialog.Panel className="w-full max-w-md bg-white rounded p-6">
              <Dialog.Title>{mode === 'create' ? 'Add Log' : 'Edit Log'}</Dialog.Title>

              <form onSubmit={handleSubmit}>
                {/* Form fields here */}
                <button type="submit">Save</button>
                <button type="button" onClick={onClose}>Cancel</button>
              </form>
            </Dialog.Panel>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
```

**Form State Management**:
- Use `useState` for controlled inputs (simple, sufficient for MVP)
- Alternatives (React Hook Form, Formik) only needed for complex validation

**Accessibility Features** (automatic with Headless UI):
- Focus trap: Tab cycles within modal
- ESC key: Closes modal
- Click outside: Configurable via `onClose`
- ARIA attributes: `role="dialog"`, `aria-modal="true"`
- Focus management: Returns focus to trigger element on close
- Body scroll lock: Prevents page scrolling when modal open

**Alternatives Considered**:
- react-modal: Older, more configuration needed
- Custom implementation: Too complex for accessibility requirements
- Material-UI Dialog: Heavy dependency for MVP

**Dependency**:
```bash
npm install @headlessui/react
```

**Styling**: Works with Tailwind CSS or plain CSS. Unstyled by default for maximum flexibility.

---

## Technology Stack Summary

Based on research findings, the final technology choices are:

### Backend (Go)
- **HTTP Framework**: net/http (standard library)
- **CORS**: github.com/rs/cors
- **UUID**: github.com/google/uuid
- **Storage**: map[int64][]Log with sync.RWMutex

### Frontend (React)
- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **Telegram SDK**: @twa-dev/sdk
- **Modal Dialogs**: @headlessui/react
- **HTTP Client**: fetch API (native)

### Development Setup
- **Backend Port**: localhost:8080
- **Frontend Port**: localhost:5173
- **Proxy**: Vite proxy for /api/* → backend
- **CORS**: Enabled for localhost:5173 origin

---

## Open Questions (None Remaining)

All unknowns from Technical Context have been resolved. Ready to proceed to Phase 1: Design & Contracts.

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-12-16 | Initial research findings for 5 technical unknowns |
