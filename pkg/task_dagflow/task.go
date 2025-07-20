package task_dagflow

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ICollection represents a collection of data to be used by tasks in the task flow.
// All inputs and outputs of tasks are expected to be of types defined in this collection.
// InputTypes(): returns the types of data that are available at the start of the task flow.
// AvailableTypes(): returns the all types of data that can be used by tasks
// TargetTypes(): returns the types of data that the task flow is expected to produce.
type ICollection interface {
	InputTypes() []reflect.Type
	// AvailableTypes() []reflect.Type
	TargetTypes() []reflect.Type
}

// ITask represents a task that can be executed within a task flow.
//
// ITask just like a function with the following signature:
//
//	func FuncName(param1 type1, param2 type2, ... , timeout time.Duration) (output outputType, error)
//
// Difference from a function:
//
//	params passed from CT (collection) instead of function parameters.
//	and output is also written to CT.
//
// Name(): returns the name of the task.
// InputTypes(): returns the types of inputs that this task depends on.
// OutputType(): returns the types of outputs that this task produces.
// Timeout(): returns the duration after which the task should be considered failed if not completed.
// Execute(ctx context.Context, collection *CT): executes the task with the provided context and collection.
type ITask[CT ICollection] interface {
	Name() string
	InputTypes() []reflect.Type
	OutputType() reflect.Type
	Timeout() time.Duration
	Execute(ctx context.Context, collection CT) error
}

type TaskCreateFunc[CT ICollection] func() (ITask[CT], error)

func CreateTask[CT ICollection](createFunc TaskCreateFunc[CT]) (ITask[CT], error) {
	task, err := createFunc()
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task created with %T returned nil", createFunc)
	}
	return task, nil
}

func AutoInputTypes[CT ICollection]() []reflect.Type {
	cElem := reflect.TypeOf((*CT)(nil)).Elem()
	var inputTypes []reflect.Type
	for i := 0; i < cElem.NumMethod(); i++ {
		method := cElem.Method(i)
		if strings.HasPrefix(method.Name, "Get") {
			inputTypes = append(inputTypes, method.Type.Out(0))
		}
	}
	return inputTypes
}

func AutoOutputType[CT ICollection]() reflect.Type {
	cElem := reflect.TypeOf((*CT)(nil)).Elem()
	for i := 0; i < cElem.NumMethod(); i++ {
		method := cElem.Method(i)
		if strings.HasPrefix(method.Name, "Set") && method.Type.NumIn() == 1 {
			return method.Type.In(0) // The first input is the collection itself
		}
	}
	return nil
}
