package lego

import (
	"net"
	"strings"
	"time"

	"github.com/bryant-rh/certcli/pkg/global"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/providers/http/webroot"
	"k8s.io/klog/v2"
)

func SetupChallenges(ctx *global.CertcliOptions, client *lego.Client) {
	if !ctx.Http && !ctx.Tls && ctx.Dns == "" {
		klog.Fatal("No challenge selected. You must specify at least one challenge: `--http`, `--tls`, `--dns`.")
	}

	if ctx.Http {
		err := client.Challenge.SetHTTP01Provider(setupHTTPProvider(ctx))
		if err != nil {
			klog.Fatal(err)
		}
	}

	if ctx.Tls {
		err := client.Challenge.SetTLSALPN01Provider(setupTLSProvider(ctx))
		if err != nil {
			klog.Fatal(err)
		}
	}

	if ctx.Dns != "" {
		setupDNS(ctx, client)
	}
}

func setupHTTPProvider(ctx *global.CertcliOptions) challenge.Provider {
	switch {
	case ctx.Http_Webroot != "":
		ps, err := webroot.NewHTTPProvider(ctx.Http_Webroot)
		if err != nil {
			klog.Fatal(err)
		}
		return ps
	// case ctx.IsSet("http.memcached-host"):
	// 	ps, err := memcached.NewMemcachedProvider(ctx.StringSlice("http.memcached-host"))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return ps
	case ctx.Http_port != "":
		iface := ctx.Http_port
		if !strings.Contains(iface, ":") {
			klog.Fatalf("The --http switch only accepts interface:port or :port for its argument.")
		}

		host, port, err := net.SplitHostPort(iface)
		if err != nil {
			klog.Fatal(err)
		}

		srv := http01.NewProviderServer(host, port)
		if header := ctx.Http_ProxyHeader; header != "" {
			srv.SetProxyHeader(header)
		}
		return srv
	case ctx.Http:
		srv := http01.NewProviderServer("", "")
		if header := ctx.Http_ProxyHeader; header != "" {
			srv.SetProxyHeader(header)
		}
		return srv
	default:
		klog.Fatal("Invalid HTTP challenge options.")
		return nil
	}
}

func setupTLSProvider(ctx *global.CertcliOptions) challenge.Provider {
	switch {
	case ctx.Tls_port != "":
		iface := ctx.Tls_port
		if !strings.Contains(iface, ":") {
			klog.Fatalf("The --tls switch only accepts interface:port or :port for its argument.")
		}

		host, port, err := net.SplitHostPort(iface)
		if err != nil {
			klog.Fatal(err)
		}

		return tlsalpn01.NewProviderServer(host, port)
	case ctx.Tls:
		return tlsalpn01.NewProviderServer("", "")
	default:
		klog.Fatal("Invalid HTTP challenge options.")
		return nil
	}
}

func setupDNS(ctx *global.CertcliOptions, client *lego.Client) {
	provider, err := dns.NewDNSChallengeProviderByName(ctx.Dns)
	if err != nil {
		klog.Fatal(err)
	}

	servers := ctx.Dns_resolvers
	err = client.Challenge.SetDNS01Provider(provider,
		dns01.CondOption(len(servers) > 0,
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(ctx.Dns_resolvers))),
		dns01.CondOption(ctx.Dns_Disablecp,
			dns01.DisableCompletePropagationRequirement()),
		dns01.CondOption(ctx.Dns_timeout != 0,
			dns01.AddDNSTimeout(time.Duration(ctx.Dns_timeout)*time.Second)),
	)
	if err != nil {
		klog.Fatal(err)
	}
}
