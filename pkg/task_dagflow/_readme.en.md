# Task DAG Flow

- Task Directed Acyclic Graph (DAG) Flow Executor
- Used to manage and execute task sets with dependencies, supporting concurrent execution and timeout control
- Automatically builds dependency graphs based on task input/output types to achieve ordered concurrent execution

## Core Concepts

### ICollection Interface
- Data collection interface that defines data management in task flows
- `InputTypes()`: Returns data types available at the start of the task flow
- `TargetTypes()`: Returns target data types that the task flow expects to produce

### ITask Interface
- Task interface, similar to a function but passes parameters and return values through data collections
- `Name()`: Returns task name
- `InputTypes()`: Returns input data types that this task depends on
- `OutputType()`: Returns output data type that this task produces
- `Timeout()`: Returns task timeout duration
- `Execute(ctx, collection)`: Executes task logic

## Main Components

### Factory[CT ICollection]
Factory class for registering tasks and creating task flows
```go
type Factory[CT ICollection] struct {}
func NewFactory[CT ICollection]() *Factory[CT] {}
func (f *Factory[CT]) RegisterTask(createFunc TaskCreateFunc[CT]) error {} // Register task
func (f *Factory[CT]) CreateGraph() {} // Create dependency graph
func (f *Factory[CT]) CreateTaskDagflow(collection CT) (*TaskDagflow[CT], error) {} // Create task flow
```

### TaskDagflow[CT ICollection]
Task flow executor that manages concurrent execution of tasks, typically created from factory
```go
type TaskDagflow[CT ICollection] struct {}
func NewTaskDagflow[CT ICollection](metas []*taskMeta[CT], collection CT) (*TaskDagflow[CT], error) {}
func (td *TaskDagflow[CT]) Execute(ctx context.Context, timeout time.Duration) error {} // Execute task flow
```

## Helper Functions

### Automatic Type Inference
```go
func AutoInputTypes[CT ICollection]() []reflect.Type {} // Auto-infer input types (based on Get methods)
func AutoOutputType[CT ICollection]() reflect.Type {} // Auto-infer output type (based on Set methods)
```

### Task Creation
```go
type TaskCreateFunc[CT ICollection] func() (ITask[CT], error) // Task creation function type
func CreateTask[CT ICollection](createFunc TaskCreateFunc[CT]) (ITask[CT], error) {} // Create task instance
```

## Features

### Automatic Dependency Resolution
- Automatically builds DAG based on task `InputTypes()` and `OutputType()`
- Supports complex multi-level dependencies
- Automatically detects circular dependencies and unsatisfied dependencies

### Concurrent Execution Optimization
- Tasks without dependencies can execute concurrently
- Dynamic scheduling, triggers subsequent executable tasks immediately after task completion
- Supports both task-level and flow-level timeout control

### Type Safety
- Uses generics to ensure compile-time type safety
- Reflection mechanism for runtime type matching
- Automatic validation of input/output type consistency

## Usage Examples

### Basic Usage Flow

```go
import (
    "context"
    "time"
    "github.com/Steve-Lee-CST/go-pico-tool/pkg/task_dagflow"
)

// 1. Define data collection
type DataCollection struct {
    goods []Goods
    shops []Shop
    result GoodsInShops
}

func (c *DataCollection) InputTypes() []reflect.Type {
    return []reflect.Type{nil} // No initial input
}

func (c *DataCollection) TargetTypes() []reflect.Type {
    return []reflect.Type{reflect.TypeOf(GoodsInShops{})}
}

// 2. Implement concrete tasks
type GetGoodsTask struct {
    name string
    timeout time.Duration
}

func (t *GetGoodsTask) Name() string { return t.name }
func (t *GetGoodsTask) InputTypes() []reflect.Type { return []reflect.Type{nil} }
func (t *GetGoodsTask) OutputType() reflect.Type { return reflect.TypeOf([]Goods{}) }
func (t *GetGoodsTask) Timeout() time.Duration { return t.timeout }
func (t *GetGoodsTask) Execute(ctx context.Context, collection *DataCollection) error {
    // Logic to fetch goods data
    collection.goods = fetchGoods()
    return nil
}

// 3. Use factory to create task flow
func main() {
    factory := task_dagflow.NewFactory[*DataCollection]()
    
    // Register tasks
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &GetGoodsTask{name: "GetGoods", timeout: 500 * time.Millisecond}, nil
    })
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &GetShopsTask{name: "GetShops", timeout: 500 * time.Millisecond}, nil
    })
    factory.RegisterTask(func() (task_dagflow.ITask[*DataCollection], error) {
        return &ProcessTask{name: "Process", timeout: 1000 * time.Millisecond}, nil
    })
    
    // Create dependency graph
    factory.CreateGraph()
    
    // Create data collection and task flow
    collection := &DataCollection{}
    taskDagflow, err := factory.CreateTaskDagflow(collection)
    if err != nil {
        panic(err)
    }
    
    // Execute task flow
    ctx := context.Background()
    if err := taskDagflow.Execute(ctx, 2*time.Second); err != nil {
        panic(err)
    }
    
    // Use results
    fmt.Printf("Execution completed, result: %+v\n", collection.result)
}
```

### Using Automatic Type Inference

```go
// Define interface for automatic type inference
type IGoodsProvider interface {
    task_dagflow.ICollection
    GetGoods() []Goods    // Automatically recognized as input type
    SetResult(result GoodsInShops) // Automatically recognized as output type
}

type ProcessTask struct {
    name string
    timeout time.Duration
}

func (t *ProcessTask) InputTypes() []reflect.Type {
    return task_dagflow.AutoInputTypes[IGoodsProvider]() // Auto-infer input types
}

func (t *ProcessTask) OutputType() reflect.Type {
    return task_dagflow.AutoOutputType[IGoodsProvider]() // Auto-infer output type
}
```

### Complex Dependency Example

A complete example is shown in `demo_test.go`:
- `GetGoodsTask`: Fetch goods data (no dependencies)
- `GetShopsTask`: Fetch shop data (no dependencies)  
- `GoodsInShopsTask`: Process goods-shop relationships (depends on the first two tasks)

This example demonstrates:
- The first two tasks can execute concurrently
- The third task waits for the first two tasks to complete before executing
- Supports both task-level and flow-level timeout control
- Automatically handles complex data dependencies

## Important Notes

- About Tasks:
  - Task output types should be unique at the Factory level:
    - In other words, the output data type can serve as a unique identifier for a task
  - A task's output type should not be one of its own input types: cannot form self-loops
  - Tasks should not communicate directly with each other, only pass data through collections
  - Ensure task timeout settings are reasonable
- About Data Collections:
  - A specific data type can only be written by its corresponding task, other tasks can only read; based on this, data collections are thread-safe
  - It is recommended to use getter and setter methods to access data in collections
- About Factory and Dagflow:
  - Currently only detects unreachability, does not distinguish the cause: could be circular dependencies or missing task inputs
  - If target data has unreachable types, an error is returned.
