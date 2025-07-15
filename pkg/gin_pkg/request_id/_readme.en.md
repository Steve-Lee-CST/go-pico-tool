# Request ID Tool

- Request unique ID tool
- Used to generate a unique Request ID for each HTTP request and automatically inject it into the request and response headers

## Config

- HeaderKey: The key name of the Request ID in the header
    - Default is `X-Request-ID`
- IDGeneratorConfig: ID generator config, see [id_generator usage](../../pico_tool/id_generator/_readme.en.md)
    - You can customize the separator, ID parts, etc.
- Default config
    - HeaderKey: `X-Request-ID`
    - IDGeneratorConfig: use the default config of id_generator
- Note
    - If you want to customize the ID format, when declaring IDGeneratorConfig, override the ID generation method via Modifier. You can use closures to bring in more information.

## Main Tool: RequestIDTool

The main structure is as follows
```go
// Main tool struct and methods
type RequestIDTool struct {}
func NewRequestIDTool(config Config) *RequestIDTool {}
func (tool *RequestIDTool) GenerateRequestID() string {} // Generate a request ID
func (tool *RequestIDTool) Middleware() gin.HandlerFunc {} // Gin middleware, auto handle Request ID
func (tool *RequestIDTool) Handler() gin.HandlerFunc {} // Gin handler, return a new Request ID
```

## Helper Functions

Used to get/set RequestID in Header
- Assume config info is available where Helper is called
- A global variable Helper is declared, you can use it directly

```go
type helper struct {}
var Helper = helper{}
func (h helper) GetRequestIDFromRequest(c *gin.Context, config Config) (string, bool) // Get RequestID from request header
func (h helper) GetRequestIDFromResponse(c *gin.Context, config Config) (string, bool) // Get RequestID from response header
func (h helper) SetRequestIDToRequest(c *gin.Context, config Config, requestID string) // Set RequestID to request header
func (h helper) SetRequestIDToResponse(c *gin.Context, config Config, requestID string) // Set RequestID to response header
```

## Usage Examples

### Using Default Config

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/request_id"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    // Use default config
    tool := request_id.NewRequestIDTool(request_id.GetDefaultConfig())
    r.Use(tool.Middleware())
    r.GET("/request_id", tool.Handler())
    r.Run()
}
```

### Using Custom Config

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/request_id"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    sep := "_" // Custom separator
    cfg := request_id.Config{
        HeaderKey: "X-Custom-Request-ID",
        IDGeneratorConfig: id_generator.Config{
            Separator: &sep,
            Modifier:  nil, // You can customize the ID parts
        },
    }
    tool := request_id.NewRequestIDTool(cfg)
    r.Use(tool.Middleware())
    r.GET("/request_id", tool.Handler())
    r.Run()
}
```

---

For more ID generator config details, see [id_generator usage](../../pico_tool/id_generator/_readme.en.md)
