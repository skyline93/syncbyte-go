package utils

import (
	"bytes"
	"log"
	"os/exec"
)

func Exec(cmd string, args ...string) (out bytes.Buffer, err error) {
	c := exec.Command(cmd, args...)
	c.Stdout = &out

	log.Printf("run cmd: [%s]", c.String())
	if err := c.Run(); err != nil {
		return out, err
	}

	return out, nil
}
