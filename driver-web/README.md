# CYTAXI Driver Web Portal

Driver-facing web application for the CYTAXI platform.

## Purpose

The Driver Web Portal is the primary interface for drivers. It communicates exclusively through the API Gateway.

## Architecture

```
Driver Web Portal (React/Next.js)
       ↓
API Gateway (port 8000) /api/v1
       ↓
Driver / Trip / Matching / Payment / Notification Engines
```

## Tech Stack
React 18 + Next.js 14 + TypeScript + Tailwind CSS + React Query

## Pages

| Route | Description |
|-------|-------------|
| `/login` | Phone-based authentication |
| `/dashboard` | Overview: availability, queue, stats |
| `/trips` | Trip requests queue |
| `/trip/current` | Active trip management |
| `/trip/history` | Past trips |
| `/vehicle` | Vehicle profile |
| `/documents` | Document management |
| `/notifications` | Notification center |
| `/settings` | Driver preferences |
| `/help` | Help & support |

## Components

| Component | Description |
|-----------|-------------|
| `Sidebar` | Navigation sidebar (desktop) |
| `Header` | Top bar with availability toggle |
| `TripCard` | Trip request card |
| `TripQueue` | Trip requests list |
| `AvailabilityToggle` | Online/offline switch |
| `VehicleCard` | Vehicle information card |
| `DocumentCard` | Document status card |

## State Contexts

| Context | Purpose |
|---------|---------|
| `AuthContext` | Driver authentication |
| `TripContext` | Queue, current trip, history |
| `AvailabilityContext` | Online/offline state |

## Development

```bash
npm install
npm run dev       # http://localhost:3001
npm run build
```
