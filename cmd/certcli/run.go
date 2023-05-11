package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bryant-rh/certcli/pkg/global"
	legox "github.com/bryant-rh/certcli/pkg/lego"
	"github.com/bryant-rh/certcli/pkg/util"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

const (
	rootPathWarningMessage = `!!!! HEADS UP !!!!

Your account credentials have been saved in your Let's Encrypt
configuration directory at "%s".

You should make a secure backup of this folder now. This
configuration directory will also contain certificates and
private keys obtained from Let's Encrypt so making regular
backups of this folder is ideal.
`
	renewEnvAccountEmail = "LEGO_ACCOUNT_EMAIL"
	renewEnvCertDomain   = "LEGO_CERT_DOMAIN"
	renewEnvCertPath     = "LEGO_CERT_PATH"
	renewEnvCertKeyPath  = "LEGO_CERT_KEY_PATH"
	renewEnvCertPEMPath  = "LEGO_CERT_PEM_PATH"
	renewEnvCertPFXPath  = "LEGO_CERT_PFX_PATH"
)

func NewCmdRun(o *global.CertcliOptions) *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "生成基于Let's Encrypt 颁发的ssl 证书(Generate ssl certificates based on Let's Encrypt)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			//校验参数
			err := validate(cmd, args, o)
			if err != nil {
				klog.Fatal(util.RedColor(err))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, d := range domainsToSlice(o.Domains) {
				message := fmt.Sprintf("Run AccountsStorage For Domain: [%s]\n", d[0])

				klog.Infof(util.GreenColor(message))
				err := run(o, d)
				if err != nil {
					klog.Fatal(err)
					//return err
				}

			}
			return nil

		},
	}

	cwd, err := os.Getwd()
	if err == nil {
		defaultPath = filepath.Join(cwd, ".certcli")
	}
	//runCmd.Flags().StringVarP(&domains, "domains", "d", "", "指定域名(Add a domain to the process)")
	runCmd.Flags().BoolVar(&o.No_bundle, "no-bundle", o.No_bundle, "Do not create a certificate bundle by adding the issuers certificate to the new certificate.")
	runCmd.Flags().BoolVar(&o.Must_staple, "must-staple", o.Must_staple, "Include the OCSP must staple TLS extension in the CSR and generated certificate.")
	runCmd.Flags().BoolVar(&o.Always_Deactivate_authorizations, "always-deactivate-authorizations", o.Always_Deactivate_authorizations, "Force the authorizations to be relinquished even if the certificate request was successful.")
	runCmd.Flags().StringVarP(&o.Path, "path", "", defaultPath, "Directory to use for storing the data.")
	runCmd.Flags().StringVar(&o.Run_hook, "run-hook", "", "Define a hook. The hook is executed when the certificates are effectively created.")
	runCmd.Flags().StringVar(&o.Preferred_Chain, "preferred-chain", "", "If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name.")

	return runCmd
}
func validate(cmd *cobra.Command, args []string, o *global.CertcliOptions) error {

	// if botKey == "" {
	// 	return fmt.Errorf("环境变量BOT_KEY为空:'%s',请设置", botKey)
	// }
	//fmt.Println(cmd.Flags().GetString("domains"))

	if o.Path == "" {
		klog.Fatal("Could not determine current working directory. Please pass --path.")
	}

	err := legox.CreateNonExistingFolder(o.Path)
	if err != nil {
		klog.Fatalf("Could not check/create path: %v", err)
	}

	if o.Acme_server == "" {
		klog.Fatal("Could not determine current working server. Please pass --server.")
	}

	// we require either domains or csr, but not both
	hasDomains := len(o.Domains) > 0
	hasCsr := len(o.Csr) > 0
	if hasDomains && hasCsr {
		klog.Fatal("Please specify either --domains/-d or --csr/-c, but not both")
	}
	if !hasDomains && !hasCsr {
		klog.Fatal("Please specify --domains/-d (or --csr/-c if you already have a CSR)")
	}
	return nil

}

func run(ctx *global.CertcliOptions, domains []string) error {
	accountsStorage := legox.NewAccountsStorage(ctx, domains)

	account, client := legox.Setup(ctx, accountsStorage)
	legox.SetupChallenges(ctx, client)

	if account.Registration == nil {
		reg, err := register(ctx, client)
		if err != nil {
			klog.Fatalf("Could not complete registration\n\t%v", err)
		}

		account.Registration = reg
		if err = accountsStorage.Save(account); err != nil {
			klog.Fatal(err)
		}

		fmt.Printf(rootPathWarningMessage, accountsStorage.GetRootPath())
	}

	certsStorage := legox.NewCertificatesStorage(ctx, domains)
	certsStorage.CreateRootFolder()

	cert, err := obtainCertificate(ctx, client, domains)
	if err != nil {
		// Make sure to return a non-zero exit code if ObtainSANCertificate returned at least one error.
		// Due to us not returning partial certificate we can just exit here instead of at the end.
		klog.Fatalf("Could not obtain certificates:\n\t%v", err)
	}

	certsStorage.SaveResource(cert)

	meta := map[string]string{
		renewEnvAccountEmail: account.Email,
		renewEnvCertDomain:   cert.Domain,
		renewEnvCertPath:     certsStorage.GetFileName(cert.Domain, ".crt"),
		renewEnvCertKeyPath:  certsStorage.GetFileName(cert.Domain, ".key"),
	}

	return legox.LaunchHook(ctx.Run_hook, meta)
}

func handleTOS(ctx *global.CertcliOptions, client *lego.Client) bool {
	// Check for a global accept override
	if ctx.Accept_tos {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	log.Printf("Please review the TOS at %s", client.GetToSURL())

	for {
		fmt.Println("Do you accept the TOS? Y/n")
		text, err := reader.ReadString('\n')
		if err != nil {
			klog.Fatalf("Could not read from console: %v", err)
		}

		text = strings.Trim(text, "\r\n")
		switch text {
		case "", "y", "Y":
			return true
		case "n", "N":
			return false
		default:
			fmt.Println("Your input was invalid. Please answer with one of Y/y, n/N or by pressing enter.")
		}
	}
}

func register(ctx *global.CertcliOptions, client *lego.Client) (*registration.Resource, error) {
	accepted := handleTOS(ctx, client)
	if !accepted {
		klog.Fatal("You did not accept the TOS. Unable to proceed.")
	}

	if ctx.Eab {
		kid := ctx.Kid
		hmacEncoded := ctx.Hmac

		if kid == "" || hmacEncoded == "" {
			klog.Fatalf("Requires arguments --kid and --hmac.")
		}

		return client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: accepted,
			Kid:                  kid,
			HmacEncoded:          hmacEncoded,
		})
	}

	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func obtainCertificate(ctx *global.CertcliOptions, client *lego.Client, domains []string) (*certificate.Resource, error) {
	bundle := !ctx.No_bundle

	//domains := ctx.Domains
	if len(domains) > 0 {
		// obtain a certificate, generating a new private key
		request := certificate.ObtainRequest{
			Domains:                        domains,
			Bundle:                         bundle,
			MustStaple:                     ctx.Must_staple,
			PreferredChain:                 ctx.Preferred_Chain,
			AlwaysDeactivateAuthorizations: ctx.Always_Deactivate_authorizations,
		}
		return client.Certificate.Obtain(request)
	}

	// read the CSR
	csr, err := legox.ReadCSRFile(ctx.Csr)
	if err != nil {
		return nil, err
	}

	// obtain a certificate for this CSR
	return client.Certificate.ObtainForCSR(certificate.ObtainForCSRRequest{
		CSR:                            csr,
		Bundle:                         bundle,
		PreferredChain:                 ctx.Preferred_Chain,
		AlwaysDeactivateAuthorizations: ctx.Always_Deactivate_authorizations,
	})
}
