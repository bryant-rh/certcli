package txcloud

import (
	"fmt"
	"time"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

//UpdateDomainConfigByDomain 为指定域名的CDN更换SSL证书
func (c *TxClient) UpdateDomainConfigByDomain(domain, certID *string) (requestid string, err error) {
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = domain
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	mes := fmt.Sprintf("Auto update by api at %s", timeStr)
	request.Https = &cdn.Https{
		Switch: common.StringPtr("on"),
		CertInfo: &cdn.ServerCert{
			CertId:  certID,
			Message: common.StringPtr(mes),
		},
	}

	resp, err := c.cdnClient.UpdateDomainConfig(request)
	if err != nil {
		return requestid, err
	}
	return *resp.Response.RequestId, nil
}
