package model

import (
	"binomeway.com/common/model"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/xhit/go-str2duration/v2"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ExecuteSystemCommand(command []string, spec model.ProcessInfo) (model.ProcessResult, error) {
	p := model.ProcessResult{}

	ctx, cancel, err := LaunchProcessWithLimits(spec)
	defer cancel()

	cmdResult := exec.CommandContext(ctx, command[0], command[1:]...)
	cmdResult.Dir = spec.WorkingDirectory
	cmdResult.Env = append(os.Environ(), ParseEnv(spec.EnvironmentVars)...)
	cmdResult.Stdin = strings.NewReader(spec.StandardInput)

	var stdout, stderr bytes.Buffer
	cmdResult.Stdout = &stdout
	cmdResult.Stderr = &stderr

	startTime := time.Now()
	err = cmdResult.Start()
	if err != nil {
		return model.ProcessResult{}, err
	}

	kbMem, err := getMemoryUsage(cmdResult)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Milliseconds()

	p.Stdout = stdout.String()
	p.Stderr = stderr.String()
	p.Output = p.Stdout + p.Stderr
	p.Time = int32(elapsedTime)
	p.Memory = kbMem

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			p.ExitCode = int32(exitErr.ExitCode())
		}
	}

	return p, err
}

// in bytes
func getMemoryUsage(cmd *exec.Cmd) (uint64, error) {
	// Get the process ID of the executed command
	pid := cmd.Process.Pid

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, err
	}

	// Get memory usage in bytes
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}

	err = cmd.Wait()
	if err != nil {
		return 0, err
	}

	return memInfo.RSS, err
}

func DefaultLimits() model.ProcessInfo {
	return model.ProcessInfo{
		CPUTime:        "50s",
		MaxOpenedFiles: 128,
		MaxProcesses:   128,
		Permissions: model.Permissions{
			CanRead:       true,
			CanWrite:      true,
			NetworkAccess: false,
		},
		EnvironmentVars:  make(map[string]string),
		WorkingDirectory: ".",
	}
}

func LaunchProcessWithLimits(lim model.ProcessInfo) (context.Context, context.CancelFunc, error) {
	if lim.CPUTime == "" {
		lim.CPUTime = DefaultLimits().CPUTime
	}

	duration, err := str2duration.ParseDuration(lim.CPUTime)
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	return ctx, cancel, err
}

func ParseEnv(env map[string]string) []string {
	var envs []string

	for key, value := range env {
		newEntry := fmt.Sprintf("%s=%s", key, value)
		envs = append(envs, newEntry)
	}

	return envs
}
