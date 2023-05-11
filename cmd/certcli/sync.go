package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bryant-rh/certcli/pkg/certxctx"
	"github.com/bryant-rh/certcli/pkg/cloud/txcloud"
	"github.com/bryant-rh/certcli/pkg/global"
	legox "github.com/bryant-rh/certcli/pkg/lego"
	"github.com/bryant-rh/certcli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"k8s.io/klog/v2"
)

func NewCmdSync(o *global.CertcliOptions) *cobra.Command {

	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "检查上传至云平台的ssl证书是否过期,并自动更新证书及关联资源",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			//校验参数
			hasDomains := len(o.Domains) > 0
			hasID := o.CertNewID != "" && o.CertOldID != ""

			if o.CertNewID != "" && o.CertOldID == "" {
				klog.Fatal(util.RedColor("Please specify both --new-id and --old-id,"))
			}
			if o.CertNewID == "" && o.CertOldID != "" {
				klog.Fatal(util.RedColor("Please specify both --new-id and --old-id,"))
			}

			if !hasDomains && o.FileName == "" && !o.AllCert && !hasID {
				cmd.Help()
				klog.Fatal(util.RedColor("Please specify either --domains/-d or --filename/-f  or --new-id and --old-id or --all-cert/-A, but not both"))
			}

			if hasDomains && o.AllCert {
				klog.Fatal(util.RedColor("Please specify either --domains/-d or --all-cert/-A, but not both"))
			}
			if o.FileName != "" && o.AllCert {
				klog.Fatal(util.RedColor("Please specify either --filename/-f or --all-cert/-A, but not both"))
			}
			if hasDomains && o.FileName != "" {
				klog.Fatal(util.RedColor("Please specify either --domains/-d or --filename/-f, but not both"))
			}
			if hasID && o.FileName != "" {
				klog.Fatal(util.RedColor("Please specify either --new-id and --old-id or --filename/-f, but not both"))
			}
			if hasID && hasDomains {
				klog.Fatal(util.RedColor("Please specify either --new-id and --old-id or --domains/-d, but not both"))
			}
			if hasID && o.AllCert {
				klog.Fatal(util.RedColor("Please specify either --new-id and --old-id or --all-cert/-A, but not both"))
			}
			if len(o.ResourceType) > 0 {
				for _, t := range o.ResourceType {
					if !util.MapKeyInIntSlice(resourceTypeSlice, t) {
						klog.Fatal(util.RedColor("--resource-type/-r accepts only clb,cdn"))
					}
				}

			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			klog.V(4).Infoln(util.GreenColor("Run Sync cert!"))
			err := sync(o)
			if err != nil {
				klog.Fatal(err)
				//return err
			}

			return nil

		},
	}

	syncCmd.Flags().StringVarP(&o.Provider, "provider", "p", "txcloud", "指定云厂商,暂时只支持(txcloud)")
	syncCmd.Flags().BoolVarP(&o.AllCert, "all-cert", "A", o.AllCert, "指定此选项,即会监听所有证书进行定时更新,谨慎使用！！！")
	syncCmd.Flags().StringVarP(&o.FileName, "filename", "f", "", "通过文件指定域名,一行一个域名(Specify the domain name by file, one domain name per line)")
	syncCmd.Flags().IntVar(&o.Days, "days", 14, "指定证书还剩下多少天可以更新(The number of days left on a certificate to renew it)")
	syncCmd.Flags().StringSliceVarP(&o.ResourceType, "resource-type", "r", o.ResourceType, "指定需要更新的资源, 可指定多个,以逗号分割(目前只支持资源类型:clb,cdn)")
	syncCmd.Flags().StringVar(&o.CertNewID, "new-id", "", "指定要更新的新证书ID")
	syncCmd.Flags().StringVar(&o.CertOldID, "old-id", "", "指定要更新的旧证书ID")
	syncCmd.MarkFlagRequired("resource-type")
	return syncCmd
}
func sync(ctx *global.CertcliOptions) error {
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

	certx = certxctx.NewClient(ctx.Provider, global.Region)

	//监听所有证书进行更新
	if ctx.AllCert {
		ids, err := certx.GetAllCertificateIDs()
		if err != nil {
			return err
		}
		for _, id := range ids {
			for _, i := range id {
				err = syncByCertID(certx, ctx, *i)
				if err != nil {
					return err
				}
			}

		}
		//根据指定域名进行更新证书
	} else if len(ctx.Domains) > 0 {
		domainsList = ctx.Domains

		err := syncByDomain(certx, ctx, domainsList)
		if err != nil {
			return err
		}
		//根据指定 新证书和旧证书ID 更新证书关联资源
	} else if ctx.CertNewID != "" && ctx.CertOldID != "" {
		err := updateByCertID(certx, ctx, ctx.CertNewID, ctx.CertOldID)
		if err != nil {
			return err
		}

	}

	return nil
}

//needRenewal 判断证书到期时间是否在days 内
func needRenewal(domain string, days int, endtime string) bool {
	now := time.Now()
	delta := time.Hour * 24 * time.Duration(days)
	certEndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endtime, time.Local)
	// 大于 days 天
	sub := certEndTime.Local().Sub(now.Local())

	if sub > delta {
		klog.Infof("域名:[%s], 证书有效期大于 %d 天, 无需更新, 剩余时间 %s", domain, days, sub.String())
		return false
	}
	return true
}

//isDeployResource 查询证书是否还存在关联资源
func isDeployResource(certx certxctx.CertClient, id string) (bool, error) {
	ResourceType := []string{"clb", "cdn", "live", "waf", "antiddos"}

	mes := fmt.Sprintf("根据证书id: [%s], 查询证书是否还存在关联资源", id)
	klog.Infoln(util.GreenColor(mes))

	for _, t := range ResourceType {
		IDSlice := [][]*string{}
		IDs := []*string{}
		IDs = append(IDs, &id)
		IDSlice = append(IDSlice, IDs)
		res, err := certx.GetDeployedResourcesByID(IDSlice, t)
		if err != nil {
			klog.Errorln(util.RedColor("根据证书id,查询关联资源 失败"))
			return false, err
		}
		// 判断证书是否有关联资源,有继续,没有则退出
		if len(res) != 0 {
			m := fmt.Sprintf("证书id: [%s], 还存在关联资源, 资源类型: [%s]", id, t)
			klog.Infoln(util.GreenColor(m))
			return true, nil

		}
	}
	return false, nil
}

//syncByDomain 根据域名更新云平台证书及关联资源
func syncByDomain(certx certxctx.CertClient, ctx *global.CertcliOptions, domains []string) error {
	for _, d := range domainsToSlice(domains) {
		if len(d) > 0 {
			mes := fmt.Sprintf("开始sync操作, 域名: [%s]", d[0])
			klog.Infoln(util.GreenColor(mes))
			klog.Infoln(util.GreenColor("根据域名查询证书id"))

			//根据域名查询证书id
			id, err := certx.GetCertificateByDomain(&d[0])
			if err != nil {
				klog.Errorln(util.RedColor("根据域名查询证书id 失败"))
				return err
			}
			err = syncByCertID(certx, ctx, id)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

//syncByCertID 根据证书id更新云平台证书及关联资源
func syncByCertID(certx certxctx.CertClient, ctx *global.CertcliOptions, id string) error {
	var newId string
	var updateflag bool
	for _, t := range ctx.ResourceType {
		//根据证书id，查询关联资源
		mes := fmt.Sprintf("根据证书id: [%s],查询关联资源类型: [%s]", id, t)
		klog.Infoln(util.GreenColor(mes))

		IDSlice := [][]*string{}
		IDs := []*string{}
		IDs = append(IDs, &id)
		IDSlice = append(IDSlice, IDs)
		res, err := certx.GetDeployedResourcesByID(IDSlice, t)
		if err != nil {
			klog.Errorln(util.RedColor("根据证书id,查询关联资源 失败"))
			return err
		}
		// 判断证书是否有关联资源,有继续,没有则退出
		if len(res) == 0 {
			mes := fmt.Sprintf("证书没有绑定类型为: [%s] 的资源", t)
			klog.Infof(util.RedColor(mes))

			continue
		}

		klog.Infoln(util.GreenColor("根据证书id,查询证书详情"))
		//根据证书id,查询证书详情
		certDetail, err := certx.GetCertificateDetailByID(&id)
		if err != nil {
			klog.Errorln(util.RedColor("根据证书id,查询证书详情 失败"))
			return err
		}
		detail := fmt.Sprintf("证书ID: [%s], 证书名称: [%s], 证书域名: [%s], 证书到期时间: [%s]", id, *certDetail.Alias, *certDetail.Domain, *certDetail.CertEndTime)
		klog.Infoln(util.GreenColor(detail))

		domain := []string{}
		domain = append(domain, *certDetail.Domain)

		if newId == "" {
			klog.Infoln(util.GreenColor("判断证书到期时间是否在days 内"))
			//判断证书到期时间是否在days 内
			if needRenewal(domain[0], ctx.Days, *certDetail.CertEndTime) {
				//生成证书
				klog.Infoln(util.GreenColor("证书即将到期, 开始生成证书"))
				cert, err := legox.GenerateCertificate(ctx, domain)
				if err != nil {
					return err
				}
				certificate := txcloud.Certificate{
					Domains:     common.StringPtr(cert.Domain),
					PrivateKey:  common.StringPtr(string(cert.PrivateKey)),
					Certificate: common.StringPtr(string(cert.Certificate)),
				}
				//upload 证书
				klog.Infoln(util.GreenColor("upload 证书"))
				newId, err = certx.UploadServerCertificate(certificate)
				if err != nil {
					klog.Errorln(util.RedColor("upload 证书失败"))
					return err
				}

				a := fmt.Sprintf("upload 证书成功, 新证书ID: [%s]", newId)
				klog.Infoln(util.GreenColor(a))

			} else {
				klog.Infoln(util.GreenColor("证书未到期, 无需更新!"))
				return nil
			}

		}
		//更新证书关联资源

		b := fmt.Sprintf("开始更新关联资源, 资源类型: [%s]", t)
		klog.Infoln(util.GreenColor(b))
		var reqID string
		switch t {
		case "clb":
			reqID, err = certx.ReplaceCertForLoadBalancerByID(id, newId)
			if err != nil {
				klog.Errorln(util.RedColor("更新关联资源: [clb] 失败"))
				return err
			}
		case "cdn":
			reqID, err = certx.UpdateDomainConfigByDomain(&domain[0], &newId)
			if err != nil {
				klog.Errorln(util.RedColor("更新关联资源: [cdn] 失败"))
				return err
			}
		}
		mes = fmt.Sprintf("关联资源[%s], 更新成功, requestID: [%s]", t, reqID)
		klog.Infoln(util.GreenColor(mes))
		updateflag = true

	}
	if updateflag {
		deploy, err := isDeployResource(certx, id)
		if err != nil {
			return err
		}

		if !deploy {
			mes := fmt.Sprintf("证书id: [%s], 不存在关联资源,准备删除该证书", id)
			klog.Infoln(util.GreenColor(mes))
			deleteResult, err := certx.DeleteCertificateByID(id)
			if err != nil {
				klog.Errorln(util.RedColor("删除证书失败!"))
				return err
			}
			if !deleteResult {
				klog.Fatal(util.RedColor("删除证书失败!"))
			} else {
				klog.Infoln(util.GreenColor("证书删除成功"))
			}

		}

	}
	return nil
}

//syncByCertID 根据手动上传的新证书id更新云平台关联资源
func updateByCertID(certx certxctx.CertClient, ctx *global.CertcliOptions, newID string, oldID string) error {
	var updateflag bool

	mes := fmt.Sprintf("开始sync操作, 新证书ID: [%s], 旧证书ID: [%s]", newID, oldID)
	klog.Infoln(util.GreenColor(mes))

	//根据证书id,查询证书详情
	klog.Infoln(util.GreenColor("根据证书id,查询证书详情"))
	certNewIdDetail, err := certx.GetCertificateDetailByID(&newID)
	if err != nil {
		klog.Errorln(util.RedColor("根据新证书id,查询证书详情 失败"))
		return err
	}
	newIDdetail := fmt.Sprintf("新证书ID: [%s], 新证书名称: [%s] , 新证书域名: [%s], 新证书到期时间: [%s]", newID, *certNewIdDetail.Alias, *certNewIdDetail.Domain, *certNewIdDetail.CertEndTime)
	klog.Infoln(util.GreenColor(newIDdetail))

	certOldIdDetail, err := certx.GetCertificateDetailByID(&oldID)
	if err != nil {
		klog.Errorln(util.RedColor("根据新证书id,查询证书详情 失败"))
		return err
	}
	oldIDdetail := fmt.Sprintf("旧证书ID: [%s], 旧证书名称: [%s] , 旧证书域名: [%s], 旧证书到期时间: [%s]", oldID, *certOldIdDetail.Alias, *certOldIdDetail.Domain, *certOldIdDetail.CertEndTime)
	klog.Infoln(util.GreenColor(oldIDdetail))

	if *certNewIdDetail.Domain != *certOldIdDetail.Domain {
		klog.Fatal(util.RedColor("新旧证书域名不一致,请确认证书ID是否正确"))
		return nil
	}

	for _, t := range ctx.ResourceType {
		//根据证书id，查询关联资源
		mes := fmt.Sprintf("根据旧证书id: [%s],查询关联资源类型: [%s]", oldID, t)
		klog.Infoln(util.GreenColor(mes))
		IDSlice := [][]*string{}
		IDs := []*string{}
		IDs = append(IDs, &oldID)
		IDSlice = append(IDSlice, IDs)
		res, err := certx.GetDeployedResourcesByID(IDSlice, t)
		if err != nil {
			klog.Errorln(util.RedColor("根据旧证书id,查询关联资源 失败"))
			return err
		}
		// 判断证书是否有关联资源,有继续,没有则退出
		if len(res) == 0 {
			mes := fmt.Sprintf("证书没有绑定类型为: [%s] 的资源", t)
			klog.Infof(util.RedColor(mes))

			continue
		}

		//更新证书关联资源

		b := fmt.Sprintf("开始更新关联资源, 资源类型: [%s]", t)
		klog.Infoln(util.GreenColor(b))
		var reqID string
		switch t {
		case "clb":
			reqID, err = certx.ReplaceCertForLoadBalancerByID(oldID, newID)
			if err != nil {
				klog.Errorln(util.RedColor("更新关联资源: [clb] 失败"))
				return err
			}
		case "cdn":
			reqID, err = certx.UpdateDomainConfigByDomain(certOldIdDetail.Domain, &newID)
			if err != nil {
				klog.Errorln(util.RedColor("更新关联资源: [cdn] 失败"))
				return err
			}
		}
		mes = fmt.Sprintf("关联资源[%s], 更新成功, requestID: [%s]", t, reqID)
		klog.Infoln(util.GreenColor(mes))
		updateflag = true

	}
	if updateflag {
		deploy, err := isDeployResource(certx, oldID)
		if err != nil {
			return err
		}

		if !deploy {
			mes := fmt.Sprintf("旧证书id: [%s], 不存在关联资源,准备删除该证书", oldID)
			klog.Infoln(util.GreenColor(mes))
			deleteResult, err := certx.DeleteCertificateByID(oldID)
			if err != nil {
				klog.Errorln(util.RedColor("删除旧证书失败!"))
				return err
			}
			if !deleteResult {
				klog.Fatal(util.RedColor("删除旧证书失败!"))
			} else {
				klog.Infoln(util.GreenColor("旧证书删除成功"))
			}

		}

	}

	return nil
}
