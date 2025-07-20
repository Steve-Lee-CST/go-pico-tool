# Task DAG Flow

- 任务有向无环图(DAG)流执行器
- 用于管理和执行具有依赖关系的任务集合，支持并发执行和超时控制
- 基于任务的输入输出类型自动构建依赖关系图，实现任务的有序并发执行

## 核心概念

### ICollection 接口
- 数据集合接口，定义了任务流中数据的管理方式
- `InputTypes()`: 返回任务流开始时可用的数据类型
- `TargetTypes()`: 返回任务流期望产生的目标数据类型

### ITask 接口
- 任务接口，类似于函数但通过数据集合传递参数和返回值
- `Name()`: 返回任务名称
- `InputTypes()`: 返回任务依赖的输入数据类型
- `OutputType()`: 返回任务产生的输出数据类型
- `Timeout()`: 返回任务超时时间
- `Execute(ctx, collection)`: 执行任务逻辑

## 主要组件

### Factory[CT ICollection]
工厂类，用于注册任务和创建任务流
```go
type Factory[CT ICollection] struct {}
func NewFactory[CT ICollection]() *Factory[CT] {}
func (f *Factory[CT]) RegisterTask(createFunc TaskCreateFunc[CT]) error {} // 注册任务
func (f *Factory[CT]) CreateGraph() {} // 创建依赖关系图
func (f *Factory[CT]) CreateTaskDagflow(collection CT) (*TaskDagflow[CT], error) {} // 创建任务流
```

### TaskDagflow[CT ICollection]
任务流执行器，管理任务的并发执行，一般从工厂创建
```go
type TaskDagflow[CT ICollection] struct {}
func NewTaskDagflow[CT ICollection](metas []*taskMeta[CT], collection CT) (*TaskDagflow[CT], error) {}
func (td *TaskDagflow[CT]) Execute(ctx context.Context, timeout time.Duration) error {} // 执行任务流
```

## 辅助函数

### 自动类型推导
```go
func AutoInputTypes[CT ICollection]() []reflect.Type {} // 自动推导输入类型(基于Get方法)
func AutoOutputType[CT ICollection]() reflect.Type {} // 自动推导输出类型(基于Set方法)
```

### 任务创建
```go
type TaskCreateFunc[CT ICollection] func() (ITask[CT], error) // 任务创建函数类型
func CreateTask[CT ICollection](createFunc TaskCreateFunc[CT]) (ITask[CT], error) {} // 创建任务实例
```

## 特性说明

### 依赖关系自动解析
- 基于任务的 `InputTypes()` 和 `OutputType()` 自动构建DAG
- 支持复杂的多层依赖关系
- 自动检测循环依赖和不满足的依赖

### 并发执行优化
- 无依赖关系的任务可并发执行
- 动态调度，任务完成后立即触发后续可执行任务
- 支持任务级别和流级别的超时控制

### 类型安全
- 使用泛型确保编译时类型安全
- 反射机制进行运行时类型匹配
- 自动验证输入输出类型的一致性

## 使用样例

### 基本使用流程

```go
import (
    "context"
    "time"
    "github.com/Steve-Lee-CST/go-pico-tool/pkg/task_dagflow"
)

// 1. 定义数据集合
type DataCollection struct {
    goods []Goods
    shops []Shop
    result GoodsInShops
}

func (c *DataCollection) InputTypes() []reflect.Type {
    return []reflect.Type{nil} // 无初始输入
}

func (c *DataCollection) TargetTypes() []reflect.Type {
    return []reflect.Type{reflect.TypeOf(GoodsInShops{})}
}

// 2. 实现具体任务
type GetGoodsTask struct {
    name string
    timeout time.Duration
}

func (t *GetGoodsTask) Name() string { return t.name }
func (t *GetGoodsTask) InputTypes() []reflect.Type { return []reflect.Type{nil} }
func (t *GetGoodsTask) OutputType() reflect.Type { return reflect.TypeOf([]Goods{}) }
func (t *GetGoodsTask) Timeout() time.Duration { return t.timeout }
func (t *GetGoodsTask) Execute(ctx context.Context, collection *DataCollection) error {
    // 获取商品数据的逻辑
    collection.goods = fetchGoods()
    return nil
}

// 3. 使用工厂创建任务流
func main() {
    factory := task_dagflow.NewFactory[*DataCollection]()
    
    // 注册任务
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &GetGoodsTask{name: "GetGoods", timeout: 500 * time.Millisecond}, nil
    })
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &GetShopsTask{name: "GetShops", timeout: 500 * time.Millisecond}, nil
    })
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &ProcessTask{name: "Process", timeout: 1000 * time.Millisecond}, nil
    })
    
    // 创建依赖图
    factory.CreateGraph()
    
    // 创建数据集合和任务流
    collection := &DataCollection{}
    taskDagflow, err := factory.CreateTaskDagflow(collection)
    if err != nil {
        panic(err)
    }
    
    // 执行任务流
    ctx := context.Background()
    if err := taskDagflow.Execute(ctx, 2*time.Second); err != nil {
        panic(err)
    }
    
    // 使用结果
    fmt.Printf("执行完成，结果: %+v\n", collection.result)
}
```

### 使用自动类型推导

```go
// 定义接口用于自动类型推导
type IGoodsProvider interface {
    task_dagflow.ICollection
    GetGoods() []Goods    // 自动识别为输入类型
    SetResult(result GoodsInShops) // 自动识别为输出类型
}

type ProcessTask struct {
    name string
    timeout time.Duration
}

func (t *ProcessTask) InputTypes() []reflect.Type {
    return task_dagflow.AutoInputTypes[IGoodsProvider]() // 自动推导输入类型
}

func (t *ProcessTask) OutputType() reflect.Type {
    return task_dagflow.AutoOutputType[IGoodsProvider]() // 自动推导输出类型
}
```

### 复杂依赖关系示例

在 `demo_test.go` 中展示了一个完整的示例：
- `GetGoodsTask`: 获取商品数据 (无依赖)
- `GetShopsTask`: 获取商店数据 (无依赖)  
- `GoodsInShopsTask`: 处理商品与商店关系 (依赖前两个任务)

这个示例展示了：
- 前两个任务可以并发执行
- 第三个任务等待前两个任务完成后执行
- 支持任务级别和流级别的超时控制
- 自动处理复杂的数据依赖关系

## 注意事项

- 关于任务：
  - 任务的输出类型应当是Factory级唯一的：
    - 可以理解为，输出的数据类型可以作为某一个任务的唯一标识
  - 任务的输出类型不应当是其自身的输入类型之一：不能自成环
  - 任务之间不应该直接通信，只通过数据集合传递数据
  - 确保任务超时时间设置合理
- 关于数据集合
  - 某一数据类型，只能被其所对应的任务写入，其余任务只能读取；在此基础上，数据集合是并发安全的
  - 建议使用 getter 和 setter 方法来访问数据集合中的数据
- 关于Factory 和 Dagflow
  - 当前只检测不可达，不区分不可达原因：可能是循环依赖、可能是某任务输入缺失
  - 如果所求数据存在不可达类型，返回错误。
