package lego

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"

	"github.com/bryant-rh/certcli/pkg/global"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func GenerateCertificate(ctx *global.CertcliOptions, domain []string) (*certificate.Resource, error) {
	var email string
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.Email == "" {
		email = global.DefaultEmail
	} else {
		email = ctx.Email
	}
	myUser := MyUser{
		Email: email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	// c := &dnspod.Config{
	// 	TTL:                600,
	// 	PropagationTimeout: dns01.DefaultPropagationTimeout,
	// 	PollingInterval:    dns01.DefaultPollingInterval,
	// 	HTTPClient: &http.Client{
	// 		Timeout: time.Second * 300,
	// 	},
	// 	//LoginToken: *global.DnsPodToken,
	// }
	// provide, err := dnspod.NewDNSProviderConfig(c)
	// if err != nil {
	// 	return nil, err
	// }
	// err = client.Challenge.SetDNS01Provider(provide, dns01.AddRecursiveNameservers([]string{"f1g1ns1.dnspod.net", "f1g1ns2.dnspod.net", "8.8.8.8"}))
	// if err != nil {
	// 	return nil, err
	// }
	SetupChallenges(ctx, client)

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: domain,
		Bundle:  true,
	}
	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
