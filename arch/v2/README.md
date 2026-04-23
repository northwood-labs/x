# Architecture String

Provides a friendly _U.S. English_ description of the current OS and CPU architecture.

```go
import (
    "fmt"
    "runtime"

    "go.nwlabs.dev/x/arch/v2"
)

func main() {
    fmt.Println(
        arch.GetFriendlyName(runtime.GOOS, runtime.GOARCH),
    )
}
```
