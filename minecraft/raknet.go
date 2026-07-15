package minecraft

import (
	"context"
	"github.com/sandertv/go-raknet"
	"log/slog"
	"net"
)

// UpstreamDialer opens the underlying transport connection for a Network. It is
// satisfied by *net.Dialer and by custom proxy dialers (e.g. SOCKS5 UDP). Set it
// on Dialer.UpstreamDialer to route the game connection through a proxy.
type UpstreamDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// upstreamDialerKey is the context key under which Dialer stashes its
// UpstreamDialer so a Network can pick it up.
type upstreamDialerKey struct{}

// upstreamDialerFrom returns the UpstreamDialer carried by ctx, or nil.
func upstreamDialerFrom(ctx context.Context) UpstreamDialer {
	d, _ := ctx.Value(upstreamDialerKey{}).(UpstreamDialer)
	return d
}

// RakNet is an implementation of a RakNet v10 Network.
type RakNet struct {
	l *slog.Logger
}

// DialContext ...
func (r RakNet) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return raknet.Dialer{
		ErrorLog:       r.l.With("net origin", "raknet"),
		UpstreamDialer: upstreamDialerFrom(ctx),
	}.DialContext(ctx, address)
}

// PingContext ...
func (r RakNet) PingContext(ctx context.Context, address string) (response []byte, err error) {
	return raknet.Dialer{
		ErrorLog:       r.l.With("net origin", "raknet"),
		UpstreamDialer: upstreamDialerFrom(ctx),
	}.PingContext(ctx, address)
}

// Listen ...
func (r RakNet) Listen(address string) (NetworkListener, error) {
	return raknet.ListenConfig{ErrorLog: r.l.With("net origin", "raknet")}.Listen(address)
}

// init registers the RakNet network.
func init() {
	RegisterNetwork("raknet", func(l *slog.Logger) Network { return RakNet{l: l} })
}
