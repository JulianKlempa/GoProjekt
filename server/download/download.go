package download

import (
	filemanager "digitalDistribution/fileManager"
	"io"
	"net/http"
	"strings"
)

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	queryStrings := strings.Split(req.URL.Query().Get("file"), "/")
	res.Header().Set("Content-Disposition", "attachment; filename=Digital.zip")
	res.Header().Set("Content-Type", "application/zip")
	io.Copy(res, filemanager.GetFile(queryStrings[len(queryStrings)-1]))
	http.Redirect(res, req, "/", http.StatusPermanentRedirect)
}
