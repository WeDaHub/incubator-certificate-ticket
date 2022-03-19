package certificate

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"qcloud-tools/core"
	"qcloud-tools/core/db"
	"regexp"
	"strings"
	"time"
)

type IssueSync struct {
	Id                           uint64
	SecretId, SecretKey          string
	CdnType                      string `default:"cdn"`
	CdnDomain                    string
	LoadBalancerId, CertName, Region     string
	IssueId                      uint64 `default:"0"`
	LastIssueTime, LastCheckTime uint
}

type IssueInfo struct {
	Id                      uint64
	DnsApi                  string
	AppIdName, AppIdValue   string
	AppKeyName, AppKeyValue string
	MainDomain, ExtraDomain string
}

func (info *IssueInfo) Add() (err error) {

	sql := "INSERT INTO issue_info (dns_api,app_id,app_id_value,app_key,app_key_value,main_domain,extra_domain) values (?, ?, ?, ?, ?, ?, ?)"
	_, err = db.QcloudToolDb.Insert(sql,
		info.DnsApi,
		info.AppIdName,
		info.AppIdValue,
		info.AppKeyName,
		info.AppKeyValue,
		info.MainDomain,
		info.ExtraDomain)

	return err
}

func (info *IssueInfo) GenerateScript() (string, error) {

	rootPath := core.GetRootPath()

	fileName := fmt.Sprintf("%s/shell/%s.sh", rootPath, info.MainDomain)
	f, err := os.Create(fileName)
	defer f.Close()

	if nil != err {
		fmt.Printf("创建文件失败: %s \n", fileName)
		return "", err
	}

	templatePath := fmt.Sprintf("%s/config/issue-template.tpl", rootPath)
	tpl, err := template.ParseFiles(templatePath)
	if nil != err {
		fmt.Printf("模板文件不存在: %s \n", templatePath)
		return "", err
	}

	if err := tpl.Execute(f, info); err != nil {
		fmt.Printf("生成脚本失败：%s \n", err)
		return "", err
	}

	if err = f.Chmod(0777); err != nil {
		fmt.Printf("更改文件权限失败：%s \n", err)
		return fileName, err
	}

	return fileName, nil
}

func (issue *IssueSync) Add() (err error) {
	sql := "INSERT INTO issue_sync (secret_id,secret_key,type,cdn_domain,issue_id) values (?, ?, ?, ?, ?)"
	_, err = db.QcloudToolDb.Insert(sql,
		issue.SecretId,
		issue.SecretKey,
		issue.CdnType,
		issue.CdnDomain,
		issue.IssueId)

	return err
}

func (issue *IssueSync) IssueCertByScript() bool {

	info := GetIssueInfoById(issue.IssueId)

	fileName, err := info.GenerateScript()
	if err != nil {
		return false
	}

	command := exec.Command(fileName)
	stdout, _ := command.Output()

	var privateKeyPath, publicKeyPath string

	content := string(stdout)

	fmt.Println("issue result", content)

	privateKeyRegexp, _ := regexp.Compile(`Your cert key is in: (.*\/\.acme\.sh\/.*[\S])`)
	publicKeyRegexp, _ := regexp.Compile(`And the full chain certs is there: (.*\/\.acme\.sh\/.*[\S])`)

	var regexpResult []string
	regexpResult = privateKeyRegexp.FindStringSubmatch(content)
	if nil != regexpResult {
		privateKeyPath = regexpResult[1]
	}
	regexpResult = publicKeyRegexp.FindStringSubmatch(content)
	if nil != regexpResult {
		publicKeyPath = regexpResult[1]
	}

	if "" == privateKeyPath || "" == publicKeyPath {
		fmt.Printf("update certificate failed,private %s, public %s \n", privateKeyPath, publicKeyPath)
		return false
	}

	publicData, _ := ioutil.ReadFile(publicKeyPath)
	privateData, _ := ioutil.ReadFile(privateKeyPath)

	publicKeyData := strings.TrimSpace(string(publicData))
	publicKeyData = strings.ReplaceAll(publicKeyData, "\n", "\\n")

	privateKeyData := strings.TrimSpace(string(privateData))
	privateKeyData = strings.ReplaceAll(privateKeyData, "\n", "\\n")

	now := uint(time.Now().Unix())

	history := IssueHistory{
		IssueDomain: info.MainDomain,
		PublicKey:   publicKeyData,
		PrivateKey:  privateKeyData,
		CreatedAt:   now,
	}
	history.Add()

	if "" != info.ExtraDomain {
		extraDomain := strings.Split(info.ExtraDomain, "-d ")
		for _, value := range extraDomain {
			value = strings.TrimSpace(value)
			if value != "" {
				history.IssueDomain = value
				history.Add()
			}
		}
	}

	// 更新证书到 cdn 或者 ecdn
	var syncInstance ISync
	sync := Sync{
		SecretId:       issue.SecretId,
		SecretKey:      issue.SecretKey,
		Domain:         issue.CdnDomain,
		PrivateKeyData: privateKeyData,
		PublicKeyData:  publicKeyData,
		Region:         issue.Region,
	}

	switch issue.CdnType {
	case "lb":
		syncInstance = LBSync{
			Sync:   sync,
			LoadBalancerId: issue.LoadBalancerId,
			CertName: issue.CertName,
		}
	default:
		syncInstance = CdnSync{sync}
	}

	return syncInstance.UpdateCredential()
}

func (issue *IssueSync) IssueCertByHistory() (bool, uint) {

	history := GetLatestValidRecord(issue.CdnDomain)
	if "" == history.PublicKey {
		return false, 0
	}

	// 更新证书到 cdn 或者 ecdn
	var syncInstance ISync
	sync := Sync{
		SecretId:       issue.SecretId,
		SecretKey:      issue.SecretKey,
		Domain:         issue.CdnDomain,
		PrivateKeyData: history.PrivateKey,
		PublicKeyData:  history.PublicKey,
		Region:         issue.Region,
	}

	switch issue.CdnType {
	case "lb":
		syncInstance = LBSync{
			Sync:   sync,
			LoadBalancerId: issue.LoadBalancerId,
			CertName: issue.CertName,
		}
	default:
		syncInstance = CdnSync{sync}
	}

	return syncInstance.UpdateCredential(), history.CreatedAt
}

func (issue *IssueSync) IssueCert() {

	result, now := issue.IssueCertByHistory()

	if !result {
		result = issue.IssueCertByScript()
		now = uint(time.Now().Unix())
	}

	// 更新数据库信息
	if result && issue.Id > 0 {
		sqlStr := "UPDATE issue_sync SET last_issue_time = ? WHERE id = ?"
		_, _ = db.QcloudToolDb.Update(sqlStr, now, issue.Id)
	}
}
