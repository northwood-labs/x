# Maps

Utilities for working with maps.

```go
import (
    "log/slog"

    "go.nwlabs.dev/x/maps/v2"
)

func main() {
    details := map[string]string{
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    }

    logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
    logger.Info("This is a message where I'd like to add some details to it.", maps.MapToLogger(details)...)
}
```
