# Files

Utilities for working with files and directories.

```go
import (
    "go.nwlabs.dev/x/files/v2"
)

func main() {
    filename := "myfile.txt"
    substr := "applesauce"

    fmt.Printf(
        "The file '%s' contains the string '%s': %+v\n",
        filename,
        substr,
        files.GrepFile(filename, substr),
    )
}
```
