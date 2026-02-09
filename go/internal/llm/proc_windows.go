//go:build windows

package llm

import (
	"os/exec"
)

func setProcGroup(cmd *exec.Cmd) {
	// No process group management on Windows
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}
