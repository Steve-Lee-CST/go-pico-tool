# Request ID Tool

- 请求唯一ID工具
- 用于为每个 HTTP 请求生成唯一的 Request ID，并自动注入到请求头和响应头中

## 配置: Config

- HeaderKey: 请求ID在 Header 中的键名
    - 默认为 `X-Request-ID`
- IDGeneratorConfig: ID 生成器配置，详见 [id_generator 使用说明](../../pico_tool/id_generator/_readme.cn.md)
    - 可自定义分隔符、ID组成部分等
- 默认配置
    - HeaderKey: `X-Request-ID`
    - IDGeneratorConfig: 使用 id_generator 的默认配置
- 备注
    - 如果想要自定义ID格式，请声明 IDGeneratorConfig 的时候，通过 Modifier 重写ID的生产方式。可以通过闭包的形式额外带进去更多信息

## 工具本体: RequestIDTool

主体结构如下
```go
// 工具结构体及主要方法
type RequestIDTool struct {}
func NewRequestIDTool(config Config) *RequestIDTool {}
func (tool *RequestIDTool) GenerateRequestID() string {} // 生成一个请求ID
func (tool *RequestIDTool) Middleware() gin.HandlerFunc {} // Gin中间件，自动处理Request ID
func (tool *RequestIDTool) Handler() gin.HandlerFunc {} // Gin处理器，返回新的Request ID
```

## 辅助函数

用于获取、设置RequestID到Header中
- 这里假定在Helper调用之处是可以拿到 config 信息的
- 声明了一个 Helper 的全局变量，可以直接使用此变量调用函数

```go
type helper struct {}
var Helper = helper{}
func (h helper) GetRequestIDFromRequest(c *gin.Context, config Config) (string, bool) // 从 Request 的 Header 中获取 RequestID
func (h helper) GetRequestIDFromResponse(c *gin.Context, config Config) (string, bool) // 从 Response 的 Header 中 获取 RequestID
func (h helper) SetRequestIDToRequest(c *gin.Context, config Config, requestID string) // 向 Request 的 Header 中 写入 RequestID
func (h helper) SetRequestIDToResponse(c *gin.Context, config Config, requestID string) // 向 Response 的 Header 中 写入 RequestID
```

## 使用样例

### 使用默认配置

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/request_id"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    // 使用默认配置
    tool := request_id.NewRequestIDTool(request_id.GetDefaultConfig())
    r.Use(tool.Middleware())
    r.GET("/request_id", tool.Handler())
    r.Run()
}
```

### 使用自定义配置

```go
import (
    "github.com/Steve-Lee-CST/go-pico-tool/gin_tool/request_id"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    sep := "_" // 自定义分隔符
    cfg := request_id.Config{
        HeaderKey: "X-Custom-Request-ID",
        IDGeneratorConfig: id_generator.Config{
            Separator: &sep,
            Modifier:  nil, // 可自定义ID组成部分
        },
    }
    tool := request_id.NewRequestIDTool(cfg)
    r.Use(tool.Middleware())
    r.GET("/request_id", tool.Handler())
    r.Run()
}
```

---

更多 ID 生成器配置说明，详见 [id_generator 使用说明](../../pico_tool/id_generator/_readme.cn.md)

