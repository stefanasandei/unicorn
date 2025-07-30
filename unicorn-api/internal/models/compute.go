package models

// ComputePreset defines hardware constraints for compute service
// e.g. Micro, Small

type ComputePreset string

const (
	PresetMicro ComputePreset = "micro"
	PresetSmall ComputePreset = "small"
)

// ComputeCreateRequest is the request body for creating a compute container
type ComputeCreateRequest struct {
	Name        string            `json:"name,omitempty"`
	Image       string            `json:"image" binding:"required"`
	Command     []string          `json:"command,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Preset      ComputePreset     `json:"preset,omitempty"`
	Ports       map[string]string `json:"ports,omitempty"`
	ExposePort  string            `json:"expose_port,omitempty"`
}

// ComputeContainerInfo holds info about a running container

type ComputeContainerInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	Ports     map[string]string `json:"ports"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
