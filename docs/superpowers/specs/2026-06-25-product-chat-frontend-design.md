# Product Chat Frontend Design

Date: 2026-06-25
Branch: codex/product-chat-frontend

## Goal

Build a clean, product-like web frontend for the existing Go IM backend. The frontend should make the core chat flow usable while respecting the backend APIs that exist today.

The first implementation will cover:

- Register and login.
- Persist authenticated session locally.
- View and edit the current user profile.
- Add chat contacts locally by `user_id`.
- Open a conversation with a saved contact.
- Load conversation history through the existing REST API.
- Send and receive real-time text messages through the existing WebSocket API.
- Show useful connection, loading, empty, and error states.

## Assumptions

- The backend remains unchanged for this frontend pass.
- Contacts and conversation list data are not available from the backend yet.
- Productized contact and conversation surfaces will therefore be backed by browser local state until backend endpoints are added.
- The frontend should follow the backend code as the source of truth when it differs from `docs/v1.0.0-api.md`.
- The backend runs at `http://localhost:8080` by default.
- The WebSocket endpoint is `ws://localhost:8080/api/v1/ws?token=<token>`.

Important backend details from code:

- Login and register responses expose `token`, not `access_token`.
- Message requests use `message_type`, not `msg_type`.
- REST profile update is `PUT /api/v1/users/me`.
- Mark conversation read is `PATCH /api/v1/messages/read` with `other_user_id`.
- WebSocket client events are `message`, `ack`, and `ping`.
- WebSocket message payload uses `receiver_id`, `content`, and `message_type`.

## Recommended Approach

Use a standalone frontend app under `web/` with:

- Vite
- React
- TypeScript
- React Router
- TanStack Query
- Zustand
- Vitest and Testing Library
- lucide-react icons

This is small enough for the current project, provides good type safety, and keeps the frontend independent from the Go backend internals. It also leaves a clean path to add generated API types or richer backend endpoints later.

Alternatives considered:

- Plain HTML/CSS/JavaScript: lowest setup cost, but the chat state and WebSocket behavior would become harder to keep clean.
- Next.js: strong app framework, but server rendering and routing conventions add unnecessary weight for this local IM client.

## Product Scope

The app will feel like a real chat product, but only claim capabilities supported by the backend.

Primary screens:

- Auth screen with login and register modes.
- Main chat shell with:
  - left sidebar for local contacts and conversations;
  - center message timeline and composer;
  - right details panel for current user and selected contact metadata.
- Profile panel for nickname and avatar updates.

Local-only product features:

- Add contact by `user_id`, display name, and optional avatar.
- Persist local contacts and selected conversation in localStorage.
- Show empty states that guide the user to add a contact.

Out of scope for this pass:

- Backend user search.
- Server-backed contact lists.
- Server-backed conversation summaries.
- File or image uploads.
- Push notifications outside the open browser tab.
- Admin or observability screens.

## Frontend Architecture

Proposed structure:

```text
web/
  src/
    app/
      App.tsx
      router.tsx
      providers.tsx
    api/
      client.ts
      auth.ts
      users.ts
      messages.ts
      types.ts
    features/
      auth/
      chat/
      profile/
    shared/
      components/
      hooks/
      lib/
      styles/
    test/
```

Responsibilities:

- `api/` owns REST calls, response envelope parsing, auth headers, and backend field names.
- `features/auth/` owns auth forms and session state integration.
- `features/chat/` owns conversation state, local contacts, WebSocket connection, message timeline, and composer.
- `features/profile/` owns current-user display and update flow.
- `shared/` owns reusable UI pieces and small utilities.

The WebSocket code should live behind a focused hook or service, for example `useChatSocket`, so UI components do not directly manage low-level connection details.

## Data Flow

Authentication:

1. User logs in or registers.
2. API layer parses `{ code, message, data }`.
3. Auth store saves `token` and user information.
4. Protected app routes become available.

Conversation:

1. User adds/selects a local contact by `user_id`.
2. TanStack Query loads `GET /api/v1/messages?user_id=<id>&limit=50&offset=0`.
3. WebSocket connects with the current token.
4. Sending a message writes an optimistic local message and sends:

```json
{
  "event": "message",
  "payload": {
    "receiver_id": "<contact-user-id>",
    "content": "<text>",
    "message_type": "text"
  }
}
```

5. `ack` updates the optimistic message with the backend `message_id`.
6. incoming `message` events are appended to the matching conversation.
7. `unread` events update the sidebar badge when possible.

## UI Direction

The visual direction is quiet, clean, and product-focused:

- Light background with crisp panels and restrained borders.
- Compact sidebar and message layout designed for repeated use.
- Clear connection state near the chat header.
- Soft accent color for the active conversation and primary actions.
- Icons for add contact, send, profile, logout, refresh, and connection state.
- No marketing hero, decorative blobs, or oversized explanatory text.

The first viewport is the actual app, not a landing page.

## Error Handling

The frontend should keep error handling focused and useful:

- Show backend `message` text when available.
- If the token is missing or rejected, clear session and return to auth screen.
- If WebSocket disconnects, show a visible reconnecting/offline state.
- Failed optimistic sends should be marked as failed with a retry affordance if the message is still available locally.
- Empty conversation, no contacts, loading, and failed history states should be explicit.

## Testing and Verification

Implementation should add focused tests before production code for:

- API response envelope parsing and auth header behavior.
- Auth store session persistence and clearing.
- Local contacts store behavior.
- Chat message state transitions for optimistic send, ack, receive, and failure.

After implementation, verify with:

- `npm test` in `web/`.
- `npm run build` in `web/`.
- Browser check of the local app through the dev server.
- Core manual flow: register or login, add contact by `user_id`, load history, send a message, observe WebSocket connection state.

The existing backend baseline is:

- `go test ./...` passes on 2026-06-25 before frontend implementation.

## Success Criteria

- A user can authenticate against the existing backend.
- A user can add a local contact and open a conversation.
- The app can load history for the selected contact.
- The app can send text messages using the backend WebSocket protocol.
- The app can receive WebSocket messages while connected.
- The app remains honest about missing backend-backed contacts and conversation lists.
- The code is organized so future backend endpoints can replace local contact storage without rewriting the UI.
