package model

import "github.com/google/uuid"

type File struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

type Permissions struct {
	CanRead       bool `json:"read,omitempty"`
	CanWrite      bool `json:"write,omitempty"`
	NetworkAccess bool `json:"network,omitempty"`
}

// ProcessInfo TODO
type ProcessInfo struct {
	StandardInput string `json:"stdin,omitempty"`
	// ExpectsFileOutput bool   `json:"file_output,omitempty"`

	CPUTime string `json:"time,omitempty"`
	// Memory         string `json:"memory,omitempty"`
	// MaxFileSize    string `json:"max_file_size,omitempty"`
	MaxOpenedFiles int32 `json:"max_opened_files,omitempty"` // TODO
	MaxProcesses   int32 `json:"max_processes,omitempty"`    // TODO

	Permissions Permissions `json:"permissions,omitempty"` // TODO

	EnvironmentVars map[string]string `json:"env,omitempty"`

	WorkingDirectory string `json:"-"`
}

type ExecutionRequest struct {
	Runtime struct {
		Name    string `json:"name"`
		Version string `json:"version,omitempty"`
	} `json:"runtime"`
	Project struct {
		Entry string `json:"entry,omitempty"`
		Files []File `json:"files"`
	} `json:"project"`
	Process ProcessInfo `json:"process,omitempty"`
}

type ExecutionRequestWrapper struct {
	Id  uuid.UUID
	Req ExecutionRequest
}

/*
Example request:
	{
		"runtime": {
			"name": "python3",
			"version": "3.12"
		},
		"project": {
			"entry": "import utils\nprint(add(1,2))",
			"files": [{
				"name": "utils.py",
				"contents": "def add(a, b):\n\treturn a + b"
			}]
		},
		"process": {
			"time": "2s",
			"permissions": {
				"read": true
			}
		}
	}

Implemented options:
	- runtime name
	- project entry point
	- project files
	- process limits
	- process permissions
	- file io

	- fix docker deployment, broken because of the common lib

Notes:
	- this is for a "execute and run" approach, small projects
	- only the "entry" file is compiled, the rest of the files are either resources or don't require compilation
	- use the interactive project runner for complex projects (websockets instead of rest, WIP)
*/
