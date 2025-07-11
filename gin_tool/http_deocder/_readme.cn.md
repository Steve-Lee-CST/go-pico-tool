# Http Decoder

- HTTP 解码器工具
- 用于解析和捕获 HTTP 请求和响应的详细信息，包括请求头、响应体、表单数据等

## 配置: Config

- HttpRequestKey: HTTP请求在 Gin Context 中的键名
    - 默认为 `http_request`
- HttpResponseKey: HTTP响应在 Gin Context 中的键名
    - 默认为 `http_response`
- HttpResponseWriterKey: HTTP响应写入器在 Gin Context 中的键名
    - 默认为 `http_response_writer`
- 默认配置
    - HttpRequestKey: `http_request`
    - HttpResponseKey: `http_response`
    - HttpResponseWriterKey: `http_response_writer`

## 工具本体: HttpDecoder

主体结构如下
```go
// 工具结构体及主要方法
type HttpDecoder struct {}
func NewHttpDecoder(config Config) *HttpDecoder {}
func (hd *HttpDecoder) Middleware() gin.HandlerFunc {} // Gin中间件，自动解析请求和响应
func (hd *HttpDecoder) Handler() gin.HandlerFunc {} // Gin处理器，返回解析后的HTTP请求信息
```

## 数据模型

### HttpRequest 结构
```go
type HttpRequest struct {
    Method   string `json:"method"`           // HTTP方法
    Protocol string `json:"protocol"`         // HTTP协议版本
    Host     string `json:"host"`             // 主机名
    URL      string `json:"url"`              // 完整URL
    Path     string `json:"path"`             // 请求路径
    FullPath string `json:"full_path"`        // Gin完整路径
    IP       string `json:"ip"`               // 客户端IP

    Header    http.Header    `json:"headers"`     // 请求头
    Cookies   []*http.Cookie `json:"cookies"`    // Cookie
    UserAgent string         `json:"user_agent"`  // 用户代理

    PathParams  gin.Params `json:"path_params,omitempty"`  // 路径参数
    QueryParams url.Values `json:"query_params,omitempty"` // 查询参数

    ContentType string `json:"content_type"` // 内容类型
    RawBody     []byte `json:"raw_body,omitempty"`     // 原始请求体
    JsonBody    map[string]any `json:"json_body,omitempty"` // JSON请求体
    FormValues  url.Values `json:"form_values,omitempty"`   // 表单数据
    FormFiles   []string `json:"form_files,omitempty"`      // 上传文件

    ReceivedTime time.Time `json:"received_time"` // 接收时间
}
```

### HttpResponse 结构
```go
type HttpResponse struct {
    Status      int    `json:"status"`       // 状态码
    StatusText  string `json:"status_text"`  // 状态文本
    Size        int    `json:"size"`         // 响应大小
    ContentType string `json:"content_type"` // 内容类型

    Header  http.Header    `json:"headers,omitempty"`  // 响应头
    Cookies []*http.Cookie `json:"cookies,omitempty"`  // Cookie

    Body     []byte         `json:"body,omitempty"`     // 响应体
    JsonBody map[string]any `json:"json_body,omitempty"` // JSON响应体

    Error      string `json:"error,omitempty"`       // 错误信息
    StackTrace string `json:"stack_trace,omitempty"` // 堆栈跟踪

    ResponseTime time.Time `json:"response_time"` // 响应时间
}
```

## 辅助函数

用于获取、设置HTTP请求和响应信息到Gin Context中
- 声明了一个 Helper 的全局变量，可以直接使用此变量调用函数

```go
type helper struct {}
var Helper = &helper{}
func (h *helper) SetHttpRequest(c *gin.Context, config Config, req *HttpRequest) // 设置HTTP请求到Context
func (h *helper) GetHttpRequest(c *gin.Context, config Config) *HttpRequest // 从Context获取HTTP请求
func (h *helper) SetHttpResponse(c *gin.Context, config Config, resp *HttpResponse) // 设置HTTP响应到Context
func (h *helper) GetHttpResponse(c *gin.Context, config Config) *HttpResponse // 从Context获取HTTP响应
func (h *helper) SetResponseWriter(c *gin.Context, config Config, writer IWrappedResponseWriter) // 设置响应写入器
func (h *helper) GetResponseWriter(c *gin.Context, config Config) IWrappedResponseWriter // 获取响应写入器
func (h *helper) DecodeRequest(c *gin.Context, config Config) *HttpRequest // 解码HTTP请求
func (h *helper) DecodeResponse(c *gin.Context, config Config) *HttpResponse // 解码HTTP响应
```

## 使用样例

### 使用默认配置

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/http_decoder"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    // 使用默认配置
    decoder := http_decoder.NewHttpDecoder(http_decoder.DefaultConfig())
    r.Use(decoder.Middleware())
    r.GET("/decode", decoder.Handler())
    r.Run()
}
```

### 使用自定义配置

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

### 在中间件中获取解析后的请求信息

```go
func customMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取解析后的HTTP请求
        req := http_decoder.Helper.GetHttpRequest(c, http_decoder.DefaultConfig())
        if req != nil {
            fmt.Printf("请求方法: %s\n", req.Method)
            fmt.Printf("请求路径: %s\n", req.Path)
            fmt.Printf("客户端IP: %s\n", req.IP)
            fmt.Printf("请求头: %v\n", req.Header)
        }
        
        c.Next()
        
        // 获取解析后的HTTP响应
        resp := http_decoder.Helper.GetHttpResponse(c, http_decoder.DefaultConfig())
        if resp != nil {
            fmt.Printf("响应状态: %d\n", resp.Status)
            fmt.Printf("响应大小: %d\n", resp.Size)
        }
    }
}
```

### 处理不同类型的请求数据

```go
func handleRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        req := http_decoder.Helper.GetHttpRequest(c, http_decoder.DefaultConfig())
        if req == nil {
            c.JSON(400, gin.H{"error": "无法解析请求"})
            return
        }
        
        // 处理JSON请求
        if len(req.JsonBody) > 0 {
            fmt.Printf("JSON数据: %v\n", req.JsonBody)
        }
        
        // 处理表单数据
        if len(req.FormValues) > 0 {
            fmt.Printf("表单数据: %v\n", req.FormValues)
        }
        
        // 处理上传文件
        if len(req.FormFiles) > 0 {
            fmt.Printf("上传文件: %v\n", req.FormFiles)
        }
        
        c.JSON(200, gin.H{"message": "请求解析成功"})
    }
}
```
