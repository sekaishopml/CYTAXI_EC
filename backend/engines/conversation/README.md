# Conversation Engine

Conversation Engine for the CYTAXI platform.

## Purpose

Handles the conversational flow with users via WhatsApp (primary) and other channels. Manages sessions, conversation context, and message pipelines.

## Architecture

```
MessageInputPort
       ↓
MessagePipeline
       ↓
  SessionUseCase →  MessageProcessors
       ↓
  Session / Conversation / Context  (domain entities)
       ↓
  Repository interfaces (Session, Conversation, Message, Context)
```

## Domain Entities

| Entity | Description |
|--------|-------------|
| `Conversation` | Long-lived conversation between user and platform |
| `Session` | Active session with state tracking and expiration |
| `Message` | Individual message within a conversation |
| `ConversationContext` | Key-value context storage per conversation |

## Conversation States

```
new → active → waiting_input → processing → waiting_input (loop)
                 ↓                               ↓
               idle → expired → closed       closed
```

| State | Description |
|-------|-------------|
| `new` | Just created |
| `active` | Session active |
| `waiting_input` | Waiting for user message |
| `processing` | Message being processed |
| `idle` | No recent activity |
| `expired` | Session timed out |
| `closed` | Conversation finished |

## Session Expiration

- Sessions auto-expire after 30 minutes of inactivity.
- `ExpireIdleSessions` can be called periodically as a cron job.
- On expiration, the session and conversation are marked as closed.

## Message Pipeline

When a message arrives:
1. Session is retrieved or created for the phone number
2. State transition validated (→ processing)
3. Message is persisted
4. MessageReceived event emitted
5. Registered processors run
6. State transitions to waiting_input

## Events

| Event | Description |
|-------|-------------|
| `conversation.started` | New conversation created |
| `conversation.closed` | Conversation closed |
| `session.created` | New session created |
| `session.expired` | Session timed out |
| `session.closed` | Session closed |
| `message.received` | New message from user |
| `message.sent` | Message sent to user |
| `conversation.state_changed` | Conversation state transition |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| POST | /messages/incoming | Receive incoming message |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `CONVERSATION_PORT` | 8081 | HTTP server port |
| `APP_ENV` | development | Environment |
| `LOG_LEVEL` | info | Log level |

## Dependencies

- Foundation: `github.com/sekaishopml/cytaxi`

## Development

```bash
go run ./cmd/conversation
```
