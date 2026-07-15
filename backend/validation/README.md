# Validation

Struct validation for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/validation"

type CreateDriverRequest struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
}

var req CreateDriverRequest
if err := validation.ValidateStruct(req); err != nil {
    // validation.Errors with field-level messages
}
```

## How it works

- Uses struct tags (`validate:"required,email"`) for rule declarations.
- Built-in rules: `required`, `email`.
- Extensible via `RegisterRule(name string, fn RuleFunc)`.
- Returns `*validation.Errors` with per-field messages.
