package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"

	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"
)

type ComputeHandler struct {
	Config   *config.Config
	IAMStore stores.IAMStore
}

func NewComputeHandler(cfg *config.Config, iamStore stores.IAMStore) *ComputeHandler {
	return &ComputeHandler{Config: cfg, IAMStore: iamStore}
}

// POST /compute/create
func (h *ComputeHandler) CreateCompute(c *gin.Context) {
	var req models.ComputeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !h.hasPermission(claims, "compute", 1) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "docker unavailable"})
		return
	}
	defer cli.Close()

	ctx := context.Background()
	// Pull image if not present
	_, err = cli.ImagePull(ctx, req.Image, types.ImagePullOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to pull image"})
		return
	}

	// Set resource limits
	resources := container.Resources{}
	switch req.Preset {
	case models.PresetMicro:
		resources.NanoCPUs = 500_000_000     // 0.5 CPU
		resources.Memory = 256 * 1024 * 1024 // 256MB
	case models.PresetSmall:
		resources.NanoCPUs = 1_000_000_000   // 1 CPU
		resources.Memory = 512 * 1024 * 1024 // 512MB
	}

	// Port bindings
	exposedPorts := natPortSet(req.Ports)
	portBindings := natPortBindings(req.Ports)

	containerName := "compute-" + claims.AccountID + "-" + randString(8)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        req.Image,
		ExposedPorts: exposedPorts,
		Labels:       map[string]string{"owner": claims.AccountID},
	}, &container.HostConfig{
		PortBindings: portBindings,
		Resources:    resources,
	}, nil, nil, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "container create failed"})
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "container start failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": resp.ID})
}

// GET /compute/list
func (h *ComputeHandler) ListCompute(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !h.hasPermission(claims, "compute", 0) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "docker unavailable"})
		return
	}
	defer cli.Close()

	ctx := context.Background()
	filter := filters.NewArgs()
	filter.Add("label", "owner="+claims.AccountID)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	var result []models.ComputeContainerInfo
	for _, ctr := range containers {
		ports := map[string]string{}
		for _, p := range ctr.Ports {
			ports[formatPort(p.PrivatePort, p.Type)] = formatPort(p.PublicPort, p.Type)
		}
		result = append(result, models.ComputeContainerInfo{
			ID:     ctr.ID,
			Image:  ctr.Image,
			Status: ctr.Status,
			Ports:  ports,
		})
	}
	c.JSON(http.StatusOK, result)
}

// Helpers
func (h *ComputeHandler) getClaims(c *gin.Context) (*auth.Claims, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token, _ = c.Cookie("token")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	return auth.ValidateToken(token, h.Config)
}

func (h *ComputeHandler) hasPermission(claims *auth.Claims, resource string, perm int) bool {
	// Look up the user's role and check if it has the required permission
	role, err := h.IAMStore.GetRoleByID(claims.RoleID)
	if err != nil {
		return false
	}
	for _, p := range role.Permissions {
		if int(p) == perm {
			return true
		}
	}
	return false
}

// Docker helpers
func natPortSet(ports map[string]string) nat.PortSet {
	ps := nat.PortSet{}
	for cport := range ports {
		p, _ := nat.NewPort("tcp", cport)
		ps[p] = struct{}{}
	}
	return ps
}

func natPortBindings(ports map[string]string) nat.PortMap {
	pm := nat.PortMap{}
	for cport, hport := range ports {
		p, _ := nat.NewPort("tcp", cport)
		pm[p] = []nat.PortBinding{{HostPort: hport}}
	}
	return pm
}

func formatPort(port uint16, typ string) string {
	if port == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%s", port, typ)
}

func randString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
