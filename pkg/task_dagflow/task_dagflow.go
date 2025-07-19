package task_dagflow

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type TaskDagflow[CT ICollection] struct {
	collection     *CT
	collectionMeta *collectionMeta[CT]

	metas        []*taskMeta[CT]
	tasks        []*taskExecutor[CT]
	inputToTasks map[reflect.Type][]*taskExecutor[CT]
	timeCost     time.Duration

	lock sync.Mutex
}

func NewTaskDagflow[CT ICollection](metas []*taskMeta[CT], collection *CT) (*TaskDagflow[CT], error) {
	if collection == nil {
		return nil, errors.New("collection cannot be nil")
	}
	collectionMeta, err := newCollectionMeta(*collection)
	if err != nil {
		return nil, err
	}

	tasks := make([]*taskExecutor[CT], 0, len(metas))
	for _, meta := range metas {
		task, err := newTaskExecutor(meta)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	inputToTasks := make(map[reflect.Type][]*taskExecutor[CT], 0)
	for _, task := range tasks {
		for inputType := range task.Meta.InputTypes.Iter() {
			inputToTasks[inputType] = append(inputToTasks[inputType], task)
		}
	}
	return &TaskDagflow[CT]{
		collection:     collection,
		collectionMeta: collectionMeta,

		metas:        metas,
		tasks:        tasks,
		inputToTasks: inputToTasks,
		timeCost:     0,

		lock: sync.Mutex{},
	}, nil
}

func (t *TaskDagflow[CT]) initUnblockTypeChan(collectionMeta *collectionMeta[CT]) chan reflect.Type {
	unblockTypeChan := make(chan reflect.Type, len(t.metas)+2) // +2 for nil: ensure no-chan-block
	unblockTypeChan <- nil
	for initType := range collectionMeta.InputTypes.Iter() {
		unblockTypeChan <- initType
	}
	return unblockTypeChan
}

func (t *TaskDagflow[CT]) Execute(ctx context.Context, timeout time.Duration) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	startTime := time.Now()

	unblockTypeChan := t.initUnblockTypeChan(t.collectionMeta)
	resultChan := make(chan *taskResult[CT], len(t.metas)+1) // ensure no-chan-block
	taskRecord, resultRecord := mapset.NewSet[reflect.Type](), mapset.NewSet[reflect.Type]()
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer func() {
		t.timeCost = time.Since(startTime)
	}()
	for {
		select {
		case <-subCtx.Done():
			return subCtx.Err()
		case <-time.After(timeout):
			return errors.New("task dagflow execution timed out")
		case unblockType := <-unblockTypeChan:
			tasks := t.inputToTasks[unblockType]
			for _, task := range tasks {
				if task.RemoveAndCheckBlock(unblockType) {
					taskRecord.Add(task.Meta.OutputType)
					go task.Execute(subCtx, t.collection, resultChan)
				}
			}
			if taskRecord.Equal(t.collectionMeta.TargetTypes) {
				return nil
			}
		case result := <-resultChan:
			if result == nil {
				return errors.New("received nil result from task execution")
			}
			resultRecord.Add(result.Meta.OutputType)
			if result.Err != nil {
				return fmt.Errorf("task %s failed: %w", result.Meta.Name, result.Err)
			}
		}
	}
}

func (t *TaskDagflow[CT]) TimeCost() time.Duration {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.timeCost
}

func (t *TaskDagflow[CT]) Tasks() []*taskExecutor[CT] {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.tasks
}
