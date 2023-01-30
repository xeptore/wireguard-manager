package wg

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/ini.v1"

	"github.com/xeptore/wireguard-manager/wgmngr/env"
)

type WGServerConfig struct {
	PublicKey      string
	PrivateKey     string
	ListenPort     string
	PublicHost     string
	ServerIPv4CIDR string
	ServerIPv6CIDR string
}

func LoadConfig(ctx context.Context, wgServerConfigFilePath string) (*WGServerConfig, error) {
	var out WGServerConfig
	conf, err := ini.LoadSources(
		ini.LoadOptions{
			AllowShadows:               true,
			AllowNonUniqueSections:     true,
			AllowDuplicateShadowValues: true,
		},
		wgServerConfigFilePath,
	)
	if nil != err {
		return nil, fmt.Errorf("failed to read wireguard server config file: %w", err)
	}
	serverAddr := conf.Section("Interface").Key("Address").String()
	ipPairs := strings.Split(serverAddr, ",")
	if len(ipPairs) != 2 {
		return nil, fmt.Errorf("expected 2 ip addresses in wireguard server config, got: %v", ipPairs)
	}
	for _, v := range ipPairs {
		v = strings.TrimSpace(v)
		ip, ipnet, err := net.ParseCIDR(v)
		if nil != err {
			return nil, fmt.Errorf("failed to parse wireguard server ip cidr address '%s': %w", v, err)
		}
		if ip.IsLoopback() || ip.IsUnspecified() || !ip.IsPrivate() || !ipnet.Contains(ip) {
			return nil, fmt.Errorf("expected wireguard server ip address '%s' to be a valid ip address", ip.String())
		}
		addr := netip.MustParseAddr(ip.String())
		if addr.Is4() {
			out.ServerIPv4CIDR = v
		} else if addr.Is6() {
			out.ServerIPv6CIDR = v
		} else {
			return nil, fmt.Errorf("expected wireguard server cidr address to either be a v6 or v4 ip address, got: %s", v)
		}
	}

	out.PublicHost = env.MustGet("WG_SERVER_HOST")
	out.ListenPort = env.MustGet("WG_SERVER_LISTEN_PORT")
	out.PrivateKey = env.MustGet("WG_SERVER_PRIVATE_KEY")
	privateKey, err := wgtypes.ParseKey(out.PrivateKey)
	if nil != err {
		return nil, fmt.Errorf("failed to parse wireguard server private key: %v", err)
	}

	out.PublicKey = privateKey.PublicKey().String()

	return &out, nil
}
