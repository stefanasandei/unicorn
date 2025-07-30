package model

import (
	"binomeway.com/common/model"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
)

type Worker struct {
	ID          uuid.UUID        `json:"id"`
	CPUUsage    float64          `json:"cpu_usage"`
	LastUpdated uint64           `json:"last_updated"`
	Languages   []model.Language `json:"-"`
}

func NewWorker(id uuid.UUID, runtimesDir string) Worker {
	worker := Worker{
		ID: id,
	}

	languages, err := QuerySupportedLanguages(runtimesDir)
	if err != nil {
		log.Panicf("Failed to build worker: %s", err)
	}

	worker.Languages = languages

	worker.Update()

	return worker
}

func (w *Worker) Update() {
	w.SetCPUUsage()
	w.LastUpdated = uint64(time.Now().UnixMilli())
}

func (w *Worker) SetCPUUsage() {
	// wait a few seconds to query CPU usage between two intervals in time
	// this makes the worker startup slower
	interval := 2 * time.Second

	// Get CPU usage measurements
	percent, err := cpu.Percent(interval, false)
	if err != nil {
		log.Panicf("Failed to read CPU usage: %s", err)
	}

	// Calculate the average CPU usage, since Percent() returns data per CPU core
	totalUsage := 0.0
	for _, p := range percent {
		totalUsage += p
	}

	w.CPUUsage = totalUsage / float64(len(percent))
}

func QuerySupportedLanguages(runtimesDir string) ([]model.Language, error) {
	var languages []model.Language

	if _, err := os.Stat(runtimesDir); os.IsNotExist(err) {
		return []model.Language{}, errors.New("runtimes directory is not present")
	}

	entries, err := os.ReadDir(runtimesDir)
	if err != nil {
		return []model.Language{}, err
	}

	for _, entry := range entries {
		if info, err := entry.Info(); err == nil {
			if !strings.HasSuffix(info.Name(), ".yaml") {
				continue
			}
		}

		newLang, err := model.NewLanguage(entry, runtimesDir)
		if err != nil {
			return []model.Language{}, err
		}

		languages = append(languages, newLang)
	}

	return languages, nil
}
