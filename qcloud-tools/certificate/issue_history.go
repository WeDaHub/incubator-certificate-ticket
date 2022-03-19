package certificate

import (
	"fmt"
	"qcloud-tools/core/db"
	"strings"
	"time"
)

type IssueHistory struct {
	IssueDomain string
	PublicKey   string
	PrivateKey  string
	CreatedAt   uint
}

func (history IssueHistory) Add() {
	sql := "INSERT INTO issue_history (issue_domain,public_key,private_key,created_at) values (?, ?, ?, ?)"
	_, _ = db.QcloudToolDb.Insert(sql,
		history.IssueDomain,
		history.PublicKey,
		history.PrivateKey,
		history.CreatedAt)
}

func GetLatestValidRecord(domain string) (history IssueHistory) {

	sqlStr := "SELECT issue_domain,public_key,private_key,created_at FROM issue_history WHERE issue_domain in ('%s') AND created_at > %d ORDER BY id DESC LIMIT 1"
	now := time.Now().Unix()

	var arr []string
	arr = append(arr, domain)

	index := strings.Index(domain, ".")

	arr = append(arr, "*"+domain[index:])

	domain = strings.Join(arr, "','")

	sqlStr = fmt.Sprintf(sqlStr, domain, now-86400*31)

	rows, err := db.QcloudToolDb.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&history.IssueDomain,
			&history.PublicKey,
			&history.PrivateKey,
			&history.CreatedAt)

		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

func GetIssueInfoById(id uint64) (info IssueInfo) {
	sqlStr := "SELECT id,dns_api,app_id,app_id_value,app_key,app_key_value,main_domain,extra_domain FROM issue_info WHERE id = ? LIMIT 1"

	rows, err := db.QcloudToolDb.Query(sqlStr, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&info.Id,
			&info.DnsApi,
			&info.AppIdName,
			&info.AppIdValue,
			&info.AppKeyName,
			&info.AppKeyValue,
			&info.MainDomain,
			&info.ExtraDomain)

		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

func GetIssueInfoList() (issueList []IssueInfo) {

	sql := "SELECT id,main_domain,extra_domain FROM issue_info"
	rows, err := db.QcloudToolDb.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var info IssueInfo
		err = rows.Scan(
			&info.Id,
			&info.MainDomain,
			&info.ExtraDomain)

		if err != nil {
			fmt.Println(err)
		}

		issueList = append(issueList, info)
	}

	return
}
