package dag_taskflow

import (
	"context"
	"fmt"
	"reflect"
	"time"

	set "github.com/deckarep/golang-set/v2"
)

type TaskMeta[CT ICollection] struct {
	CreateFunc     TaskCreateFunc[CT]
	CreateFuncType reflect.Type // Type of the create function

	Name     string
	Requires set.Set[reflect.Type] // Types that this task requires
	Provide  reflect.Type
	Timeout  time.Duration
}

func NewTaskMeta[CT ICollection](createFunc TaskCreateFunc[CT]) (*TaskMeta[CT], error) {
	createFuncType := reflect.TypeOf(createFunc)
	task, err := createFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to create task with %s: %w", createFuncType.Name(), err)
	}
	if task == nil {
		return nil, fmt.Errorf("task create with %s returned nil", createFuncType.Name())
	}
	if task.Provide() == nil {
		return nil, fmt.Errorf("task %s must provide a non-nil type", task.Name())
	}
	requires := set.NewSet[reflect.Type]()
	requires.Append(task.Requires()...)
	requires.Remove(nil)
	if requires.Contains(task.Provide()) {
		return nil, fmt.Errorf("task %s cannot require its own provide type %s", task.Name(), task.Provide())
	}

	return &TaskMeta[CT]{
		CreateFunc:     createFunc,
		CreateFuncType: createFuncType,
		Name:           task.Name(),
		Requires:       requires,
		Provide:        task.Provide(),
		Timeout:        task.Timeout(),
	}, nil
}

type TaskExecutor[CT ICollection] struct {
	Meta *TaskMeta[CT]
	Task ITask[CT]

	Blocks set.Set[reflect.Type]
}

func NewTaskExecutor[CT ICollection](meta *TaskMeta[CT]) (*TaskExecutor[CT], error) {
	task, err := meta.CreateFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to create task with %s: %w", meta.CreateFuncType.Name(), err)
	}
	if task == nil {
		return nil, fmt.Errorf("task create with %s returned nil", meta.CreateFuncType.Name())
	}

	return &TaskExecutor[CT]{
		Meta:   meta,
		Task:   task,
		Blocks: meta.Requires.Clone(),
	}, nil
}

func (te *TaskExecutor[CT]) RemoveAndCheckBlock(blockType reflect.Type) bool {
	if te.Blocks.Contains(blockType) {
		te.Blocks.Remove(blockType)
		if te.Blocks.IsEmpty() {
			return true
		}
	}
	return false
}

func (te *TaskExecutor[CT]) Execute(
	ctx context.Context, collection *CT, resultChan chan<- TaskResult[CT],
) {
	executeStatusChan := make(chan ExecuteStatus, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				executeStatusChan <- ExecuteStatusEnum.Panic.WithError(
					fmt.Errorf("task %s panic: %v", te.Meta.Name, r),
				)
			}
		}()
		err := te.Task.Execute(ctx, collection)
		if err != nil {
			executeStatusChan <- ExecuteStatusEnum.Error.WithError(err)
			return
		}
		executeStatusChan <- ExecuteStatusEnum.Success
	}()
	select {
	case <-ctx.Done():
		resultChan <- TaskResult[CT]{
			Meta:   te.Meta,
			Status: ExecuteStatusEnum.Cancelled.WithError(ctx.Err()),
		}
	case <-time.After(te.Meta.Timeout):
		resultChan <- TaskResult[CT]{
			Meta:   te.Meta,
			Status: ExecuteStatusEnum.Timeout,
		}
	case status := <-executeStatusChan:
		resultChan <- TaskResult[CT]{
			Meta:   te.Meta,
			Status: status,
		}
	}
}
