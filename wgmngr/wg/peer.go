package wg

import (
	"bytes"
	"fmt"
	"os/exec"
)

func EnablePeer(iface, pubkey, pskey string) error {
	cmd := exec.Command("wg", "set", iface, "peer", pubkey, "preshared-key", pskey, "allowed-ips", "0.0.0.0/0,::/0")
	var stdout, stderr bytes.Buffer
	stdout.Grow(44)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set peer: %w %s %s", err, stderr.String(), stdout.String())
	}

	return nil
}

func DisablePeer(iface, pubkey string) error {
	cmd := exec.Command("wg", "set", iface, "peer", pubkey, "remove")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set peer: %w", err)
	}

	return nil
}
