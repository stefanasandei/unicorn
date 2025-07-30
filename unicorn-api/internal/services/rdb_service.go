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
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// RDBService handles database container operations
type RDBService struct {
	validator *validation.Validator
}

// NewRDBService creates a new RDB service
func NewRDBService() *RDBService {
	return &RDBService{
		validator: validation.NewValidator(),
	}
}

// CreateRDBInstance creates a new database container
func (s *RDBService) CreateRDBInstance(userID uuid.UUID, req models.RDBCreateRequest) (*models.RDBInstanceInfo, error) {
	// Set defaults if not provided
	if req.Preset == "" {
		req.Preset = models.RDBPresetMicro
	}
	if req.Port == "" {
		req.Port = s.getDefaultPort(req.Type)
	}
	if req.Database == "" {
		req.Database = "main"
	}
	if req.Username == "" {
		req.Username = "user"
	}
	if req.Password == "" {
		req.Password = s.generatePassword()
	}

	// Validate port
	if err := s.validator.ValidatePort(req.Port); err != nil {
		return nil, err
	}

	// Validate volumes
	for _, volume := range req.Volumes {
		if volume.Size < 1 || volume.Size > 100000 {
			return nil, errors.ErrBadRequest.WithDetails("Volume size must be between 1MB and 100GB")
		}
		if volume.Name == "" {
			return nil, errors.ErrBadRequest.WithDetails("Volume name is required")
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

	// Get image and configuration based on database type
	image, envVars := s.getDatabaseConfig(req)

	// Pull image with retry logic
	if err := s.pullImageWithRetry(ctx, cli, image); err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to pull Docker image: " + err.Error())
	}

	// Set resource limits based on preset
	resources := s.getResourceLimits(req.Preset)

	// Configure port bindings
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	port, _ := nat.NewPort("tcp", req.Port)
	exposedPorts[port] = struct{}{}
	hostPort := fmt.Sprintf("%d", 10000+rand.Intn(10000))
	portBindings[port] = []nat.PortBinding{{HostPort: hostPort}}

	// Configure volumes with automatic mount paths for database storage
	var mounts []mount.Mount
	for _, volume := range req.Volumes {
		volumeName := fmt.Sprintf("rdb-%s-%s", userID.String()[:8], volume.Name)

		// Determine mount path based on database type
		var mountPath string
		if req.Type == models.RDBTypePostgreSQL {
			mountPath = "/var/lib/postgresql/data"
		} else {
			mountPath = "/var/lib/mysql"
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: volumeName,
			Target: mountPath,
		})
	}

	// Set environment variables
	envVars = append(envVars, fmt.Sprintf("POSTGRES_DB=%s", req.Database))
	envVars = append(envVars, fmt.Sprintf("POSTGRES_USER=%s", req.Username))
	envVars = append(envVars, fmt.Sprintf("POSTGRES_PASSWORD=%s", req.Password))

	// Add custom environment variables
	for key, value := range req.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Generate container name
	containerName := req.Name
	if containerName == "" {
		containerName = fmt.Sprintf("rdb-%s-%s", userID.String()[:8], s.randString(8))
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Env:          envVars,
		ExposedPorts: exposedPorts,
		Labels:       map[string]string{"owner": userID.String(), "type": "rdb"},
	}, &container.HostConfig{
		PortBindings: portBindings,
		Resources:    resources,
		Mounts:       mounts,
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
	now := time.Now().Format(time.RFC3339)
	return &models.RDBInstanceInfo{
		ID:          resp.ID,
		Name:        containerName,
		Type:        req.Type,
		Status:      containerInfo.State.Status,
		Port:        hostPort,
		Host:        "localhost",
		Database:    req.Database,
		Username:    req.Username,
		Volumes:     req.Volumes,
		Environment: req.Environment,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ListRDBInstances lists all database containers for a user
func (s *RDBService) ListRDBInstances(userID uuid.UUID) ([]models.RDBInstanceInfo, error) {
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
	filter.Add("label", "type=rdb")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to list containers: " + err.Error())
	}

	var result []models.RDBInstanceInfo
	for _, ctr := range containers {
		// Get detailed container info
		_, err = cli.ContainerInspect(ctx, ctr.ID)
		if err != nil {
			continue
		}

		// Extract port information
		var port string
		for _, p := range ctr.Ports {
			if p.PublicPort != 0 {
				port = fmt.Sprintf("%d", p.PublicPort)
				break
			}
		}

		// Determine database type from image
		dbType := s.getDatabaseTypeFromImage(ctr.Image)

		now := time.Now().Format(time.RFC3339)
		result = append(result, models.RDBInstanceInfo{
			ID:        ctr.ID,
			Name:      ctr.Names[0][1:], // Remove leading slash
			Type:      dbType,
			Status:    ctr.Status,
			Port:      port,
			Host:      "localhost",
			Database:  "main", // Default, could be extracted from env vars
			Username:  "user", // Default, could be extracted from env vars
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	return result, nil
}

// DeleteRDBInstance deletes a database container
func (s *RDBService) DeleteRDBInstance(userID uuid.UUID, containerID string) error {
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

	// Check if container belongs to user and is an RDB container
	if containerInfo.Config.Labels["owner"] != userID.String() ||
		containerInfo.Config.Labels["type"] != "rdb" {
		return errors.ErrForbidden.WithDetails("Container does not belong to user or is not an RDB container")
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

// getDatabaseConfig returns the appropriate Docker image and environment variables for the database type
func (s *RDBService) getDatabaseConfig(req models.RDBCreateRequest) (string, []string) {
	switch req.Type {
	case models.RDBTypePostgreSQL:
		return "postgres:15-alpine", []string{
			"POSTGRES_HOST_AUTH_METHOD=trust",
		}
	case models.RDBTypeMySQL:
		return "mysql:8.0", []string{
			"MYSQL_ROOT_PASSWORD=" + req.Password,
			"MYSQL_DATABASE=" + req.Database,
			"MYSQL_USER=" + req.Username,
			"MYSQL_PASSWORD=" + req.Password,
		}
	default:
		return "postgres:15-alpine", []string{
			"POSTGRES_HOST_AUTH_METHOD=trust",
		}
	}
}

// getDefaultPort returns the default port for the database type
func (s *RDBService) getDefaultPort(dbType models.RDBType) string {
	switch dbType {
	case models.RDBTypePostgreSQL:
		return "5432"
	case models.RDBTypeMySQL:
		return "3306"
	default:
		return "5432"
	}
}

// getDatabaseTypeFromImage determines the database type from the Docker image
func (s *RDBService) getDatabaseTypeFromImage(image string) models.RDBType {
	if len(image) >= 8 && image[:8] == "mysql" {
		return models.RDBTypeMySQL
	}
	return models.RDBTypePostgreSQL
}

// getResourceLimits returns resource limits based on preset
func (s *RDBService) getResourceLimits(preset models.RDBPreset) container.Resources {
	switch preset {
	case models.RDBPresetMicro:
		return container.Resources{
			NanoCPUs: 500_000_000,       // 0.5 CPU
			Memory:   512 * 1024 * 1024, // 512MB
		}
	case models.RDBPresetSmall:
		return container.Resources{
			NanoCPUs: 1_000_000_000,          // 1 CPU
			Memory:   1 * 1024 * 1024 * 1024, // 1GB
		}
	case models.RDBPresetMedium:
		return container.Resources{
			NanoCPUs: 2_000_000_000,          // 2 CPU
			Memory:   2 * 1024 * 1024 * 1024, // 2GB
		}
	default:
		return container.Resources{}
	}
}

// pullImageWithRetry pulls a Docker image with retry logic
func (s *RDBService) pullImageWithRetry(ctx context.Context, cli *client.Client, image string) error {
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

// generatePassword generates a random password
func (s *RDBService) generatePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// randString generates a random string
func (s *RDBService) randString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
