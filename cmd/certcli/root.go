package cmd

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bryant-rh/certcli/pkg/certxctx"
	"github.com/bryant-rh/certcli/pkg/global"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/xenolf/lego/lego"
	"k8s.io/klog/v2"
	"software.sslmate.com/src/go-pkcs12"
)

var (
	defaultPath string
	version     string
	domainSlice [][]string

	certx       certxctx.CertClient
	domainsList []string

	resourceTypeSlice = []string{"clb", "cdn"}
)

// versionString returns the version prefixed by 'v'
// or an empty string if no version has been populated by goreleaser.
// In this case, the --version flag will not be added by cobra.
func versionString() string {
	if len(version) == 0 {
		return ""
	}
	return "v" + version
}

func NewCmd(o *global.CertcliOptions) *cobra.Command {

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     "certcli",
		Short:   "certcli 是一个基于Let's Encrypt 命令行https 证书申请工具 ",
		Version: versionString(),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// 	err := Validate(cmd, args)
		// 	if err != nil {
		// 		klog.Fatal(util.RedColor(err))
		// 	}
		// 	// Client = pkg.NewReqClient()
		// 	// if enableDebug { // Enable debug mode if `--enableDebug=true` or `DEBUG=true`.
		// 	// 	Client.SetDebug(true)
		// 	// }

		// },
	}
	rootCmd.PersistentFlags().BoolVar(&o.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&o.Eab, "eab", o.Eab, "Use External Account Binding for account registration. Requires --kid and --hmac.")
	rootCmd.PersistentFlags().BoolVar(&o.Tls, "tls", o.Tls, "Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.")
	rootCmd.PersistentFlags().BoolVar(&o.Http, "http", o.Http, "Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.")
	rootCmd.PersistentFlags().BoolVar(&o.Dns_Disablecp, "dns.disable-cp", o.Dns_Disablecp, "By setting this flag to true, disables the need to await propagation of the TXT record to all authoritative name servers.")
	rootCmd.PersistentFlags().BoolVar(&o.Pem, "pem", o.Pem, "Generate an additional .pem (base64) file by concatenating the .key and .crt files together.")
	rootCmd.PersistentFlags().BoolVar(&o.Pfx, "pfx", o.Pfx, "Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.")
	rootCmd.PersistentFlags().BoolVarP(&o.Accept_tos, "accept-tos", "y", o.Accept_tos, "By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.")

	rootCmd.PersistentFlags().StringSliceVarP(&o.Domains, "domains", "d", o.Domains, "指定域名,可指定多个以逗号分割(Specify the domain name, you can specify multiple separated by commas)")
	rootCmd.PersistentFlags().StringSliceVarP(&o.Dns_resolvers, "dns.resolvers", "", o.Dns_resolvers, "Set the resolvers to use for performing (recursive) CNAME resolving and apex domain determination.")

	rootCmd.PersistentFlags().StringVar(&o.Dns, "dns", "", "Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges. Run 'certcli dnshelp' for help on usage.")
	rootCmd.PersistentFlags().StringVarP(&o.Email, "email", "m", "", "Email used for registration and recovery contact.")
	rootCmd.PersistentFlags().StringVarP(&o.Key_type, "key_type", "k", "rsa2048", "Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384.")
	rootCmd.PersistentFlags().StringVarP(&o.Acme_server, "server", "s", lego.LEDirectoryProduction, "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.")
	rootCmd.PersistentFlags().StringVarP(&o.User_agent, "user-agent", "", "", "Add to the user-agent sent to the CA to identify an application embedding lego-cli")
	rootCmd.PersistentFlags().StringVar(&o.Csr, "csr", "", "Certificate signing request filename, if an external CSR is to be used.")
	rootCmd.PersistentFlags().StringVar(&o.Kid, "kid", "", "Key identifier from External CA. Used for External Account Binding.")
	rootCmd.PersistentFlags().StringVar(&o.Hmac, "hmac", "", "MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.")
	rootCmd.PersistentFlags().StringVar(&o.Http_port, "http.port", ":80", "Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port.")
	rootCmd.PersistentFlags().StringVar(&o.Http_ProxyHeader, "http.proxy-header", "host", "Validate against this HTTP header when solving HTTP-01 based challenges behind a reverse proxy.")
	rootCmd.PersistentFlags().StringVar(&o.Http_Webroot, "http.webroot", "", "Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file.")
	rootCmd.PersistentFlags().StringVar(&o.Tls_port, "tls.port", ":443", "Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port.")
	rootCmd.PersistentFlags().StringVar(&o.Pfx_Pass, "pfx-pass", pkcs12.DefaultPassword, "The password used to encrypt the .pfx (PCKS#12) file.")

	rootCmd.PersistentFlags().IntVarP(&o.Cert_timeout, "cert.timeout", "", 30, "Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates.")
	rootCmd.PersistentFlags().IntVarP(&o.Dns_timeout, "dns-timeout", "", 10, "Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries.")
	rootCmd.PersistentFlags().IntVarP(&o.Http_timeout, "http-timeout", "", o.Http_timeout, "Set the HTTP timeout value to a specific value in seconds.")

	rootCmd.AddCommand(NewCmdRun(o))
	rootCmd.AddCommand(NewCmdSync(o))
	rootCmd.AddCommand(NewCmdUpload(o))
	rootCmd.AddCommand(NewDNSHelp(o))

	return rootCmd
}

func initLog() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{}) // convince pkg/flag we parsed the flags

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Name != "v" { // hide all glog flags except for -v
			pflag.Lookup(f.Name).Hidden = true
		}
	})
	if err := flag.Set("logtostderr", "true"); err != nil {
		fmt.Printf("can't set log to stderr %+v", err)
		os.Exit(1)
	}
}

func init() {
	initLog()
}

func domainsToSlice(d []string) [][]string {
	for _, i := range d {
		var data []string
		data = append(data, i)
		domainSlice = append(domainSlice, data)

	}
	return domainSlice
}
