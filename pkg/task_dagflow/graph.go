package task_dagflow

import (
	"errors"
	"reflect"

	"github.com/Steve-Lee-CST/go-pico-tool/tools"
	mapset "github.com/deckarep/golang-set/v2"
)

type node[CT ICollection] struct {
	Meta      *taskMeta[CT]
	Reachable bool
}

func newNode[CT ICollection](meta *taskMeta[CT]) *node[CT] {
	return &node[CT]{
		Meta:      meta,
		Reachable: false,
	}
}

type graph[CT ICollection] struct {
	outputToNode   map[reflect.Type]*node[CT]
	inputToNodes   map[reflect.Type][]*node[CT]
	outputToInputs map[reflect.Type]mapset.Set[reflect.Type]
}

func newGraph[CT ICollection](metas []*taskMeta[CT]) *graph[CT] {
	g := &graph[CT]{
		outputToNode:   make(map[reflect.Type]*node[CT]),
		inputToNodes:   make(map[reflect.Type][]*node[CT]),
		outputToInputs: make(map[reflect.Type]mapset.Set[reflect.Type]),
	}
	for _, meta := range metas {
		n := newNode(meta)
		g.outputToNode[meta.OutputType] = n
		for inputType := range meta.InputTypes.Iter() {
			g.inputToNodes[inputType] = append(g.inputToNodes[inputType], n)
		}
	}
	for output := range g.outputToNode {
		g.outputToInputs[output] = g.getInputs(output)
	}

	return g
}

// 并查集查找 outputType 的所有输入类型
func (g *graph[CT]) getInputs(outputType reflect.Type) mapset.Set[reflect.Type] {
	inputs, exists := g.outputToInputs[outputType]
	if !exists {
		return nil
	}

	inputList := tools.NewQueue[reflect.Type]()
	for inputType := range inputs.Iter() {
		inputList.Enqueue(inputType)
	}
	for inputList.Size() > 0 {
		inputType, exist := inputList.Dequeue()
		if !exist {
			break
		}
		// find in graph
		if meta, ok := g.outputToNode[inputType]; ok {
			for childInputType := range meta.Meta.InputTypes.Iter() {
				if !inputs.Contains(childInputType) {
					inputs.Add(childInputType)
					inputList.Enqueue(childInputType)
				}
			}
		}
	}
	return inputs
}

func (g *graph[CT]) calReachStatus(
	collectionMeta *collectionMeta[CT],
) (reachableTypes, unReachableTypes mapset.Set[reflect.Type]) {
	inputs := mapset.NewSet[reflect.Type]()
	inputs.Append(collectionMeta.InputTypes.ToSlice()...)
	inputs.Add(nil)
	targets := mapset.NewSet[reflect.Type]()
	targets.Append(collectionMeta.TargetTypes.ToSlice()...)

	reachableTypes = mapset.NewSet[reflect.Type]()
	for _, task := range g.outputToNode {
		task.Reachable = false
	}

	for inputType := range inputs.Iter() {
		reachableTypes.Add(inputType)
	}
	newNodeTag := true
	for newNodeTag {
		newNodeTag = false
		for output, node := range g.outputToNode {
			if node.Meta.InputTypes.IsSubset(reachableTypes) && !reachableTypes.Contains(output) {
				reachableTypes.Add(output)
				node.Reachable = true
				newNodeTag = true
			}
		}
	}
	unReachableTypes = targets.Difference(reachableTypes)
	return
}

func (g *graph[CT]) GetMinTaskMetas(collection CT) ([]*taskMeta[CT], error) {
	collectionMeta, err := newCollectionMeta(collection)
	if err != nil {
		return nil, err
	}

	reachableTypes, unReachableTypes := g.calReachStatus(collectionMeta)
	if unReachableTypes.Cardinality() > 0 {
		return nil, errors.New("task flow has unreachable output types: " + reflect.TypeOf(unReachableTypes).String())
	}

	taskMetas := make([]*taskMeta[CT], 0)
	for outputType := range reachableTypes.Iter() {
		if node, exists := g.outputToNode[outputType]; exists {
			taskMetas = append(taskMetas, node.Meta)
		}
	}
	return taskMetas, nil
}
