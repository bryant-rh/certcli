package global

import "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"

//定义所有命令行参数
type CertcliOptions struct {
	//botKey        = os.Getenv("BOT_KEY")
	EnableDebug                      bool
	Eab                              bool
	Tls                              bool
	Http                             bool
	Dns_Disablecp                    bool
	Pem                              bool
	Pfx                              bool
	Accept_tos                       bool
	No_bundle                        bool
	Must_staple                      bool
	Always_Deactivate_authorizations bool
	AllCert                          bool

	FileName         string
	Pfx_Pass         string
	Http_port        string
	Tls_port         string
	Http_ProxyHeader string
	Http_Webroot     string
	Dns              string
	Key_type         string
	Email            string
	Path             string
	Acme_server      string
	Csr              string
	Kid              string
	Hmac             string
	User_agent       string
	Run_hook         string
	Preferred_Chain  string
	Provider         string
	Code             string
	CertNewID        string
	CertOldID        string

	ResourceType  []string
	Dns_resolvers []string
	Domains       []string

	Cert_timeout int
	Dns_timeout  int
	Http_timeout int
	Days         int
}

var (
	Region       = regions.Guangzhou
	DefaultEmail = "op@sensorsdata.cn"
)

// var (
// 	AccessKey string
// 	SecretKey string
// 	Zone      string
// 	Txclient  common.Client
// )
