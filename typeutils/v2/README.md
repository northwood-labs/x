# types

## RFC3339Time and UnixTime

In structs that represent JSON documents, instead of specifying `string` for RFC 3339-formatted dates, or `int` for Unix-formatted dates, use `RFC3339Time` and `UnixTime` instead. These will unmarshal those values into `time.Time` objects, and will re-marshal them into the appropriate format in JSON form.

```go
import (
    "go.nwlabs.dev/x/typeutils/v2"
)

type Response struct {
    Results       []string              `json:"results,omitempty"`
    DownloadCount int64                 `json:"downloads"`
    CreatedAt     typeutils.RFC3339Time `json:"createdAt,omitempty"`
}

func (response *Response) FromJSON(data []byte) error {
    return json.Unmarshal(data, response)
}

func (response *Response) AsJSON() ([]byte, error) {
    return json.Marshal(response)
}

func main() {
    jsonData := `{
        results: ["a", "b", "c"],
        downloads: 123,
        createdAt: "2023-01-01T00:00:00Z"
    }`

    response := new(Response)
    _ = response.FromJSON(jsonData)

    // response.CreatedAt.Time is a time.Time object.

    data, _ := response.AsJSON()

    // The date serializes back to RFC 3339 format.
    //=> data = `{results: ["a", "b", "c"], downloads: 123, createdAt: "2023-01-01T00:00:00Z"}`
}
```
