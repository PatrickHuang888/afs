package daemon

import (
	"context"
	"fmt"
	"github.com/PatrickHuang888/afs/logging"
	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	nilrouting "github.com/ipfs/go-ipfs-routing/none"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	bsnet "github.com/ipfs/go-bitswap/network"
)

type Host struct {
	ip   string
	port string
	h    host.Host
}

type Daemon struct {
	ctx context.Context
	host Host
	exchange exchange.Interface
}

func New(ctx context.Context, ip, port string) (*Daemon, error) {
	d := &Daemon{ctx: ctx}

	h := Host{ip: ip, port: port}

	a := fmt.Sprintf("/ip4/%s/tcp/%s", h.ip, h.port)
	//a := fmt.Sprintf("/ip4/%s/udp/%s/quic", h.ip, h.port)
	addr, err := ma.NewMultiaddr(a)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	h.h, err = libp2p.New(ctx, libp2p.ListenAddrs(addr))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logging.Infof("listening on address %s ...", addr)

	routing, err := nilrouting.ConstructNilRouting(ctx, nil, nil, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	net := bsnet.NewFromIpfsHost(h.h, routing)
	bstore:= blockstore.NewBlockstore(datastore.NewMapDatastore())
	exchange := bitswap.New(ctx, net, bstore)

	d.host= h
	d.exchange= exchange
	return d, nil
}

func (d *Daemon) Shutdown() {
}

func (d *Daemon) Run() error {
	<-d.ctx.Done()

	logging.Info("shutting down daemonn")
	d.Shutdown()
	logging.Info("daemon stopped")
	return nil
}
