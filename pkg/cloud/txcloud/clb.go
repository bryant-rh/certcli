package txcloud

import (
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type CertificateBinding struct {
	LoadbalancerListenerID string
	LoadBalancerID         string
}

// GetServerCertificateBindingByID 使用 证书ID 查询所绑定的 LoadBalance Listener 和 LB
func (c *TxClient) GetServerCertificateBindingByID(scID []*string) (res *clb.DescribeLoadBalancerListByCertIdResponse, err error) {
	request := clb.NewDescribeLoadBalancerListByCertIdRequest()
	request.CertIds = scID
	resp, err := c.clbClient.DescribeLoadBalancerListByCertId(request)
	if err != nil {
		return nil, err
	}

	// for _, lb := range resp.Response.CertSet {
	// 	bindings = append(bindings, CertificateBinding{lb.LoadBalancers., lbl.LoadbalancerID})
	// }

	return resp, nil
}

//ReplaceCertForLoadBalancerByID 替换负载均衡实例关联的服务端证书
func (c *TxClient) ReplaceCertForLoadBalancerByID(oldID, newId string) (requestid string, err error) {
	request := clb.NewReplaceCertForLoadBalancersRequest()
	request.OldCertificateId = common.StringPtr(oldID)
	request.Certificate = &clb.CertificateInput{
		CertId: common.StringPtr(newId),
	}
	//request.Certificate.CertId = common.StringPtr(newId)
	resp, err := c.clbClient.ReplaceCertForLoadBalancers(request)

	if err != nil {
		return requestid, err
	}
	return *resp.Response.RequestId, nil

}
