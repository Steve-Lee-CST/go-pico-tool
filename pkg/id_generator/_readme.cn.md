# ID Generator

- ID 生成器
- 使用 分隔符 将多个ID组成部分连接在一起

## 配置: Config

- Separator: 分隔符, 用于连接多个ID组成部分
    - 如果为 `nil`, 使用默认配置的 Separator = `-`
- Modifier: ID组成部分编辑器,
    - 参数:
        - timestamp int64: 当前时间-时间戳(秒级);
        - microSecond int64: 当前时间-毫秒部分;
        - randSegment string: 随机部分(uuid的前八个字符)
    - 返回值:
        - parts []string: ID组成部分
    - 说明:
        - 此函数提供了一个自定义ID组成形式的接口。
        - 可以额外添加一些信息到ID中
        - 如果为 `nil`, 使用默认配置的 Modifier
        - Modifier 应当为 **并发安全** 的函数
- 默认配置
    - Separator: `-`
    - Modifier:
        - timestamp: 转为可读的本地时间, 格式: 20060102150405
        - microSecond: 无特殊操作
        - randSegment: 无特殊操作
        - 返回: `[]string{ formatted_timestamp, microSecond, randSegment }`

## 工具本体: IDGenerator

主体结构如下
```go
type IDGenerator struct {}
func NewIDGenerator(config Config) *IDGenerator {}
func (tool *IDGenerator) Generate() string {} // 生成一个ID, 并发安全依赖于 Modifier 
```

## 辅助函数: 无

## 使用样例

### 使用默认配置

```go
import (
    "fmt"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
)

func foo() {
    // 使用默认配置
    gen := id_generator.NewIDGenerator(id_generator.GetDefaultConfig())
    id := gen.Generate()
    fmt.Println("默认ID:", id)
}
```

### 使用自定义配置

```go
import (
    "fmt"
    "strings"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
)

func bar() {
    // 自定义分隔符和 Modifier
    sep := "_" // 使用下划线作为分隔符
    customModifier := func(timestamp int64, microSecond int64, randSegment string) []string {
        // 只保留时间戳和自定义前缀
        return []string{
            "CUSTOM", // 增加自定义前缀
            fmt.Sprintf("%d", timestamp),
            // 去掉 microSecond
            randSegment,
        }
    }
    cfg := id_generator.Config{
        Separator: &sep,
        Modifier:  customModifier,
    }
    gen := id_generator.NewIDGenerator(cfg)
    id := gen.Generate()
    fmt.Println("自定义ID:", id)
}
```




