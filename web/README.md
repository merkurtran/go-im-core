# Go IM Core Web

React frontend for the Go IM Core backend.

## Commands

```bash
npm install
npm run dev
npm test
npm run build
```

## Environment

Create `web/.env.local` when the backend is not using the defaults:

```text
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_BASE_URL=ws://localhost:8080/api/v1/ws
```

The first chat version stores contacts in browser localStorage because the current backend does not expose contact or conversation-list endpoints yet.
