package txcloud

import (
	"fmt"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type TxClient struct {
	sslClient *ssl.Client
	clbClient *clb.Client
	cdnClient *cdn.Client
}

func NewTxClient(region string) (*TxClient, error) {
	credential, cpf, err := getCredential()
	if err != nil {
		return nil, err
	}

	sslclient, err := ssl.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("error creating txcloud ssl client: '%v'", err)
	}

	clbclient, err := clb.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("error creating txcloud clb client: '%v'", err)
	}
	cdnclient, err := cdn.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("error creating txcloud clb client: '%v'", err)
	}

	return &TxClient{
		sslClient: sslclient,
		clbClient: clbclient,
		cdnClient: cdnclient,
	}, nil

}

//getCredential/
func getCredential() (common.CredentialIface, *profile.ClientProfile, error) {

	provider := common.DefaultEnvProvider()
	credential, err := provider.GetCredential()
	if err != nil {
		return nil, nil, err

	}
	cpf := profile.NewClientProfile()

	return credential, cpf, nil
}
