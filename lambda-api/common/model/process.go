package model

import "github.com/google/uuid"

type ProcessResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Output string `json:"output"`

	Time   int32  `json:"time"`   // ms
	Memory uint64 `json:"memory"` // bytes

	ExitCode int32 `json:"exit_code"`
}

type WorkerResponse struct {
	Compile ProcessResult `json:"compile"`
	Run     ProcessResult `json:"run"`
}

type WorkerResponseWrapper struct {
	Id  uuid.UUID
	Res WorkerResponse
}
