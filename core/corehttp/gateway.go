package corehttp

import (
	"fmt"
	"net"
	"net/http"

	version "github.com/ipfs/go-ipfs"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"

	id "gx/ipfs/QmUDTcnDp2WssbmiDLC6aYurUeyt7QeRakHUQMxA2mZ5iB/go-libp2p/p2p/protocol/identify"
)

type GatewayConfig struct {
	Headers      map[string][]string
	Writable     bool
	PathPrefixes []string
}

func GatewayOption(writable bool, paths ...string) ServeOption {
	return func(n *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		cfg, err := n.Repo.Config()
		if err != nil {
			return nil, err
		}

		api := coreapi.NewCoreAPI(n)
		if cfg.Gateway.NoFetch {
			api = api.(*coreapi.CoreAPI).Offline()
		}

		gateway := newGatewayHandler(n, GatewayConfig{
			Headers:      cfg.Gateway.HTTPHeaders,
			Writable:     writable,
			PathPrefixes: cfg.Gateway.PathPrefixes,
		}, api)

		for _, p := range paths {
			mux.Handle(p+"/", gateway)
		}
		return mux, nil
	}
}

func VersionOption() ServeOption {
	return func(_ *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Commit: %s\n", version.CurrentCommit)
			fmt.Fprintf(w, "Client Version: %s\n", id.ClientVersion)
			fmt.Fprintf(w, "Protocol Version: %s\n", id.LibP2PVersion)
		})
		return mux, nil
	}
}
