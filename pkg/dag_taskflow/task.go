package dag_taskflow

import (
	"context"
	"reflect"
	"time"
)

type ICollection interface {
	InitialProvides() []reflect.Type
	TotalProvides() []reflect.Type
	FinalRequires() []reflect.Type
}

type ITask[CT ICollection] interface {
	Name() string
	Requires() []reflect.Type
	Provide() reflect.Type
	Timeout() time.Duration
	Execute(ctx context.Context, collection *CT) error
}

type TaskCreateFunc[CT ICollection] func() (ITask[CT], error)
