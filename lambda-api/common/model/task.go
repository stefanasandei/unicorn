package model

type ExecutionTaskStatus string

const (
	StatusDone   ExecutionTaskStatus = "successful" // everything worked
	StatusError                      = "error"      // something went wrong in the worker, likely a bug
	StatusFailed                     = "failed"     // either compilation error or runtime error
)

type ResponseTask struct {
	Status ExecutionTaskStatus `json:"status"`
	Output WorkerResponse      `json:"output"`
}
