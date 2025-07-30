package services

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"unicorn-api/internal/common/errors"
	"unicorn-api/internal/common/validation"
	"unicorn-api/internal/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// ComputeService handles Docker container operations
type ComputeService struct {
	validator *validation.Validator
}

// NewComputeService creates a new compute service
func NewComputeService() *ComputeService {
	return &ComputeService{
		validator: validation.NewValidator(),
	}
}

// CreateContainer creates a new Docker container
func (s *ComputeService) CreateContainer(userID uuid.UUID, req models.ComputeCreateRequest) (*models.ComputeContainerInfo, error) {
	// Validate Docker image
	if err := s.validator.ValidateDockerImage(req.Image); err != nil {
		return nil, err
	}

	// Set defaults if not provided
	if req.Preset == "" {
		req.Preset = models.PresetMicro
	}
	if req.Ports == nil {
		req.Ports = make(map[string]string)
	}
	if req.ExposePort == "" {
		req.ExposePort = "80"
	}

	// Validate ports
	for port := range req.Ports {
		if err := s.validator.ValidatePort(port); err != nil {
			return nil, err
		}
	}

	// Create Docker client with proper socket path for macOS
	var cli *client.Client
	var err error
	if runtime.GOOS == "darwin" {
		// On macOS, Docker Desktop uses a different socket path
		cli, err = client.NewClientWithOpts(
			client.WithHost("unix:///Users/asandeistefan/.docker/run/docker.sock"),
			client.WithAPIVersionNegotiation(),
		)
	} else {
		cli, err = client.NewClientWithOpts(client.FromEnv)
	}
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Docker client unavailable: " + err.Error())
	}
	defer cli.Close()

	ctx := context.Background()

	// Pull image with retry logic
	if err := s.pullImageWithRetry(ctx, cli, req.Image); err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to pull Docker image: " + err.Error())
	}

	// Set resource limits based on preset
	resources := s.getResourceLimits(req.Preset)

	// Configure port bindings
	exposedPorts := s.natPortSet(req.Ports)
	portBindings := s.natPortBindings(req.Ports)

	// Handle exposed port
	if req.ExposePort != "" {
		if _, exists := req.Ports[req.ExposePort]; !exists {
			hostPort := fmt.Sprintf("%d", 10000+rand.Intn(10000))
			req.Ports[req.ExposePort] = hostPort
			exposedPorts = s.natPortSet(req.Ports)
			portBindings = s.natPortBindings(req.Ports)
		}
	}

	// Set environment variables
	envVars := []string{}
	for key, value := range req.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Generate container name
	containerName := req.Name
	if containerName == "" {
		containerName = fmt.Sprintf("compute-%s-%s", userID.String()[:8], s.randString(8))
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        req.Image,
		Cmd:          req.Command,
		Env:          envVars,
		ExposedPorts: exposedPorts,
		Labels:       map[string]string{"owner": userID.String()},
	}, &container.HostConfig{
		PortBindings: portBindings,
		Resources:    resources,
	}, nil, nil, containerName)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Container creation failed: " + err.Error())
	}

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, errors.ErrInternalError.WithDetails("Container start failed: " + err.Error())
	}

	// Get container info
	containerInfo, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to inspect container: " + err.Error())
	}

	// Build response
	ports := make(map[string]string)
	for port, bindings := range containerInfo.NetworkSettings.Ports {
		if len(bindings) > 0 {
			ports[port.Port()] = bindings[0].HostPort
		}
	}

	now := time.Now().Format(time.RFC3339)
	return &models.ComputeContainerInfo{
		ID:        resp.ID,
		Name:      containerName,
		Image:     req.Image,
		Status:    containerInfo.State.Status,
		Ports:     ports,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ListContainers lists all containers for a user
func (s *ComputeService) ListContainers(userID uuid.UUID) ([]models.ComputeContainerInfo, error) {
	var cli *client.Client
	var err error
	if runtime.GOOS == "darwin" {
		// On macOS, Docker Desktop uses a different socket path
		cli, err = client.NewClientWithOpts(
			client.WithHost("unix:///Users/asandeistefan/.docker/run/docker.sock"),
			client.WithAPIVersionNegotiation(),
		)
	} else {
		cli, err = client.NewClientWithOpts(client.FromEnv)
	}
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Docker client unavailable: " + err.Error())
	}
	defer cli.Close()

	ctx := context.Background()
	filter := filters.NewArgs()
	filter.Add("label", "owner="+userID.String())

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to list containers: " + err.Error())
	}

	var result []models.ComputeContainerInfo
	for _, ctr := range containers {
		ports := make(map[string]string)
		for _, p := range ctr.Ports {
			if p.PublicPort != 0 {
				ports[fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)] = fmt.Sprintf("%d", p.PublicPort)
			}
		}

		now := time.Now().Format(time.RFC3339)
		result = append(result, models.ComputeContainerInfo{
			ID:        ctr.ID,
			Name:      ctr.Names[0][1:], // Remove leading slash
			Image:     ctr.Image,
			Status:    ctr.Status,
			Ports:     ports,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	return result, nil
}

// DeleteContainer deletes a Docker container
func (s *ComputeService) DeleteContainer(userID uuid.UUID, containerID string) error {
	var cli *client.Client
	var err error
	if runtime.GOOS == "darwin" {
		// On macOS, Docker Desktop uses a different socket path
		cli, err = client.NewClientWithOpts(
			client.WithHost("unix:///Users/asandeistefan/.docker/run/docker.sock"),
			client.WithAPIVersionNegotiation(),
		)
	} else {
		cli, err = client.NewClientWithOpts(client.FromEnv)
	}
	if err != nil {
		return errors.ErrInternalError.WithDetails("Docker client unavailable: " + err.Error())
	}
	defer cli.Close()

	ctx := context.Background()

	// Verify the container belongs to the user
	containerInfo, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return errors.ErrNotFound.WithDetails("Container not found")
	}

	// Check if container belongs to user
	if containerInfo.Config.Labels["owner"] != userID.String() {
		return errors.ErrForbidden.WithDetails("Container does not belong to user")
	}

	// Stop the container if it's running
	if containerInfo.State.Running {
		if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
			return errors.ErrInternalError.WithDetails("Failed to stop container: " + err.Error())
		}
	}

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		return errors.ErrInternalError.WithDetails("Failed to remove container: " + err.Error())
	}

	return nil
}

// pullImageWithRetry pulls a Docker image with retry logic
func (s *ComputeService) pullImageWithRetry(ctx context.Context, cli *client.Client, image string) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		_, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return fmt.Errorf("failed to pull image after %d attempts", maxRetries)
}

// getResourceLimits returns resource limits based on preset
func (s *ComputeService) getResourceLimits(preset models.ComputePreset) container.Resources {
	switch preset {
	case models.PresetMicro:
		return container.Resources{
			NanoCPUs: 500_000_000,       // 0.5 CPU
			Memory:   256 * 1024 * 1024, // 256MB
		}
	case models.PresetSmall:
		return container.Resources{
			NanoCPUs: 1_000_000_000,     // 1 CPU
			Memory:   512 * 1024 * 1024, // 512MB
		}
	default:
		return container.Resources{}
	}
}

// natPortSet converts port map to Docker port set
func (s *ComputeService) natPortSet(ports map[string]string) nat.PortSet {
	ps := nat.PortSet{}
	for cport := range ports {
		p, _ := nat.NewPort("tcp", cport)
		ps[p] = struct{}{}
	}
	return ps
}

// natPortBindings converts port map to Docker port bindings
func (s *ComputeService) natPortBindings(ports map[string]string) nat.PortMap {
	pm := nat.PortMap{}
	for cport, hport := range ports {
		p, _ := nat.NewPort("tcp", cport)
		pm[p] = []nat.PortBinding{{HostPort: hport}}
	}
	return pm
}

// randString generates a random string
func (s *ComputeService) randString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
