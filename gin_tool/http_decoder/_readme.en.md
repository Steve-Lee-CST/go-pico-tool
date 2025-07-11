# Http Decoder

- HTTP decoder tool
- Used to parse and capture detailed information of HTTP requests and responses, including request headers, response body, form data, etc.

## Configuration: Config

- HttpRequestKey: Key name for HTTP request in Gin Context
    - Default: `http_request`
- HttpResponseKey: Key name for HTTP response in Gin Context
    - Default: `http_response`
- HttpResponseWriterKey: Key name for HTTP response writer in Gin Context
    - Default: `http_response_writer`
- Default configuration
    - HttpRequestKey: `http_request`
    - HttpResponseKey: `http_response`
    - HttpResponseWriterKey: `http_response_writer`

## Tool Entity: HttpDecoder

Main structure as follows
```go
// Tool struct and main methods
type HttpDecoder struct {}
func NewHttpDecoder(config Config) *HttpDecoder {}
func (hd *HttpDecoder) Middleware() gin.HandlerFunc {} // Gin middleware, automatically parse requests and responses
func (hd *HttpDecoder) Handler() gin.HandlerFunc {} // Gin handler, return parsed HTTP request information
```

## Data Models

### HttpRequest Structure
```go
type HttpRequest struct {
    Method   string `json:"method"`           // HTTP method
    Protocol string `json:"protocol"`         // HTTP protocol version
    Host     string `json:"host"`             // Host name
    URL      string `json:"url"`              // Complete URL
    Path     string `json:"path"`             // Request path
    FullPath string `json:"full_path"`        // Gin full path
    IP       string `json:"ip"`               // Client IP

    Header    http.Header    `json:"headers"`     // Request headers
    Cookies   []*http.Cookie `json:"cookies"`    // Cookies
    UserAgent string         `json:"user_agent"`  // User agent

    PathParams  gin.Params `json:"path_params,omitempty"`  // Path parameters
    QueryParams url.Values `json:"query_params,omitempty"` // Query parameters

    ContentType string `json:"content_type"` // Content type
    RawBody     []byte `json:"raw_body,omitempty"`     // Raw request body
    JsonBody    map[string]any `json:"json_body,omitempty"` // JSON request body
    FormValues  url.Values `json:"form_values,omitempty"`   // Form data
    FormFiles   []string `json:"form_files,omitempty"`      // Uploaded files

    ReceivedTime time.Time `json:"received_time"` // Received time
}
```

### HttpResponse Structure
```go
type HttpResponse struct {
    Status      int    `json:"status"`       // Status code
    StatusText  string `json:"status_text"`  // Status text
    Size        int    `json:"size"`         // Response size
    ContentType string `json:"content_type"` // Content type

    Header  http.Header    `json:"headers,omitempty"`  // Response headers
    Cookies []*http.Cookie `json:"cookies,omitempty"`  // Cookies

    Body     []byte         `json:"body,omitempty"`     // Response body
    JsonBody map[string]any `json:"json_body,omitempty"` // JSON response body

    Error      string `json:"error,omitempty"`       // Error message
    StackTrace string `json:"stack_trace,omitempty"` // Stack trace

    ResponseTime time.Time `json:"response_time"` // Response time
}
```

## Helper Functions

Used to get and set HTTP request and response information to Gin Context
- Declared a global Helper variable, can directly use this variable to call functions

```go
type helper struct {}
var Helper = &helper{}
func (h *helper) SetHttpRequest(c *gin.Context, config Config, req *HttpRequest) // Set HTTP request to Context
func (h *helper) GetHttpRequest(c *gin.Context, config Config) *HttpRequest // Get HTTP request from Context
func (h *helper) SetHttpResponse(c *gin.Context, config Config, resp *HttpResponse) // Set HTTP response to Context
func (h *helper) GetHttpResponse(c *gin.Context, config Config) *HttpResponse // Get HTTP response from Context
func (h *helper) SetResponseWriter(c *gin.Context, config Config, writer IWrappedResponseWriter) // Set response writer
func (h *helper) GetResponseWriter(c *gin.Context, config Config) IWrappedResponseWriter // Get response writer
func (h *helper) DecodeRequest(c *gin.Context, config Config) *HttpRequest // Decode HTTP request
func (h *helper) DecodeResponse(c *gin.Context, config Config) *HttpResponse // Decode HTTP response
```

## Usage Examples

### Using Default Configuration

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/http_decoder"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    // Use default configuration
    decoder := http_decoder.NewHttpDecoder(http_decoder.DefaultConfig())
    r.Use(decoder.Middleware())
    r.GET("/decode", decoder.Handler())
    r.Run()
}
```

### Using Custom Configuration

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/http_decoder"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    cfg := http_decoder.Config{
        HttpRequestKey:        "custom_request",
        HttpResponseKey:       "custom_response",
        HttpResponseWriterKey: "custom_writer",
    }
    decoder := http_decoder.NewHttpDecoder(cfg)
    r.Use(decoder.Middleware())
    r.GET("/decode", decoder.Handler())
    r.Run()
}
```

### Getting Parsed Request Information in Middleware

```go
func customMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get parsed HTTP request
        req := http_decoder.Helper.GetHttpRequest(c, http_decoder.DefaultConfig())
        if req != nil {
            fmt.Printf("Request method: %s\n", req.Method)
            fmt.Printf("Request path: %s\n", req.Path)
            fmt.Printf("Client IP: %s\n", req.IP)
            fmt.Printf("Request headers: %v\n", req.Header)
        }
        
        c.Next()
        
        // Get parsed HTTP response
        resp := http_decoder.Helper.GetHttpResponse(c, http_decoder.DefaultConfig())
        if resp != nil {
            fmt.Printf("Response status: %d\n", resp.Status)
            fmt.Printf("Response size: %d\n", resp.Size)
        }
    }
}
```

### Handling Different Types of Request Data

```go
func handleRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        req := http_decoder.Helper.GetHttpRequest(c, http_decoder.DefaultConfig())
        if req == nil {
            c.JSON(400, gin.H{"error": "Unable to parse request"})
            return
        }
        
        // Handle JSON requests
        if len(req.JsonBody) > 0 {
            fmt.Printf("JSON data: %v\n", req.JsonBody)
        }
        
        // Handle form data
        if len(req.FormValues) > 0 {
            fmt.Printf("Form data: %v\n", req.FormValues)
        }
        
        // Handle uploaded files
        if len(req.FormFiles) > 0 {
            fmt.Printf("Uploaded files: %v\n", req.FormFiles)
        }
        
        c.JSON(200, gin.H{"message": "Request parsed successfully"})
    }
}
```
