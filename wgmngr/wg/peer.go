package wg

import (
	"fmt"
	"os/exec"
)

func EnablePeer(iface, pubkey, pskey string) error {
	cmd := exec.Command("wg", "set", iface, "peer", pubkey, "preshared-key", pskey, "allowed-ips", "0.0.0.0/0,::/0")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set peer: %w", err)
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
