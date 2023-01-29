package wg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Pubkey(privateKey string) (string, error) {
	cmd := exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBufferString(privateKey)
	var stdout, stderr bytes.Buffer
	stdout.Grow(44)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate peer public key: %v", err)
	}
	if len(stderr.Bytes()) != 0 {
		return "", errors.New("expected peer public key generation to not to output anything on stderr")
	}

	return strings.TrimSpace(stdout.String()), nil
}
