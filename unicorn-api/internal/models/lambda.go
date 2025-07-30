package models

// LambdaFile represents a file in a Lambda function
type LambdaFile struct {
	Name     string `json:"name" binding:"required" example:"utils.py"`
	Contents string `json:"contents" binding:"required" example:"def add(a, b):\n\treturn a + b"`
}

// LambdaPermissions defines what the Lambda function is allowed to do
type LambdaPermissions struct {
	CanRead       bool `json:"read,omitempty" example:"true"`
	CanWrite      bool `json:"write,omitempty" example:"false"`
	NetworkAccess bool `json:"network,omitempty" example:"false"`
}

// LambdaProcessInfo contains execution parameters for the Lambda function
type LambdaProcessInfo struct {
	StandardInput    string            `json:"stdin,omitempty" example:"test input"`
	CPUTime          string            `json:"time,omitempty" example:"2s"`
	MaxOpenedFiles   int32             `json:"max_opened_files,omitempty" example:"10"`
	MaxProcesses     int32             `json:"max_processes,omitempty" example:"5"`
	Permissions      LambdaPermissions `json:"permissions,omitempty"`
	EnvironmentVars  map[string]string `json:"env,omitempty" example:"{\"DEBUG\":\"true\"}"`
	WorkingDirectory string            `json:"-"`
}

// LambdaExecuteRequest is the request body for executing a Lambda function
type LambdaExecuteRequest struct {
	Runtime struct {
		Name    string `json:"name" binding:"required" example:"python3"`
		Version string `json:"version,omitempty" example:"3.12"`
	} `json:"runtime"`
	Project struct {
		Entry string       `json:"entry,omitempty" example:"import utils\nprint(utils.add(1,2))"`
		Files []LambdaFile `json:"files" binding:"required"`
	} `json:"project"`
	Process LambdaProcessInfo `json:"process,omitempty"`
}

// LambdaExecuteResponse is the response from executing a Lambda function
type LambdaExecuteResponse struct {
	Status  string `json:"status" example:"success"`
	Output  string `json:"output" example:"3\n"`
	Runtime string `json:"runtime,omitempty" example:"python3"`
	Time    string `json:"time,omitempty" example:"0.023s"`
}
