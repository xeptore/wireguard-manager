package api

import (
	"bytes"
	"context"
	"fmt"

	m "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/model"
	t "github.com/xeptore/wireguard-manager/wgmngr/db/gen/latest/wgmngr/table"
	"gopkg.in/ini.v1"
)

func (h *Handler) GetActivePeerConfigs(ctx context.Context) ([]byte, error) {
	rows, err := t.PeerConfigs.
		SELECT(t.PeerConfigs.PublicKey, t.PeerConfigs.PresharedKey, t.PeerConfigs.Ipv4, t.PeerConfigs.Ipv6, t.PeerConfigs.ID).
		WHERE(t.PeerConfigs.IsActive.IS_TRUE()).
		LIMIT(1_000).
		ORDER_BY(t.PeerConfigs.GeneratedAt.ASC()).
		Rows(ctx, h.db)
	if nil != err {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); nil != closeErr {
			if nil != err {
				err = fmt.Errorf("closing rows failed while there was already an error: %v %w", err, closeErr)
				return
			}

			err = closeErr
		}
	}()

	var configs []m.PeerConfigs
	for rows.Next() {
		var cfg m.PeerConfigs
		if err := rows.Scan(&cfg); nil != err {
			return nil, err
		}
		configs = append(configs, cfg)
	}

	var out bytes.Buffer
	cfg := ini.Empty()
	peerSec := cfg.Section("Interface")
	peerSec.NewKey("PrivateKey", h.wgServerConf.PrivateKey)
	peerSec.NewKey("ListenPort", h.wgServerConf.ListenPort)
	if _, err := cfg.WriteTo(&out); nil != err {
		return nil, fmt.Errorf("failed to serialize wireguard server interface section config: %w", err)
	}

	for i, c := range configs {
		cfg := ini.Empty()
		peerSec := cfg.Section("Peer")
		peerSec.NewKey("PublicKey", c.PublicKey)
		peerSec.NewKey("PresharedKey", c.PresharedKey)
		peerSec.NewKey("AllowedIPs", fmt.Sprintf("%s/32, %s/128", c.Ipv4, c.Ipv6))

		if _, err := cfg.WriteTo(&out); nil != err {
			return nil, fmt.Errorf("failed to serialize wireguard server peer %d (%s) config: %w", i, c.ID, err)
		}
	}

	return out.Bytes(), nil
}
