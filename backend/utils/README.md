# Utils

General-purpose utilities for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/utils"

id := utils.NewID()       // 32-char hex ID
sid := utils.NewShortID() // 12-char hex ID

val := utils.Coalesce(input, fallback) // First non-zero value
ptr := utils.Ptr(value)                // Pointer to value
```
