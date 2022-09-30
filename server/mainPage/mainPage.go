package mainpage

import (
	filemanager "digitalDistribution/fileManager"
	"html/template"
	"net/http"
	"strings"
)

type TableRow struct {
	RevisionNumber string
	BuildDate      string
	DownloadFile   string
}

type Data struct {
	Items []TableRow
}

func getCurrentData() Data {
	dataRaw := filemanager.GetCurrentData()
	var items []TableRow
	for _, item := range dataRaw[:] {
		var tableRow TableRow
		tableRow.RevisionNumber = item[1]
		tableRow.BuildDate = item[2]
		fileStrings := strings.Split(item[3], "/")
		tableRow.DownloadFile = fileStrings[len(fileStrings)-1]
		items = append(items, tableRow)
	}
	return Data{Items: items}
}

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	tmpl, _ := template.ParseFiles("./server/views/mainPage.html")
	tmpl.Execute(res, getCurrentData())
}
