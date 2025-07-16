package dag_taskflow

// 一次性的运行容器, 包括任务图的构建和执行
type DagTaskflow[CT ICollection] struct {
	collection *CT // 任务执行的集合
}

func NewDagTaskflow[CT ICollection](collection *CT) *DagTaskflow[CT] {
	return &DagTaskflow[CT]{collection: collection}
}
