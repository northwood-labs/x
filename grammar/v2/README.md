# Spelling and grammar

```go
import (
    "fmt"

    "go.nwlabs.dev/x/grammar/v2"
)

func main() {
    count := 2

    fmt.Printf(
        "I have %d %s.",
        count,
        grammar.Pluralize(count, "apple", "apples")
    )
}
```
