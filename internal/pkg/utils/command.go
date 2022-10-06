package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func Exec(cmd string, args ...string) (out bytes.Buffer, err error) {
	var stderr bytes.Buffer

	c := exec.Command(cmd, args...)
	c.Stdout = &out
	c.Stderr = &stderr

	log.Printf("run cmd: [%s]", c.String())
	if err := c.Run(); err != nil {
		return out, fmt.Errorf(fmt.Sprintf("%v%v", stderr.String(), err))
	}

	return out, nil
}
