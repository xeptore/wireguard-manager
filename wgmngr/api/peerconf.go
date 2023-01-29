package api

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gopkg.in/ini.v1"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/table"
	"github.com/xeptore/wireguard-manager/wgmngr/wg"
)

var (
	ErrNoMoreIPv4Address = errors.New("no more ipv4 address is available in the range")
	ErrNoMoreIPv6Address = errors.New("no more ipv6 address is available in the range")
)

type CreatePeerConfigReq struct {
	Name        string
	Description string
	ResellerID  string
	ServerIPv4  string
	ServerIPv6  string
}

func getLastPeerIPs(ctx context.Context, db *sql.DB) (string, string, error) {
	var c m.PeerConfigs
	err := t.PeerConfigs.SELECT(t.PeerConfigs.Ipv4, t.PeerConfigs.Ipv6).LIMIT(1).ORDER_BY(t.PeerConfigs.GeneratedAt.DESC()).QueryContext(ctx, db, &c)
	return c.Ipv4, c.Ipv6, err
}

func findNextIpv4(src string) (string, bool) {
	return "", false
}

func findNextIpv6(src string) (string, bool) {
	return "", false
}

func (h *Handler) CreatePeerConfig(ctx context.Context, req CreatePeerConfigReq) (string, error) {
	id, err := gonanoid.New(64)
	if nil != err {
		return "", fmt.Errorf("failed to generate id for the peer config: %w", err)
	}

	lastIpv4, lastIpv6, err := getLastPeerIPs(ctx, h.db)
	if nil != err {
		if !errors.Is(err, qrm.ErrNoRows) {
			return "", fmt.Errorf("failed to get last peer ip addresses: %v", err)
		}

		lastIpv4 = req.ServerIPv4
		lastIpv6 = req.ServerIPv6
	}

	ipv4, ok := findNextIpv4(lastIpv4)
	if !ok {
		return "", ErrNoMoreIPv4Address
	}

	ipv6, ok := findNextIpv6(lastIpv6)
	if !ok {
		return "", ErrNoMoreIPv6Address
	}

	cmd := exec.Command("wg", "genkey")
	var stdout, stderr bytes.Buffer
	stdout.Grow(44)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate peer private key: %v", err)
	}
	if len(stderr.Bytes()) != 0 {
		return "", errors.New("expected peer private key generation to not to output anything on stderr")
	}
	privateKey := stdout.String()

	cmd = exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBufferString(privateKey)
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate peer public key: %v", err)
	}
	if len(stderr.Bytes()) != 0 {
		return "", errors.New("expected peer public key generation to not to output anything on stderr")
	}
	publicKey := stdout.String()

	cmd = exec.Command("wg", "genpsk")
	cmd.Stdin = bytes.NewBufferString(privateKey)
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate peer preshared key: %v", err)
	}
	if len(stderr.Bytes()) != 0 {
		return "", errors.New("expected peer preshared key generation to not to output anything on stderr")
	}
	presharedKey := stdout.String()

	c := m.PeerConfigs{
		ID:            id,
		GeneratedAt:   time.Now(),
		Name:          req.Name,
		Description:   req.Description,
		GeneratedByID: req.ResellerID,
		Ipv4:          ipv4,
		Ipv6:          ipv6,
		PrivateKey:    privateKey,
		PresharedKey:  presharedKey,
	}
	res, err := t.PeerConfigs.INSERT(t.PeerConfigs.AllColumns).MODEL(c).ExecContext(ctx, h.db)
	if nil != err {
		return "", fmt.Errorf("failed to insert peer config: %v", err)
	}
	rows, err := res.RowsAffected()
	if nil != err {
		return "", fmt.Errorf("failed to get number of inserted rows: %v", err)
	}
	if rows != 1 {
		return "", fmt.Errorf("expected 1 row to be inserted, got %d", rows)
	}

	if err := wg.EnablePeer("wg0", publicKey, presharedKey); nil != err {
		return "", fmt.Errorf("failed to enable peer: %v", err)
	}

	return c.ID, nil
}

type GetPeerConfigReq struct {
	ResellerID string
	ConfigID   string
}

type WGConfig struct {
	InterfacePrivateKey string
	InterfaceAddress    string
	InterfaceDNS        string
	PeerPublicKey       string
	PeerPresharedKey    string
	PeerEndpoint        string
	PeerAllowedIPs      string
}

func (h *Handler) GetPeerConfig(ctx context.Context, req GetPeerConfigReq) ([]byte, error) {
	var c m.PeerConfigs
	err := t.PeerConfigs.SELECT(t.PeerConfigs.Ipv4, t.PeerConfigs.Ipv6, t.PeerConfigs.PresharedKey, t.PeerConfigs.PrivateKey).WHERE(t.PeerConfigs.ID.EQ(mysql.String(req.ConfigID)).AND(t.PeerConfigs.GeneratedByID.EQ(mysql.String(req.ResellerID)))).LIMIT(1).QueryContext(ctx, h.db, &c)
	if nil != err {
		return nil, err
	}

	/*
	   [Interface]
	   PrivateKey = SKeW6PUVxfry4gm3u7zZl4xRzJagrpzWXD2r3D4K+l4=
	   Address = 10.66.66.64/32, fd4f:fb5c:33d6:36d4::40/128
	   DNS = 1.1.1.1, 1.0.0.1

	   [Peer]
	   PublicKey = +jbF/7wsCxZLeaJ4ABV3tGiL6mkYaLDL1rZZJy7L9T4=
	   PresharedKey = LMXhQh+mBRbGPG6iICDntJxrrnfnqZNCTCVN8RI4qu8=
	   Endpoint = db.conisma.com:16438
	   AllowedIPs = 0.0.0.0/0, ::/0
	*/
	cfg := ini.Empty()
	ifaceSec := cfg.Section("Interface")
	ifaceSec.NewKey("PrivateKey", c.PrivateKey)
	ifaceSec.NewKey("Address", c.PrivateKey)
	ifaceSec.NewKey("DNS", c.PrivateKey)

	peerSec := cfg.Section("Peer")
	peerSec.NewKey("PublicKey", c.PresharedKey)
	peerSec.NewKey("PresharedKey", c.PresharedKey)
	peerSec.NewKey("Endpoint", c.PresharedKey)
	peerSec.NewKey("AllowedIPs", c.PresharedKey)

	var out bytes.Buffer
	if _, err := cfg.WriteTo(&out); nil != err {
		return nil, fmt.Errorf("failed to write ini config file: %w", err)
	}

	return out.Bytes(), nil
}
