package services

import (
	"fmt"
	"html/template"
	"net/http"
	"qcloud-tools/certificate"
	"qcloud-tools/core"
	"qcloud-tools/core/db"
)

type SyncList struct {
	Item []certificate.IssueSync
}

func GetList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	rootPath := core.GetRootPath()
	templatePath := fmt.Sprintf("%s/web/list.html", rootPath)
	tpl, err := template.ParseFiles(templatePath)
	if nil != err {
		fmt.Fprint(writer, "<div>Error ~~</div>")
		return
	}

	sqlStr := "SELECT id,type,cdn_domain,issue_id,last_issue_time,last_check_time FROM issue_sync order by id desc"
	rows, err := db.QcloudToolDb.Query(sqlStr)
	if err != nil {
		fmt.Fprint(writer, "<div>Error Query~~</div>")
		return
	}
	defer rows.Close()

	var list SyncList
	for rows.Next() {
		var issue certificate.IssueSync
		err = rows.Scan(
			&issue.Id,
			&issue.CdnType,
			&issue.CdnDomain,
			&issue.IssueId,
			&issue.LastIssueTime,
			&issue.LastCheckTime)
		if nil != err {
			fmt.Println("scan row error:", err)
			continue
		}

		list.Item = append(list.Item, issue)
	}

	_ = tpl.Execute(writer, list)
}
