# CYTAXI Customer MiniWeb

Conversation-First Mobility Platform — Customer Web Application.

## Purpose

The MiniWeb is the primary web interface for CYTAXI customers. It communicates exclusively through the API Gateway — never directly with Engines.

## Architecture

```
Customer MiniWeb (React/Next.js)
       ↓ (API Gateway: http://localhost:8000/api/v1)
API Gateway (port 8000)
       ↓ (reverse proxy)
Trip / Pricing / Customer / Notification / ... Engines
```

## Tech Stack

- **Framework:** Next.js 14 + React 18
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **State:** React Context + React Query
- **HTTP:** Fetch API

## Project Structure

```
miniweb/
├── src/
│   ├── app/              # App layout + providers
│   ├── components/
│   │   ├── layout/       # Layout, Header, Footer
│   │   └── ui/           # TripCard, BookingForm, ProfileCard, Loading, Error
│   ├── contexts/         # Auth, Trip state
│   ├── hooks/            # useTripRequest, useProfile
│   ├── pages/            # Home, Login, Profile, TripHistory, Notifications, Help
│   ├── services/         # API client (fetches API Gateway)
│   ├── api/              # OpenAPI documentation
│   ├── styles/           # Global styles + Tailwind
│   └── shared/           # Shared types
├── package.json
├── tsconfig.json
└── next.config.js
```

## Pages

| Route | Component | Description |
|-------|-----------|-------------|
| `/` | Home | Landing + booking form + active trip |
| `/login` | Login | Phone-based authentication |
| `/profile` | Profile | User profile + settings links |
| `/trip/history` | TripHistory | Past trips list |
| `/notifications` | Notifications | Notification center |
| `/help` | Help | Help & support |

## Components

| Component | Type | Description |
|-----------|------|-------------|
| `Layout` | Layout | Page wrapper with Header + Footer |
| `Header` | Layout | Navigation bar |
| `Footer` | Layout | Copyright footer |
| `TripCard` | UI | Trip card with status, origin, destination, fare |
| `BookingForm` | UI | Trip request form with pickup, destination, vehicle type |
| `ProfileCard` | UI | Profile summary with avatar, name, phone, email, trip count |
| `Loading` | UI | Loading spinner |
| `ErrorMessage` | UI | Error display with retry |

## State Management

| Context | Purpose |
|---------|---------|
| `AuthContext` | User authentication state |
| `TripContext` | Active trip + history state |

## Development

```bash
npm install
npm run dev        # http://localhost:3000
npm run build      # Production build
npm run typecheck  # TypeScript check
```
