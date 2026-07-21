# CYTAXI

Conversation-First Mobility Platform.

## Architecture

- **Primary Channel:** WhatsApp
- **Secondary Channels:** Miniweb, Driver Web App, Admin Dashboard
- **Backend:** Go
- **Architecture:** Event-Driven Microservices, DDD, Clean Architecture, CQRS

## Repository Structure

```
CYTAXI_EC/
├── .editorconfig
├── .env.example
├── .golangci.yml
├── .github/workflows/ci.yml
├── AGENTS.md
├── CHANGELOG.md
├── Dockerfile.engine
├── EXECUTIVE_SUMMARY.md
├── Makefile
├── RELEASE_NOTES.md
├── VERSION
├── docker-compose.dev.yml / .beta.yml / .prod.yml / .quick.yml
├── go.work
├── opencode.jsonc
├── package.json
│
├── backend/                         # Go backend (monorepo)
│   ├── go.mod
│   ├── auth/                        # JWT auth
│   ├── config/                      # Config loader
│   ├── containers/                  # DI container
│   ├── engine-template/             # Engine scaffold (Clean Arch + DDD)
│   ├── engines/                     # Microservices
│   │   ├── admin/                   #   Admin panel engine
│   │   ├── analytics/               #   Analytics engine
│   │   ├── conversation/            #   WhatsApp/AI conversation engine
│   │   ├── customer/                #   Customer engine
│   │   ├── driver/                  #   Driver engine
│   │   ├── geospatial/              #   Geospatial engine (OSRM)
│   │   ├── matching/                #   Driver-rider matching
│   │   ├── mobility/                #   Mobility services
│   │   ├── notification/            #   Push notifications
│   │   ├── payment/                 #   Payment processing
│   │   ├── policy/                  #   Business rules engine
│   │   ├── pricing/                 #   Dynamic pricing
│   │   ├── trip/                    #   Trip lifecycle
│   │   └── trust/                   #   Trust & safety
│   ├── errors/                      # Error handling
│   ├── events/                      # Event bus
│   ├── flow/                        # Journey flow engine
│   ├── foundation/                  # Shared foundation
│   ├── gateway/                     # API gateway (chi router)
│   ├── gateways/whatsapp/           # WhatsApp provider
│   ├── http/                        # HTTP server/client
│   ├── integration/                 # Event bus, saga, outbox, contracts
│   ├── llm/                         # LLM orchestrator
│   ├── logger/                      # Structured logging
│   ├── middleware/                  # HTTP middleware (CORS, correlation, etc.)
│   ├── mobile/                      # Mobile push manager
│   ├── observability/              # Health, metrics
│   ├── telemetry/                   # Distributed tracing
│   ├── testing/                     # Test utilities
│   ├── utils/                       # General utilities
│   ├── validation/                  # Validation
│   └── version/                     # Version info
│
├── frontend/                        # Frontend apps
│   ├── travel/                      #   Travel — Next.js 14 (App Router)
│   │   ├── src/
│   │   │   ├── app/                 #   App Router pages
│   │   │   ├── components/          #   States, MapPreview, BottomSheet, etc.
│   │   │   ├── services/            #   API, tracking, offline-queue
│   │   │   ├── hooks/useJourneyEngine.ts
│   │   │   └── ...
│   │   ├── tailwind.config.js
│   │   └── vitest.config.ts
│   ├── driver-web/                  #   Driver web app — Next.js 14 (App Router)
│   │   └── src/
│   │       ├── pages/               #   dashboard, trips, earnings, profile, etc.
│   │       ├── components/
│   │       ├── contexts/            #   auth, availability, trip
│   │       ├── hooks/
│   │       └── services/
│   └── dashboard/                   #   Admin dashboard — Next.js (Pages Router)
│       └── src/
│           ├── components/Layout.tsx
│           ├── pages/               #   billing, system, tenants, index
│           └── styles/
│
├── packages/                        # Shared TypeScript packages
│   ├── api-client/                  #   API client SDK
│   ├── events/                      #   Event bus (typed)
│   ├── ride-machine/                #   Trip state machine
│   ├── map-engine/                  #   Map abstraction layer
│   ├── design-tokens/               #   Colors, typography, spacing, shadows
│   ├── ui/                          #   UI components (Button, Card, Modal, etc.)
│   ├── i18n/                        #   Internationalization
│   ├── offline/                     #   Offline queue
│   ├── realtime/                    #   Real-time client
│   ├── multi-tenant/                #   Multi-tenant utilities
│   ├── billing/                     #   Billing helpers
│   ├── webhooks/                    #   Webhook utilities
│   ├── security/                    #   Security helpers
│   ├── observability/               #   Frontend observability
│   ├── trust-score/                 #   Trust score client
│   ├── llm-conversation/            #   LLM conversation client
│   ├── analytics/                   #   Analytics client
│   ├── ai/                          #   AI service client
│   ├── sounds/                      #   Sound effects
│   └── fonts/                       #   Geist, Inter CSS
│
├── docs/                            # Architecture docs, ADRs, operations
│   ├── adr/                         # Architecture Decision Records (13)
│   ├── BACKEND_ARCHITECTURE.md
│   ├── API_GUIDE.md
│   ├── DEPLOYMENT.md
│   ├── OPERATIONS_GUIDE.md
│   ├── SCALING_ARCHITECTURE.md
│   └── ...
│
├── knowledge/                       # Logseq knowledge graph
│   ├── pages/                       # architecture, domain, decisions, sprints…
│   ├── journals/
│   └── logseq/config.edn
│
├── infra/                           # Infrastructure
│   ├── grafana/dashboards/
│   ├── prometheus/prometheus.yml
│   ├── security/                    # fail2ban, SSL, production checklist
│   └── backup/
│
├── nginx/                           # Nginx configs
│   ├── nginx.conf
│   ├── static.conf
│   └── swagger/
│
├── deploy/                          # Deployment runbooks
├── scripts/                         # CI/CD, smoke tests, verification
├── screenshots/                     # UI screenshots
└── public/stitch/                   # Static HTML mockups
```

## Development

```bash
make lint    # Run linters
make test    # Run tests
make build   # Build project
```

## Principles

- Documentation before implementation
- Architecture before optimization
- Business before technology
- Consistency over speed
- Replaceability
- Long-term thinking

## Source of Truth

The [CYTAXI-BLUEPRINT](https://github.com/sekaishopml/cydigital-blueprint) repository is the single source of truth for architecture, business rules, and engineering standards.
