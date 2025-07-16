package dag_taskflow

type ExecuteStatus struct {
	Code int
	Msg  string
	Err  error
}

func (s ExecuteStatus) WithError(err error) ExecuteStatus {
	s.Err = err
	return s
}

var ExecuteStatusEnum = struct {
	Success   ExecuteStatus
	Error     ExecuteStatus
	Panic     ExecuteStatus
	Timeout   ExecuteStatus
	Cancelled ExecuteStatus
}{
	Success:   ExecuteStatus{Code: 0, Msg: "success"},
	Error:     ExecuteStatus{Code: 1, Msg: "error"},
	Panic:     ExecuteStatus{Code: 2, Msg: "panic"},
	Timeout:   ExecuteStatus{Code: 3, Msg: "timeout"},
	Cancelled: ExecuteStatus{Code: 4, Msg: "cancelled"},
}

type TaskResult[CT ICollection] struct {
	Meta   *TaskMeta[CT]
	Status ExecuteStatus
}
