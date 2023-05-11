package certxctx

import (
	"github.com/bryant-rh/certcli/pkg/cloud"
	"github.com/bryant-rh/certcli/pkg/cloud/txcloud"
	"k8s.io/klog/v2"
)

type CertClient interface {
	//GetCertificateByID 根据证书id 查询证书详情
	GetCertificateDetailByID(scID *string) (sc *cloud.CertificateDetail, err error)

	//GetAllCertificateIDs 查询所有证书的id ,20个id为一组
	GetAllCertificateIDs() (ids [][]*string, err error)

	//GetCertificateByDomain 根据域名查询证书ID
	GetCertificateByDomain(domain *string) (id string, err error)

	//GetDeployedResourcesByID 根据证书id 查询证书关联资源
	GetDeployedResourcesByID(scID [][]*string, resourceType string) (resp []*cloud.ResourcesItem, err error)

	//DeleteCertificateByID 根据证书id
	DeleteCertificateByID(scID string) (res bool, err error)

	//ReplaceCertForLoadBalancerByID 替换负载均衡实例关联的服务端证书
	ReplaceCertForLoadBalancerByID(oldID, newId string) (requestid string, err error)

	//UploadCertificate 上传证书
	UploadServerCertificate(cert txcloud.Certificate) (string, error)

	//UpdateDomainConfigByDomain 为指定域名的CDN更换SSL证书
	UpdateDomainConfigByDomain(domain, certID *string) (requestid string, err error)
}

// NewClient 根据 Provider 返回相应 DNS 客户端
func NewClient(provider string, region string) CertClient {

	switch provider {
	// case "aliyun":
	// 	return aliyun.NewClient()
	case "txcloud":
		txclient, err := txcloud.NewTxClient(region)
		if err != nil {
			klog.Fatalf("Could not create  NewTxClient:\n\t%v", err)
		}
		return txclient
	default:
		klog.Fatalf("Provider(%s) : 暂不支持不支持该供应商", provider)
	}

	return nil
}
