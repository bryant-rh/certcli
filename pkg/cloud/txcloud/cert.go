package txcloud

import (
	"fmt"
	"time"

	"github.com/bryant-rh/certcli/pkg/cloud"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Certificate struct {
	Domains     *string
	PrivateKey  *string
	Certificate *string
}

//GetAllCertificateIDs 查询所有证书的id ,20个id为一组
func (c *TxClient) GetAllCertificateIDs() (ids [][]*string, err error) {

	idList := []*string{}
	request := ssl.NewDescribeCertificatesRequest()
	request.Limit = common.Uint64Ptr(100)

	resp, err := c.sslClient.DescribeCertificates(request)
	if err != nil {
		return nil, err
	}

	for _, sc := range resp.Response.Certificates {
		idList = append(idList, sc.CertificateId)
	}

	for i := range idList {
		if i%20 != 0 {
			continue
		}
		if i+20 >= len(idList) {
			ids = append(ids, idList[i:])
		} else {
			ids = append(ids, idList[i:i+20])
		}
	}

	return

}

//GetCertificateByDomain 根据域名查询证书ID
func (c *TxClient) GetCertificateByDomain(domain *string) (id string, err error) {
	request := ssl.NewDescribeCertificatesRequest()
	request.SearchKey = domain

	resp, err := c.sslClient.DescribeCertificates(request)
	if err != nil {
		return id, err
	}

	if *resp.Response.TotalCount != 1 {
		return id, fmt.Errorf("no such Server Certificate domain : %s", *domain)
	}

	return *resp.Response.Certificates[0].CertificateId, nil
}

//GetCertificateByID 根据证书id 查询证书详情
func (c *TxClient) GetCertificateDetailByID(scID *string) (sc *cloud.CertificateDetail, err error) {
	request := ssl.NewDescribeCertificateDetailRequest()
	request.CertificateId = scID

	resp, err := c.sslClient.DescribeCertificateDetail(request)
	if err != nil {
		return nil, err
	}

	// if *resp.Response.TotalCount != 1 {
	// 	return sc, fmt.Errorf("no such Server Certificate id : %s", scID)
	// }

	sc = &cloud.CertificateDetail{
		From:                  resp.Response.From,
		CertificateType:       resp.Response.CertificateType,
		Domain:                resp.Response.Domain,
		Alias:                 resp.Response.Alias,
		Status:                resp.Response.Status,
		CertBeginTime:         resp.Response.CertBeginTime,
		CertEndTime:           resp.Response.CertEndTime,
		ValidityPeriod:        resp.Response.ValidityPeriod,
		CertificatePrivateKey: resp.Response.CertificatePrivateKey,
		CertificatePublicKey:  resp.Response.CertificatePublicKey,
		CertificateId:         resp.Response.CertificateId,
		IsWildcard:            resp.Response.IsWildcard,
		RequestId:             resp.Response.RequestId,
	}
	return
}

//GetDeployedResourcesByID 根据证书id 查询证书关联资源
func (c *TxClient) GetDeployedResourcesByID(scID [][]*string, resourceType string) (resp []*cloud.ResourcesItem, err error) {
	for _, id := range scID {
		request := ssl.NewDescribeDeployedResourcesRequest()
		request.CertificateIds = id
		request.ResourceType = common.StringPtr(resourceType)

		response, err := c.sslClient.DescribeDeployedResources(request)
		if err != nil {
			return resp, err
		}

		// if len(response.Response.DeployedResources) == 0 {
		// 	return nil, fmt.Errorf("no such resource type : %s", resourceType)
		// }
		for _, i := range response.Response.DeployedResources {
			if *i.Count != 0 {
				resp = append(resp, &cloud.ResourcesItem{
					CertificateId: i.CertificateId,
					Count:         i.Count,
					Type:          i.Type,
					Resources:     i.Resources,
				})
			}
		}
	}

	return resp, nil
}

//UploadCertificate 上传证书
func (c *TxClient) UploadServerCertificate(cert Certificate) (string, error) {
	prefix := time.Now().Format("2006-01-02")
	alias := fmt.Sprintf("%s-%s", prefix, *cert.Domains)

	request := ssl.NewUploadCertificateRequest()
	request.Alias = common.StringPtr(alias)
	request.CertificatePrivateKey = cert.PrivateKey
	request.CertificatePublicKey = cert.Certificate

	resp, err := c.sslClient.UploadCertificate(request)
	if err != nil {
		return "", err
	}

	return *resp.Response.CertificateId, nil
}

//DeleteCertificateByID 根据证书id删除证书
func (c *TxClient) DeleteCertificateByID(scID string) (res bool, err error) {
	request := ssl.NewDeleteCertificateRequest()
	request.CertificateId = common.StringPtr(scID)

	resp, err := c.sslClient.DeleteCertificate(request)

	if err != nil {
		return res, err
	}

	res = *resp.Response.DeleteResult
	return
}
