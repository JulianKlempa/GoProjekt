package mainpage

import (
	filemanager "digitalDistribution/fileManager"
	"fmt"
	"html/template"
	"net/http"
)

type TableRow struct {
	RevisionNumber string
	BuildDate      string
}

type Data struct {
	Items []TableRow
}

func getCurrentData() Data {
	dataRaw := filemanager.GetStorage().Files
	var items []TableRow
	for _, item := range dataRaw[:] {
		var tableRow TableRow
		tableRow.RevisionNumber = fmt.Sprintf("%d.%d", item.Version.MajorVersion, item.Version.MinorVersion)
		tableRow.BuildDate = item.BuildDate
		items = append(items, tableRow)
	}
	return Data{Items: items}
}

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	tmpl, _ := template.ParseFiles("./server/views/mainPage.html")
	tmpl.Execute(res, getCurrentData())
}
