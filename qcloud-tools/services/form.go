package services

import (
	"fmt"
	"html/template"
	"net/http"
	"qcloud-tools/certificate"
	"strconv"
)

const (
	GET  = "GET"
	POST = "POST"
)

type DnsApi struct {
	Id         string
	Name       string
	AppIdName  string
	AppKeyName string
}

type CdnType struct {
	Id   string
	Name string
}

var DnsApiList = map[string]DnsApi{
	"dns_dp": {
		Id:         "dns_dp",
		Name:       "dnspod",
		AppIdName:  "DP_Id",
		AppKeyName: "DP_Key",
	},
	"dns_cf": {
		Id:         "dns_cf",
		Name:       "Cloudflare",
		AppIdName:  "CF_Token",
		AppKeyName: "CF_Account_ID",
	},
	"dns_gd": {
		Id:         "dns_cf",
		Name:       "GoDaddy.com",
		AppIdName:  "GD_Key",
		AppKeyName: "GD_Secret",
	},
	"dns_aws": {
		Id:         "dns_aws",
		Name:       "Amazon Route53",
		AppIdName:  "AWS_ACCESS_KEY_ID",
		AppKeyName: "AWS_SECRET_ACCESS_KEY",
	},
	"dns_ali": {
		Id:         "dns_ali",
		Name:       "Aliyun",
		AppIdName:  "Ali_Key",
		AppKeyName: "Ali_Secret",
	},
}

var CdnTypeList = []CdnType{
	{
		Id:   "cdn",
		Name: "内容分发网络",
	},
	{
		Id:   "ecdn",
		Name: "全站加速网络",
	},
}

func AddDomain(writer http.ResponseWriter, request *http.Request) {

	if POST == request.Method {
		_ = request.ParseForm()

		var info certificate.IssueInfo
		info.DnsApi = request.Form.Get("dns_api")
		info.AppIdValue = request.Form.Get("app_id_value")
		info.AppKeyValue = request.Form.Get("app_key_value")
		info.MainDomain = request.Form.Get("main_domain")
		info.ExtraDomain = request.Form.Get("extra_domain")

		dnsApi, ok := DnsApiList[info.DnsApi]
		if !ok {
			msg := fmt.Sprintf("%s 不存在", info.DnsApi)
			fmt.Fprintf(writer, fmt.Sprintf(`{"code":1,"msg":"%s"}`, msg))
			return
		}

		info.AppIdName = dnsApi.AppIdName
		info.AppKeyName = dnsApi.AppKeyName

		if err := info.Add(); err != nil {
			fmt.Fprintf(writer, fmt.Sprintf(`{"code":1,"msg":"%s"}`, err))
		} else {
			fmt.Fprintf(writer, `{"code":0}`)
		}
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	rootPath := tools.GetRootPath()
	templatePath := fmt.Sprintf("%s/web/add-domain.html", rootPath)
	tpl, _ := template.ParseFiles(templatePath)

	var form = struct {
		DnsApiList map[string]DnsApi
	}{
		DnsApiList,
	}

	_ = tpl.Execute(writer, form)
}

func AddSync(writer http.ResponseWriter, request *http.Request) {
	if POST == request.Method {
		_ = request.ParseForm()

		var sync certificate.IssueSync

		sync.SecretId = request.Form.Get("secret_id")
		sync.SecretKey = request.Form.Get("secret_key")
		sync.CdnType = request.Form.Get("type")
		sync.CdnDomain = request.Form.Get("cdn_domain")
		sync.IssueId, _ = strconv.ParseUint(request.Form.Get("issue_id"), 10, 64)

		if err := sync.Add(); err != nil {
			fmt.Fprintf(writer, fmt.Sprintf(`{"code":1,"msg":"%s"}`, err))
		} else {
			fmt.Fprintf(writer, `{"code":0}`)
		}

		return
	}

	writer.Header().Set("Content-Type", "text/html")
	rootPath := tools.GetRootPath()
	templatePath := fmt.Sprintf("%s/web/add.html", rootPath)
	tpl, _ := template.ParseFiles(templatePath)

	issueList := certificate.GetIssueInfoList()

	var form = struct {
		IssueInfoList []certificate.IssueInfo
		CdnTypeList   []CdnType
	}{
		issueList,
		CdnTypeList,
	}

	_ = tpl.Execute(writer, form)
}

func CheckLogin(writer http.ResponseWriter, request *http.Request) {

}
