package lego

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bryant-rh/certcli/pkg/global"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"k8s.io/klog/v2"
)

const filePerm os.FileMode = 0o600

func Setup(ctx *global.CertcliOptions, accountsStorage *AccountsStorage) (*Account, *lego.Client) {
	keyType := GetKeyType(ctx)
	privateKey := accountsStorage.GetPrivateKey(keyType)

	var account *Account
	if accountsStorage.ExistsAccountFilePath() {
		account = accountsStorage.LoadAccount(privateKey)
	} else {
		account = &Account{Email: accountsStorage.GetUserID(), key: privateKey}
	}

	client := newClient(ctx, account, keyType)

	return account, client
}

func newClient(ctx *global.CertcliOptions, acc registration.User, keyType certcrypto.KeyType) *lego.Client {
	config := lego.NewConfig(acc)
	config.CADirURL = ctx.Acme_server

	config.Certificate = lego.CertificateConfig{
		KeyType: keyType,
		Timeout: time.Duration(ctx.Cert_timeout) * time.Second,
	}
	config.UserAgent = getUserAgent(ctx)

	if ctx.Http_timeout != 0 {
		config.HTTPClient.Timeout = time.Duration(ctx.Http_timeout) * time.Second
	}

	client, err := lego.NewClient(config)
	if err != nil {
		klog.Fatalf("Could not create client: %v", err)
	}

	if client.GetExternalAccountRequired() && !ctx.Eab {
		klog.Fatal("Server requires External Account Binding. Use --eab with --kid and --hmac.")
	}

	return client
}

// getKeyType the type from which private keys should be generated.
func GetKeyType(ctx *global.CertcliOptions) certcrypto.KeyType {
	//keyType :=
	switch strings.ToUpper(ctx.Key_type) {
	case "RSA2048":
		return certcrypto.RSA2048
	case "RSA4096":
		return certcrypto.RSA4096
	case "RSA8192":
		return certcrypto.RSA8192
	case "EC256":
		return certcrypto.EC256
	case "EC384":
		return certcrypto.EC384
	}

	klog.Fatalf("Unsupported KeyType: %s", ctx.Key_type)
	return ""
}

func getEmail(ctx *global.CertcliOptions) string {
	email := ctx.Email
	if email == "" {
		klog.Fatal("You have to pass an account (email address) to the program using --email or -m")
	}
	return email
}

func getUserAgent(ctx *global.CertcliOptions) string {
	return strings.TrimSpace(fmt.Sprintf("%s cert-cli", ctx.User_agent))
}

func CreateNonExistingFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o700)
	} else if err != nil {
		return err
	}
	return nil
}

func ReadCSRFile(filename string) (*x509.CertificateRequest, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	raw := bytes

	// see if we can find a PEM-encoded CSR
	var p *pem.Block
	rest := bytes
	for {
		// decode a PEM block
		p, rest = pem.Decode(rest)

		// did we fail?
		if p == nil {
			break
		}

		// did we get a CSR?
		if p.Type == "CERTIFICATE REQUEST" || p.Type == "NEW CERTIFICATE REQUEST" {
			raw = p.Bytes
		}
	}

	// no PEM-encoded CSR
	// assume we were given a DER-encoded ASN.1 CSR
	// (if this assumption is wrong, parsing these bytes will fail)
	return x509.ParseCertificateRequest(raw)
}
