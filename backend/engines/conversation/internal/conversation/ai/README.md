# AI Orchestrator

AI orchestration layer for the CYTAXI Conversation Engine.

## Purpose

Decides when to use an LLM model and when to resolve a request via deterministic rules. Provides provider abstraction so that Qwen, DeepSeek, GPT, Claude, or any other model can be swapped without affecting the platform.

## Architecture

```
MessagePipeline
       ↓
  Orchestrator.Process()
       ↓
  Router.Route() → Classifier.Classify() → Intent
       ↓
  PolicyEvaluator.ShouldUseLLM()
       ↓
  (yes) → ProviderRegistry.Select() → Provider.Complete()
       ↓                                            ↓
  (no)  → Deterministic path              FallbackManager.Handle() (on error)
```

## Components

| Component | Description |
|-----------|-------------|
| `Orchestrator` | Entry point; makes the LLM vs deterministic decision |
| `Router` | Routes input through classifiers to determine intent |
| `Classifier` | Interface for intent classification (LLM-based or rule-based) |
| `Provider` | Interface for LLM providers (Qwen, DeepSeek, GPT, Claude) |
| `ProviderRegistry` | Registry and selection of available providers |
| `PromptBuilder` | Builds system and user prompts from intent and context |
| `ContextBuilder` | Builds conversation context for LLM requests |
| `PolicyEvaluator` | Determines if LLM should be used for a given intent |
| `FallbackManager` | Handles LLM failures with retry/fallback/static strategies |
| `MetricsCollector` | Tracks decisions, provider errors, latency |

## Provider Adapter

The `Provider` interface allows swapping LLM providers:

```go
type Provider interface {
    Name() string
    Kind() ProviderKind
    Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
}
```

Implementations:
- `NewMockProvider(kind)` — mock for development/testing
- Future: `NewQwenProvider()`, `NewDeepSeekProvider()`, `NewGPTProvider()`

## Policy Rules

Defined per intent: whether to use LLM, minimum confidence threshold, priority.

| Intent | Use LLM | Min Confidence |
|--------|---------|----------------|
| greeting | No | — |
| trip_request | Yes | 0.7 |
| trip_status | No | — |
| support | Yes | 0.6 |
| cancel | No | — |
| unknown | Yes | 0.4 |

## Fallback Strategies

1. Retry (up to 2 times)
2. Static response ("Lo siento, no puedo procesar...")
3. Fallback provider (switch to alternative LLM)

## Observability

- `MetricsCollector` tracks: total requests, LLM vs deterministic split, fallback count, provider errors
- `AILogger` interface for tracing AI decisions
- Snapshot method for health/metrics endpoints
