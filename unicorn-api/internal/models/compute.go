package models

// ComputePreset defines hardware constraints for compute service
// e.g. Micro, Small

type ComputePreset string

const (
	PresetMicro ComputePreset = "micro"
	PresetSmall ComputePreset = "small"
)

// ComputeCreateRequest is the request body for creating a compute container
// Ports is a map of containerPort:hostPort
// ExposePort is the container port to expose

type ComputeCreateRequest struct {
	Image      string            `json:"image" binding:"required"`
	Preset     ComputePreset     `json:"preset" binding:"required,oneof=micro small"`
	Ports      map[string]string `json:"ports" binding:"required"`
	ExposePort string            `json:"expose_port" binding:"required"`
}

// ComputeContainerInfo holds info about a running container

type ComputeContainerInfo struct {
	ID     string            `json:"id"`
	Image  string            `json:"image"`
	Status string            `json:"status"`
	Ports  map[string]string `json:"ports"`
}
