# Debug viewer

Spew is a library that can print the contents of a data structure, much like `PrettyPrint` in Python, or `print_r()` or `var_dump()` in PHP. This implementation provides a specific configuration for Spew that is repeatable.

```go
import (
    "json"
    "log"

    "go.nwlabs.dev/x/debugutils/v2"
)

func main() {
    var results map[string]any

    // Do some stuff to a variable.
    err := json.Unmarshal(data, &results)
    if err != nil {
        log.Fatalf("there was an error: %w", err)
    }

    // Pretty-print the contents of the variable.
    pp := debugutils.GetSpew()
    pp.Dump(results)
}
```
