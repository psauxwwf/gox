package socks

import (
	"fmt"
	"log/slog"
	"maps"

	"gox/internal/server/socks/ruler"

	"github.com/things-go/go-socks5"
)

type Socks struct {
	server *socks5.Server
	listen string
}

func New(
	_listen string,
	creds map[string]string,
) *Socks {
	opt := []socks5.Option{
		socks5.WithLogger(slogLogger{}),
		socks5.WithRule(ruler.New()),
		socks5.WithResolver(socks5.DNSResolver{}),
	}
	if len(creds) != 0 {
		opt = append(opt,
			socks5.WithAuthMethods(
				[]socks5.Authenticator{
					socks5.UserPassAuthenticator{
						Credentials: toCreds(creds),
					},
				},
			),
		)
	}
	return &Socks{
		listen: _listen,
		server: socks5.NewServer(
			opt...,
		),
	}
}

func (s *Socks) Listen() error {
	slog.Info("listen socks", "listen", s.listen)
	return s.server.ListenAndServe(
		"tcp",
		s.listen,
	)
}

type slogLogger struct{}

func (slogLogger) Errorf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}

func toCreds(creds map[string]string) socks5.StaticCredentials {
	var staticCreds = make(socks5.StaticCredentials)
	maps.Copy(staticCreds, creds)
	return staticCreds
}

// type SystemResolver struct{}

// func (r *SystemResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
// 	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", name)
// 	if err != nil {
// 		return ctx, nil, err
// 	}
// 	return ctx, ips[0], nil
// }

// func (r *SystemResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
// 	resolver := &net.Resolver{
// 		PreferGo: true,
// 		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
// 			return net.Dial("udp", "127.0.0.11:53")
// 		},
// 	}
// 	ips, err := resolver.LookupIP(ctx, "ip", name)
// 	if err != nil {
// 		return ctx, nil, err
// 	}
// 	return ctx, ips[0], nil
// }
