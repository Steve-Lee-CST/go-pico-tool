# ID Generator

- ID generator
- Use a separator to join multiple ID parts together

## Config

- Separator: Separator used to join multiple ID parts
    - If `nil`, the default Separator = `-` will be used
- Modifier: ID parts editor function
    - Parameters:
        - timestamp int64: current time - timestamp (in seconds)
        - microSecond int64: current time - microsecond part
        - randSegment string: random part (first 8 chars of uuid)
    - Return value:
        - parts []string: ID parts
    - Description:
        - This function provides an interface to customize the ID composition.
        - You can add extra information to the ID
        - If `nil`, the default Modifier will be used
        - Modifier should be a **concurrent-safe** function
- Default config
    - Separator: `-`
    - Modifier:
        - timestamp: formatted as readable local time, format: 20060102150405
        - microSecond: no special operation
        - randSegment: no special operation
        - Returns: `[]string{ timestamp, microSecond, randSegment }`

## Main Tool: IDGenerator

The main structure is as follows
```go
type IDGenerator struct {}
func NewIDGenerator(config Config) *IDGenerator {}
func (tool *IDGenerator) Generate() string {}
```

## Helper Functions: None

## Usage Examples

### Using Default Config

```go
import (
    "fmt"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
)

func foo() {
    // Using default config
    gen := id_generator.NewIDGenerator(id_generator.GetDefaultConfig())
    id := gen.Generate()
    fmt.Println("Default ID:", id)
}
```

### Using Custom Config

```go
import (
    "fmt"
    "strings"
    "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
)

func bar() {
    // Custom separator and Modifier
    sep := "_" // Use underscore as separator
    customModifier := func(timestamp int64, microSecond int64, randSegment string) []string {
        // Only keep timestamp and custom prefix
        return []string{
            "CUSTOM", // Add custom prefix
            fmt.Sprintf("%d", timestamp),
            // Remove microSecond
            randSegment,
        }
    }
    cfg := id_generator.Config{
        Separator: &sep,
        Modifier:  customModifier,
    }
    gen := id_generator.NewIDGenerator(cfg)
    id := gen.Generate()
    fmt.Println("Custom ID:", id)
}
```
