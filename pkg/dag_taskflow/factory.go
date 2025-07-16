package dag_taskflow

import (
	"fmt"
	"reflect"
)

type DagTaskflowFactory[CT ICollection] struct {
	ProvideToMeta  map[reflect.Type]*TaskMeta[CT]
	RequireToMetas map[reflect.Type][]*TaskMeta[CT]
}

func NewDagTaskflowFactory[CT ICollection]() *DagTaskflowFactory[CT] {
	return &DagTaskflowFactory[CT]{
		ProvideToMeta:  make(map[reflect.Type]*TaskMeta[CT]),
		RequireToMetas: make(map[reflect.Type][]*TaskMeta[CT]),
	}
}

func (f *DagTaskflowFactory[CT]) RegisterTask(taskCreator TaskCreateFunc[CT]) error {
	meta, err := NewTaskMeta(taskCreator)
	if err != nil {
		return err
	}
	if _, exists := f.ProvideToMeta[meta.Provide]; exists {
		return fmt.Errorf("task with provide type %s already registered", meta.Provide)
	}
	// update maps
	f.ProvideToMeta[meta.Provide] = meta
	for requireType := range meta.Requires.Iter() {
		f.RequireToMetas[requireType] = append(f.RequireToMetas[requireType], meta)
	}
	// especially for requireType == nil or empty
	if meta.Requires.Cardinality() == 0 {
		f.RequireToMetas[nil] = append(f.RequireToMetas[nil], meta)
	}

	return nil
}

func (f *DagTaskflowFactory[CT]) Create(collection *CT) (*DagTaskflow[CT], error) {

	return NewDagTaskflow[CT](collection), nil
}
