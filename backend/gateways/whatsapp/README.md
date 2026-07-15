# WhatsApp Gateway

WhatsApp integration gateway for the CYTAXI platform.

## Purpose

Provides an abstracted interface for connecting to WhatsApp, managing sessions via QR code, sending and receiving messages. The gateway is provider-agnostic and allows swapping the underlying WhatsApp library without affecting the rest of the system.

## Architecture

```
Conversation Engine
       ↕ (ports.MessageInputPort / ports.ConversationOutputPort)
WhatsApp Gateway
       ↕ (ProviderAdapter interface)
WhatsMeow / WAWebJS / Business API
```

## Components

| Component | Description |
|-----------|-------------|
| `Client` | High-level interface for WhatsApp operations |
| `ProviderAdapter` | Abstraction over WhatsApp provider libraries |
| `SessionManager` | Manages sessions and QR code lifecycle |
| `EventBus` | Publishes gateway events (connected, disconnected, qr_received, message_received) |
| `HealthCheck` | Connection health verification |

## Provider Support

| Provider | Kind | Status |
|----------|------|--------|
| whatsmeow (Go) | `whatsmeow` | Implemented (placeholder) |
| wa-web-js (Node) | `wawebjs` | Planned |
| Business API | `business_api` | Planned |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `WHATSAPP_PROVIDER` | whatsmeow | Provider kind |
| `WHATSAPP_SESSION_ID` | cytaxi-main | Session identifier |
| `WHATSAPP_RECONNECT` | true | Auto-reconnect on disconnect |
| `WHATSAPP_AUTO_LOAD_QR` | true | Auto-generate QR on startup |
| `WHATSAPP_WEBHOOK_URL` | "" | Webhook for incoming messages |

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/gateways/whatsapp/internal/whatsapp"

adapter, _ := whatsapp.NewAdapter(whatsapp.ProviderWhatsMeow)
client := whatsapp.NewClient(whatsapp.ClientConfig{
    SessionID: "cytaxi-main",
}, adapter)

client.Connect(ctx)
qr, _ := client.GetQRCode(ctx)
// Display QR code for scanning
```

## Events

| Event | Description |
|-------|-------------|
| `whatsapp.connected` | Connected to WhatsApp |
| `whatsapp.disconnected` | Disconnected from WhatsApp |
| `whatsapp.qr_received` | New QR code available |
| `whatsapp.message_received` | New message received |
| `whatsapp.message_sent` | Message sent successfully |
| `whatsapp.error` | Error occurred |
