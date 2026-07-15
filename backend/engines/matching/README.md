# Matching Engine

Driver matching and candidate selection for the CYTAXI platform.

## Purpose

The Matching Engine selects the best available driver for a trip by evaluating candidates, ranking them, and applying configurable matching strategies. It does NOT dispatch — that responsibility belongs to the Mobility Decision Engine.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **Matching** — Root aggregate managing the matching lifecycle (7 statuses)
- **CandidateRanking** — Ranked list of driver candidates
- **AssignmentAttempt** — Individual assignment delivery record

### Status Lifecycle
```
pending → searching → evaluating → ranking → selecting → completed
                                                       → failed
                                                       → cancelled
```

### Entities
- **DriverCandidate** — Driver with snapshot, score, rank
- **CandidateSet** — Collection with selection methods
- **AssignmentResult** — Final assignment outcome

### Value Objects
MatchingID, DriverID, TripID, Distance, ETA, MatchingScore, Priority, AvailabilityStatus, CandidateRank

## Strategies
- balanced — weighted combination of distance and score
- nearest — closest driver first
- highest_rated — best driver score first
- priority — driver priority-based

## CQRS

**Commands:** StartMatching, EvaluateCandidates, RankCandidates, SelectCandidate, RetryMatching, CancelMatching

**Queries:** GetMatching, GetCandidates, GetAssignmentHistory, GetRanking, PreviewCandidates

## Events

| Event | Description |
|-------|-------------|
| `matching.started` | Matching process started |
| `matching.candidates_found` | Driver candidates found |
| `matching.candidate_rejected` | Candidate rejected |
| `matching.driver_assigned` | Driver selected |
| `matching.driver_declined` | Driver declined |
| `matching.retried` | Matching retried |
| `matching.cancelled` | Matching cancelled |
| `matching.completed` | Matching completed |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /matching/{id} | Get matching session |
| GET | /matching/{id}/candidates | Get candidates |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MATCHING_PORT` | 8089 | HTTP server port |

## Development

```bash
go run ./cmd/matching
```
