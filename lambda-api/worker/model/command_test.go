package model

import (
	"binomeway.com/common/model"
	"fmt"
	"testing"
)

func TestExecuteSystemCommand(t *testing.T) {
	var tests = []struct {
		command  []string
		exitCode int32
	}{
		{[]string{"ls", "-l"}, 0},
		{[]string{"echo", "test"}, 0},
		{[]string{"grep", "waits_for_stdin"}, 1},
		{[]string{"cat", "404"}, 1},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%d", tt.command, tt.exitCode)
		t.Run(testName, func(t *testing.T) {
			ans, _ := ExecuteSystemCommand(tt.command, model.ProcessInfo{})

			t.Logf("mem: %d", ans.Memory)

			if ans.ExitCode != tt.exitCode {
				t.Errorf("Bad exit code: %d, want %d, with output: %s", ans.ExitCode, tt.exitCode, ans.Output)
			}
		})
	}
}
