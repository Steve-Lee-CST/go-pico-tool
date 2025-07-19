package task_dagflow

import (
	"context"
	"fmt"
	"reflect"
	"time"

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
	ctx context.Context, collection *CT, resultChan chan *taskResult[CT],
) {
	startTime := time.Now()
	result := te.executeWithTimeout(ctx, collection, te.Meta.Timeout)
	result.TimeCost = time.Since(startTime)
	resultChan <- result
}

func (te *taskExecutor[CT]) executeWithTimeout(
	ctx context.Context, collection *CT, timeout time.Duration,
) *taskResult[CT] {
	resultChan := make(chan *taskResult[CT], 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultChan <- &taskResult[CT]{
					Meta: te.Meta,
					Err:  fmt.Errorf("task %s panicked: %v", te.Meta.Name, r),
				}
			}
		}()
		if err := te.Task.Execute(ctx, collection); err != nil {
			resultChan <- &taskResult[CT]{
				Meta: te.Meta,
				Err:  fmt.Errorf("task %s failed: %w", te.Meta.Name, err),
			}
			return
		}
		resultChan <- &taskResult[CT]{
			Meta: te.Meta,
			Err:  nil,
		}
	}()
	select {
	case <-ctx.Done():
		return &taskResult[CT]{
			Meta: te.Meta,
			Err:  fmt.Errorf("task %s cancelled: %w", te.Meta.Name, ctx.Err()),
		}
	case <-time.After(timeout):
		return &taskResult[CT]{
			Meta: te.Meta,
			Err:  fmt.Errorf("task %s timed out after %s", te.Meta.Name, timeout),
		}
	case result := <-resultChan:
		return result
	}
}
