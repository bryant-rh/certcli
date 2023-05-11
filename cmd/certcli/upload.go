package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bryant-rh/certcli/pkg/certxctx"
	"github.com/bryant-rh/certcli/pkg/cloud/txcloud"
	"github.com/bryant-rh/certcli/pkg/global"
	legox "github.com/bryant-rh/certcli/pkg/lego"
	"github.com/bryant-rh/certcli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"k8s.io/klog/v2"
)

func NewCmdUpload(o *global.CertcliOptions) *cobra.Command {

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "上传基于Let's Encrypt 颁发的ssl证书至云平台(Upload the ssl certificate issued based on Let's Encrypt to the cloud platform)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			//校验参数
			hasDomains := len(o.Domains) > 0
			if !hasDomains && o.FileName == "" {
				cmd.Help()
				klog.Fatal(util.RedColor("Please specify either --domains/-d or --filename/-f, but not both"))
			}

			if hasDomains && o.FileName != "" {
				klog.Fatal(util.RedColor("Please specify either --domains/-d or --filename/-f, but not both"))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			klog.V(4).Infoln(util.GreenColor("Run upload cert!"))
			err := upload(o)
			if err != nil {
				klog.Fatal(err)
				//return err
			}

			return nil

		},
	}

	uploadCmd.Flags().StringVarP(&o.Provider, "provider", "p", "txcloud", "指定云厂商,暂时只支持(txcloud)")
	uploadCmd.Flags().StringVarP(&o.FileName, "filename", "f", "", "通过文件指定域名,一行一个域名(Specify the domain name by file, one domain name per line)")

	return uploadCmd
}
func upload(ctx *global.CertcliOptions) error {
	//获取域名
	if ctx.FileName != "" {
		// 读取文件解析
		bytes, err := ioutil.ReadFile(ctx.FileName)
		if err != nil {
			klog.Fatalf("文件: %s 读取失败.", ctx.FileName)
		}
		str := string(bytes)
		str = strings.Replace(str, "\r", "", -1)
		str = strings.Trim(str, " ")
		domainsList = strings.Split(str, "\n")
		if len(domainsList) == 0 {
			klog.Fatal("文件内容为空")
		}
	}
	if len(ctx.Domains) > 0 {
		domainsList = ctx.Domains
	}

	if ctx.AllCert {
		fmt.Println(ctx.AllCert)

	} else {
		certx = certxctx.NewClient(ctx.Provider, global.Region)
		for _, d := range domainsToSlice(domainsList) {
			if len(d) > 0 {
				//根据域名查询证书id
				mes := fmt.Sprintf("开始upload操作, 域名: [%s]", d[0])
				klog.Infoln(util.GreenColor(mes))
				klog.Infoln(util.GreenColor("根据域名查询证书id"))
				id, err := certx.GetCertificateByDomain(&d[0])
				if err != nil {
					if strings.Contains(err.Error(), "no such Server Certificate domain") {
						if id == "" {
							//生成证书
							klog.Infoln(util.GreenColor("生成证书"))

							cert, err := legox.GenerateCertificate(ctx, d)
							if err != nil {
								klog.Errorln(util.RedColor("生成证书失败"))
								return err
							}

							certificate := txcloud.Certificate{
								Domains:     common.StringPtr(cert.Domain),
								PrivateKey:  common.StringPtr(string(cert.PrivateKey)),
								Certificate: common.StringPtr(string(cert.Certificate)),
							}

							//upload 证书
							klog.Infoln(util.GreenColor("upload 证书"))

							newId, err := certx.UploadServerCertificate(certificate)
							if err != nil {
								klog.Errorln(util.RedColor("upload 证书失败"))
								return err
							}
							//if !util.IsNil(newId) {
							if newId != "" {
								//根据证书id,查询证书详情
								klog.Infoln(util.GreenColor("根据证书id,查询证书详情"))
								cert, err := certx.GetCertificateDetailByID(&newId)
								if err != nil {
									klog.Errorln(util.RedColor("根据证书id,查询证书详情失败"))
									return err
								}
								message := fmt.Sprintf("域名: [%s] upload 成功\n证书ID: [%s]\n证书生效时间: [%s]\n证书失效时间: [%s]\n", d[0], newId, *cert.CertBeginTime, *cert.CertEndTime)
								fmt.Println(util.GreenColor(message))

							}

						}
					} else {
						klog.Errorln(util.RedColor("根据域名查询证书id 失败"))
						return err
					}
				}
				if id != "" {
					message := fmt.Sprintf("域名: [%s] 证书已存在!", d[0])
					fmt.Println(util.RedColor(message))
				}

			}

		}

	}

	return nil
}
