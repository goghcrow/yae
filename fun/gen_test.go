package fun

import (
	"os/exec"
	"testing"
)

func TestGenFun(t *testing.T) {
	cmd := exec.Command("/bin/bash", "gen.sh")
	_, _ = cmd.Output()
}
