package task_dagflow

import (
	"fmt"
	"reflect"
)

type Factory[CT ICollection] struct {
	outputToTaskMeta map[reflect.Type]*taskMeta[CT]
	graph            *graph[CT]
}

func NewFactory[CT ICollection]() *Factory[CT] {
	return &Factory[CT]{
		outputToTaskMeta: make(map[reflect.Type]*taskMeta[CT]),
	}
}

func (f *Factory[CT]) RegisterTask(createFunc TaskCreateFunc[CT]) error {
	meta, err := newTaskMeta(createFunc)
	if err != nil {
		return err
	}
	if _, exists := f.outputToTaskMeta[meta.OutputType]; exists {
		return fmt.Errorf("task with output type %s already registered", meta.OutputType)
	}
	f.outputToTaskMeta[meta.OutputType] = meta
	return nil
}

func (f *Factory[CT]) CreateGraph() {
	tasks := make([]*taskMeta[CT], 0)
	for _, meta := range f.outputToTaskMeta {
		tasks = append(tasks, meta)
	}

	if f.graph == nil {
		f.graph = newGraph(tasks)
	}
}

func (f *Factory[CT]) CreateTaskDagflow(collection CT) (*TaskDagflow[CT], error) {
	metas, err := f.graph.GetMinTaskMetas(collection)
	if err != nil {
		return nil, err
	}

	return NewTaskDagflow(metas, collection)
}
