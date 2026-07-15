# Config

Configuration management for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/config"

cfg, err := config.Load()
if err != nil {
    // handle error
}
fmt.Println(cfg.App.Port)
```

## How it works

- Reads environment variables with sensible defaults.
- Supports `APP_ENV`, `DATABASE_URL`, `LOG_LEVEL`, etc.
- Returns a typed `*config.Config` struct.
- `config.Load()` uses `os.Getenv` — zero external dependencies.

## Environment variables

All variables are documented in `.env.example` at the repository root.
