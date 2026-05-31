package sysreq

import (
	"fmt"
	"os/exec"
)

func RequireCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("missing required command: %s", name)
	}
	return nil
}
