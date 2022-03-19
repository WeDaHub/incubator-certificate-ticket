package certificate

import (
	"fmt"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ecdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ecdn/v20191012"
	"time"
)

type ISync interface {
	UpdateCredential() bool
	GetCredential() (*common.Credential, *profile.ClientProfile)
	GetCertRequestParam() (params string)
}

type Sync struct {
	SecretId       string
	SecretKey      string
	Domain         string
	PrivateKeyData string
	PublicKeyData  string
	Region         string
}

func (sync Sync) GetCredential() (*common.Credential, *profile.ClientProfile) {
	credential := common.NewCredential(
		sync.SecretId,
		sync.SecretKey,
	)

	cpf := profile.NewClientProfile()

	return credential, cpf
}

func (sync Sync) GetCertRequestParam() (params string) {

	params = "{\"Domain\":\"%s\",\"Https\":{\"Switch\":\"on\",\"Http2\":\"on\",\"CertInfo\":{\"Certificate\":\"%s\",\"PrivateKey\":\"%s\",\"Message\":\"%s\"}}}"
	params = fmt.Sprintf(params, sync.Domain, sync.PublicKeyData, sync.PrivateKeyData, time.Now().Format("2006-01-02"))

	return
}

type CdnSync struct {
	Sync
}

func (sync CdnSync) UpdateCredential() bool {
	credential, cpf := sync.GetCredential()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"

	client, _ := cdn.NewClient(credential, sync.Region, cpf)
	request := cdn.NewUpdateDomainConfigRequest()

	params := sync.GetCertRequestParam()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	request.ForceRedirect = &cdn.ForceRedirect{
		Switch:             common.StringPtr("on"),
		RedirectType:       common.StringPtr("https"),
		RedirectStatusCode: common.Int64Ptr(301),
	}

	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	}
	if err != nil {
		fmt.Printf("UpdateDomainConfig returned: %s", err)
		return false
	}
	fmt.Printf("%s \n", response.ToJsonString())

	return true
}

type EcdnSync struct {
	Sync
}

func (sync EcdnSync) UpdateCredential() bool {
	credential, cpf := sync.GetCredential()
	cpf.HttpProfile.Endpoint = "ecdn.tencentcloudapi.com"

	client, _ := ecdn.NewClient(credential, sync.Region, cpf)
	request := ecdn.NewUpdateDomainConfigRequest()

	params := sync.GetCertRequestParam()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	request.ForceRedirect = &ecdn.ForceRedirect{
		Switch:             common.StringPtr("on"),
		RedirectType:       common.StringPtr("https"),
		RedirectStatusCode: common.Int64Ptr(301),
	}

	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Println("An API error has returned: ", err)
	}
	if err != nil {
		fmt.Println("An API error has returned: ", err)
		return false
	}
	fmt.Println(response.ToJsonString())

	return true
}

type LBSync struct {
	Sync
	LoadBalancerId string
	CertName       string
}

func (sync LBSync) UpdateCredential() bool {
	credential, cpf := sync.GetCredential()
	cpf.HttpProfile.Endpoint = "clb.tencentcloudapi.com"

	client, _ := clb.NewClient(credential, sync.Region, cpf)

	// CertId 用接口去查询，因为每次变更后就变了
	var certId string

	query := clb.NewDescribeListenersRequest()
	query.LoadBalancerId = common.StringPtr(sync.LoadBalancerId)
	query.Protocol = common.StringPtr("HTTPS")

	queryRsp, err := client.DescribeListeners(query)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return false
	}

	if err != nil {
		panic(err)
	}

	if *queryRsp.Response.TotalCount <= 0 {
		fmt.Printf("未查询到证书ID")
		return false
	}

	for _, item := range queryRsp.Response.Listeners {
		if nil != item.Certificate && "" != *item.Certificate.CertId {
			certId = *item.Certificate.CertId
		}
	}

	if certId == "" {
		return false
	}

	request := clb.NewReplaceCertForLoadBalancersRequest()
	params := "{\"OldCertificateId\":\"%s\",\"Certificate\":{\"CertName\":\"%s\",\"CertContent\":\"%s\",\"CertKey\":\"%s\"}}"
	params = fmt.Sprintf(params, certId, sync.CertName+time.Now().Format("2006-01-02"), sync.PublicKeyData, sync.PrivateKeyData)
	err = request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	response, err := client.ReplaceCertForLoadBalancers(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Println("An API error has returned: ", err)
	}
	if err != nil {
		fmt.Println("An API error has returned: ", err)
		return false
	}
	fmt.Println(response.ToJsonString())

	return true
}
