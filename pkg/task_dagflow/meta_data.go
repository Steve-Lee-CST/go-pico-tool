package task_dagflow

import (
	"errors"
	"reflect"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type collectionMeta[CT ICollection] struct {
	InputTypes mapset.Set[reflect.Type]
	// AvailableTypes mapset.Set[reflect.Type]
	TargetTypes mapset.Set[reflect.Type]
}

func newCollectionMeta[CT ICollection](collection CT) (*collectionMeta[CT], error) {
	// nil is valid type for inputs and available
	inputs := mapset.NewSet[reflect.Type]()
	inputs.Append(collection.InputTypes()...)
	inputs.Add(nil)
	// available := mapset.NewSet[reflect.Type]()
	// available.Append(collection.AvailableTypes()...)
	// available.Add(nil)
	// outputs remove nil, though nil is valid
	targets := mapset.NewSet[reflect.Type]()
	targets.Append(collection.TargetTypes()...)
	targets.Remove(nil)
	// Validate the collection metadata
	if targets.IsEmpty() {
		return nil, errors.New("task flow must produce at least one output type")
	}
	// if !available.IsSuperset(inputs) {
	// 	return nil, errors.New("task flow inputs must be a subset of available types")
	// }
	// if !available.IsSuperset(targets) {
	// 	return nil, errors.New("task flow outputs must be a subset of available types")
	// }
	return &collectionMeta[CT]{
		InputTypes: inputs,
		// AvailableTypes: available,
		TargetTypes: targets,
	}, nil
}

type taskMeta[CT ICollection] struct {
	CreateFunc TaskCreateFunc[CT]
	Name       string
	InputTypes mapset.Set[reflect.Type]
	OutputType reflect.Type
	Timeout    time.Duration
}

func newTaskMeta[CT ICollection](createFunc TaskCreateFunc[CT]) (*taskMeta[CT], error) {
	task, err := CreateTask(createFunc)
	if err != nil {
		return nil, err
	}

	inputs := mapset.NewSet[reflect.Type]()
	inputs.Append(task.InputTypes()...)
	inputs.Add(nil) // nil is a valid input type
	outputType := task.OutputType()

	if outputType == nil {
		return nil, errors.New("task must produce a non-nil output type")
	}
	if inputs.Contains(outputType) {
		return nil, errors.New("task output type cannot be one of the input types")
	}

	return &taskMeta[CT]{
		CreateFunc: createFunc,
		Name:       task.Name(),
		InputTypes: inputs,
		OutputType: outputType,
		Timeout:    task.Timeout(),
	}, nil
}
