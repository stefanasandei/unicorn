package model

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	var tests = []struct {
		runtimesDir string
		shouldWork  bool
	}{
		{"../../runtimes", true},
		{"./runtimes", false},
	}

	for _, tt := range tests {
		t.Run(tt.runtimesDir, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil && tt.shouldWork {
					t.Errorf("Failed to create a new worker: %s", err)
				}
			}()

			worker := NewWorker(uuid.New(), tt.runtimesDir)

			assert.Greater(t, worker.LastUpdated, uint64(0), "The worker should have a valid timestamp")
			assert.Less(t, worker.LastUpdated, uint64(time.Now().UnixMilli()+1), "The worker should have been created in the past")

			assert.Greater(t, len(worker.Languages), 1, "The worker should have at least one language")

			assert.Greaterf(t, worker.CPUUsage, 0.0, "The CPU usage should be greater than 0.0f")
		})
	}
}
