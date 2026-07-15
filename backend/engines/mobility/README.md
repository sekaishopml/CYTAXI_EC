# Mobility Decision Engine

Coordinates mobility decisions for the CYTAXI platform.

## Purpose

The Mobility Decision Engine is responsible for coordinating the dispatch of drivers to trip requests. It evaluates candidates through a pipeline of filters and scorers, then applies a configurable strategy to select the best driver.

## Architecture

```
Conversation Engine → Trip Request
       ↓
DispatcherCoordinator.Dispatch(ctx, DecisionContext)
       ↓
CandidateFinder.FindCandidates()
       ↓
DecisionPipeline.Execute()
  ├── PipelineStep 1: ProximityFilter
  ├── PipelineStep 2: AvailabilityFilter
  ├── PipelineStep 3: VehicleFilter
  └── ...
       ↓
Strategy.Select()  (NearestDriver, HighestRated, BalancedScore)
       ↓
Assignment Decision
```

## Components

| Component | Description |
|-----------|-------------|
| **DispatcherCoordinator** | Entry point; orchestrates the full dispatch flow |
| **DecisionPipeline** | Chain of pipeline steps + final strategy selection |
| **PipelineStep** | Individual filter or scorer in the pipeline |
| **CandidateBuilder** | Builds candidate set from raw drivers with filters and scorers |
| **Strategy** | Selection algorithm (swapable at runtime) |
| **StrategyRegistry** | Registry of available strategies |
| **DecisionContext** | Trip and user context for decision making |

## Strategies

| Strategy | Description |
|----------|-------------|
| **NearestDriver** | Selects the driver closest to the pickup location |
| **HighestRated** | Selects the driver with the highest score |
| **BalancedScore** | Weighted combination of distance and score (configurable) |

## Pipeline

The pipeline executes steps in order. Each step receives the current candidate set and can filter/modify it. After all steps pass, the strategy selects the final driver.

Built-in filter steps:
- `ProximityFilter` — filters by distance radius
- `AvailabilityFilter` — filters by driver availability status

## Events

| Event | Description |
|-------|-------------|
| `mobility.dispatch_started` | Dispatch process started |
| `mobility.dispatch_completed` | Driver assigned successfully |
| `mobility.dispatch_failed` | Dispatch failed |
| `mobility.candidate_found` | Candidates found for trip |
| `mobility.candidate_selected` | Final driver selected |
| `mobility.no_candidates` | No candidates available |
| `mobility.strategy_applied` | Strategy executed |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOBILITY_PORT` | 8084 | HTTP server port |
| `MOBILITY_DEFAULT_STRATEGY` | balanced_score | Default dispatch strategy |
| `MOBILITY_MAX_CANDIDATES` | 20 | Max candidates per dispatch |
| `MOBILITY_DISPATCH_TIMEOUT` | 30 | Dispatch timeout in seconds |

## Domain Model

```
DecisionContext (trip, user, origin, destination, requirements)
       ↓
CandidateSet
  └── Candidate[] (driver, location, vehicle, distance, score)
       ↓
Decision (status, selected_driver, strategy, score, pipeline_summary)
  └── Assignment (status, driver_id, strategy, score)
```

## Development

```bash
go run ./cmd/mobility
```
