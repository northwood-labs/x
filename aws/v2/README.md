# AWS Utilities

## `awsutils.GetAWSConfig`

```go
import (
    awsutils "go.nwlabs.dev/x/aws/v2"
)

func main() {
    ctx := context.Background()

    config, err := awsutils.GetAWSConfig(ctx, awsutils.AWSConfigOptions{
        Region: "us-west-2",
        Retries: 5,
        Verbose: false,
    })
    if err != nil {
        fmt.Println(err)
        os.exit(1)
    }
}
```
