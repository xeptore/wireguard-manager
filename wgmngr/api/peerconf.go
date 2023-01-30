package api

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/ini.v1"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/table"
)

var (
	ErrNoMoreIPv4Address = errors.New("no more ipv4 address is available in the range")
	ErrNoMoreIPv6Address = errors.New("no more ipv6 address is available in the range")
)

type CreatePeerConfigReq struct {
	Name        string
	Description string
	ResellerID  string
}

func getLastPeerIPs(ctx context.Context, db *sql.DB) (ipv4 string, ipv6 string, err error) {
	var c m.PeerConfigs
	err = t.PeerConfigs.SELECT(t.PeerConfigs.Ipv4, t.PeerConfigs.Ipv6).LIMIT(1).ORDER_BY(t.PeerConfigs.GeneratedAt.DESC()).QueryContext(ctx, db, &c)
	return c.Ipv4, c.Ipv6, err
}

func invalidIPv4AddressErr(addr string) error {
	return fmt.Errorf("%w: %s", ErrInvalidIPv4Address, addr)
}

func invalidIPv6AddressErr(addr string) error {
	return fmt.Errorf("%w: %s", ErrInvalidIPv6Address, addr)
}

var (
	ErrInvalidIPv4Address = errors.New("invalid ipv4 address")
	ErrInvalidIPv6Address = errors.New("invalid ipv6 address")
)

func findNextIpv4(srcIPv4 string, serverIPv4Net *net.IPNet) (string, error) {
	ip, err := netip.ParseAddr(srcIPv4)
	if nil != err || !ip.Is4() || !ip.IsValid() || !ip.IsPrivate() || ip.IsLoopback() || ip.IsMulticast() || ip.IsUnspecified() {
		return "", invalidIPv4AddressErr(srcIPv4)
	}

	nextIP := ip.Next()
	if !nextIP.IsValid() || !nextIP.IsPrivate() || nextIP.IsUnspecified() || !serverIPv4Net.Contains(nextIP.AsSlice()) {
		return "", ErrNoMoreIPv4Address
	}

	nextIPstr := nextIP.String()
	if nextIPstr == "invalid IP" {
		return "", ErrNoMoreIPv4Address
	}

	return nextIPstr, nil
}

func findNextIpv6(srcIPv6 string, serverIPv6Net *net.IPNet) (string, error) {
	ip, err := netip.ParseAddr(srcIPv6)
	if nil != err || !ip.Is6() || !ip.IsValid() || !ip.IsPrivate() || ip.IsLoopback() || ip.IsMulticast() || ip.IsUnspecified() {
		return "", invalidIPv6AddressErr(srcIPv6)
	}

	nextIP := ip.Next()
	if !nextIP.IsValid() || !nextIP.IsPrivate() || nextIP.IsUnspecified() || !serverIPv6Net.Contains(nextIP.AsSlice()) {
		return "", ErrNoMoreIPv6Address
	}

	nextIPstr := nextIP.StringExpanded()
	if nextIPstr == "invalid IP" {
		return "", ErrNoMoreIPv6Address
	}

	return nextIPstr, nil
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
	}

	serverIPv4, serverIPv4Net, err := net.ParseCIDR(h.wgServerConf.ServerIPv4CIDR)
	if nil != err {
		return "", err
	}
	if lastIpv4 == "" {
		lastIpv4 = serverIPv4.String()
	}

	ipv4, err := findNextIpv4(lastIpv4, serverIPv4Net)
	if nil != err {
		return "", err
	}

	serverIPv6, serverIPv6Net, err := net.ParseCIDR(h.wgServerConf.ServerIPv6CIDR)
	if nil != err {
		return "", err
	}
	if lastIpv6 == "" {
		lastIpv6 = serverIPv6.String()
	}

	ipv6, err := findNextIpv6(lastIpv6, serverIPv6Net)
	if nil != err {
		return "", err
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", err
	}

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return "", err
	}

	c := m.PeerConfigs{
		ID:            id,
		GeneratedAt:   time.Now(),
		Name:          req.Name,
		Description:   req.Description,
		GeneratedByID: req.ResellerID,
		Ipv4:          ipv4,
		Ipv6:          ipv6,
		PrivateKey:    key.String(),
		PublicKey:     key.PublicKey().String(),
		PresharedKey:  presharedKey.String(),
		IsActive:      true,
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

	// TODO: Signal wireguard to reload config

	return c.ID, nil
}

type GetPeerConfigReq struct {
	ResellerID string
	ConfigID   string
}

func (h *Handler) GetPeerConfig(ctx context.Context, req GetPeerConfigReq) ([]byte, error) {
	var c m.PeerConfigs
	err := t.PeerConfigs.SELECT(t.PeerConfigs.Ipv4, t.PeerConfigs.Ipv6, t.PeerConfigs.PresharedKey, t.PeerConfigs.PrivateKey).WHERE(t.PeerConfigs.ID.EQ(mysql.String(req.ConfigID)).AND(t.PeerConfigs.GeneratedByID.EQ(mysql.String(req.ResellerID)))).LIMIT(1).QueryContext(ctx, h.db, &c)
	if nil != err {
		return nil, err
	}

	cfg := ini.Empty()
	ifaceSec := cfg.Section("Interface")
	ifaceSec.NewKey("PrivateKey", c.PrivateKey)
	ifaceSec.NewKey("Address", fmt.Sprintf("%s/32, %s/128", c.Ipv4, c.Ipv6))
	ifaceSec.NewKey("DNS", "1.1.1.1, 1.0.0.1")
	peerSec := cfg.Section("Peer")
	peerSec.NewKey("PublicKey", h.wgServerConf.PublicKey)
	peerSec.NewKey("PresharedKey", c.PresharedKey)
	peerSec.NewKey("Endpoint", fmt.Sprintf("%s:%s", h.wgServerConf.PublicHost, h.wgServerConf.ListenPort))
	peerSec.NewKey("AllowedIPs", "0.0.0.0/0, ::/0")

	var out bytes.Buffer
	if _, err := cfg.WriteTo(&out); nil != err {
		return nil, fmt.Errorf("failed to serialize ini config: %w", err)
	}

	return out.Bytes(), nil
}
