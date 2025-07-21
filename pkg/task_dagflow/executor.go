package task_dagflow

import (
	"context"
	"fmt"
	"reflect"
	"time"

	tools "github.com/Steve-Lee-CST/go-pico-tool/tools"

	mapset "github.com/deckarep/golang-set/v2"
)

type taskResult[CT ICollection] struct {
	Meta     *taskMeta[CT]
	TimeCost time.Duration
	Err      error
}

type taskExecutor[CT ICollection] struct {
	Meta   *taskMeta[CT]
	Task   ITask[CT]
	Blocks mapset.Set[reflect.Type]
}

func newTaskExecutor[CT ICollection](meta *taskMeta[CT]) (*taskExecutor[CT], error) {
	task, err := CreateTask(meta.CreateFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create task %s: %w", meta.Name, err)
	}
	return &taskExecutor[CT]{
		Meta:   meta,
		Task:   task,
		Blocks: meta.InputTypes.Clone(),
	}, nil
}

func (te *taskExecutor[CT]) RemoveAndCheckBlock(block reflect.Type) bool {
	if te.Blocks.Contains(block) {
		te.Blocks.Remove(block)
		if te.Blocks.IsEmpty() {
			return true
		}
	}
	return false
}

func (te *taskExecutor[CT]) Execute(
	ctx context.Context, collection CT, resultChan chan *taskResult[CT],
) {
	startTime := time.Now()
	_, err := tools.RunFuncWithTimeout(
		ctx, te.Meta.Timeout,
		func(subCtx context.Context) (interface{}, error) {
			return struct{}{}, te.Task.Execute(ctx, collection)
		},
	)
	resultChan <- &taskResult[CT]{
		Meta:     te.Meta,
		Err:      err,
		TimeCost: time.Since(startTime),
	}
}
