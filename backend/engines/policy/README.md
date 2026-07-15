# Policy Engine

Centralized business rule evaluation for the CYTAXI platform.

## Purpose

The Policy Engine centralizes all business decision rules (pricing, dispatch, zone restrictions, etc.) so they can be defined, versioned, and evaluated independently of any specific Engine. This keeps business logic decoupled, auditable, and hot-reloadable.

## Architecture

```
Engine (Conversation, Geospatial, etc.)
       ↓
PolicyEngine.Evaluate(ctx, DecisionContext, domains...)
       ↓
PolicyRegistry.GetPolicies(domains...)
       ↓
RuleEvaluator.EvaluatePolicy(policy, decisionCtx)
       ↓
ConditionEvaluator.Evaluate(condition, decisionCtx)
       ↓
Decision (matched actions)
```

## Components

| Component | Description |
|-----------|-------------|
| **PolicyEngine** | Entry point; evaluates all policies for given domains |
| **PolicyRegistry** | In-memory registry with domain-based lookup |
| **RuleEvaluator** | Evaluates all rules in a policy against a decision context |
| **ConditionEvaluator** | Evaluates individual conditions with typed operators |
| **PolicyLoader** | Loads policies from file, memory, or other sources |
| **PolicyRepository** | Interface for persistent policy storage |
| **VersionRepository** | Interface for version history tracking |

## Domain Model

```
Policy (id, name, domain, version, priority, enabled)
  └── Rules []
        ├── Conditions []  (field, operator, value)
        └── Actions  []    (type, params)

DecisionContext (user, role, location, trip, pricing, time, custom)
  └── EvaluationResult
        └── Decisions []   (policy, rule, actions, matched)
```

## Condition Operators

| Operator | Description |
|----------|-------------|
| `equals` | Exact match |
| `not_equals` | Not equal |
| `greater_than` | Numeric > |
| `less_than` | Numeric < |
| `in` | Value in slice |
| `not_in` | Value not in slice |
| `contains` | String contains |
| `starts_with` | String prefix |
| `is_true` | Boolean true |
| `is_false` | Boolean false |

## Policy Versioning

- Semantic versioning (Major.Minor.Patch)
- Version lifecycle: draft → active → deprecated → archived
- `VersionRepository` for persistent version tracking
- Policies can be hot-reloaded without restart

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `POLICY_PORT` | 8083 | HTTP server port |
| `POLICIES_DIR` | ./policies | Directory for JSON policy files |
| `POLICY_AUTO_RELOAD` | true | Auto-reload on file changes |
| `POLICY_RELOAD_SEC` | 60 | Reload interval in seconds |

## Events

| Event | Description |
|-------|-------------|
| `policy.created` | New policy registered |
| `policy.updated` | Policy updated |
| `policy.deleted` | Policy removed |
| `policy.evaluated` | Policy evaluated (matched/not) |
| `policy.rule_matched` | A rule matched during evaluation |
| `policy.reloaded` | All policies reloaded from source |

## Development

```bash
go run ./cmd/policy
```
