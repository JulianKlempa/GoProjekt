package download

import (
	filemanager "digitalDistribution/fileManager"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	queryStrings := strings.Split(req.URL.Query().Get("version"), "/")
	res.Header().Set("Content-Disposition", "attachment; filename=Digital.zip")
	res.Header().Set("Content-Type", "application/zip")
	versionStrings := strings.Split(queryStrings[len(queryStrings)-1], ".")
	version := filemanager.Version{}
	majorVersion, err := strconv.ParseInt(strings.TrimLeft(versionStrings[0], "v"), 0, 64)
	if err != nil {
		panic(err)
	}
	minorVersion, err := strconv.ParseInt(versionStrings[1], 0, 64)
	if err != nil {
		panic(err)
	}
	version.MajorVersion = int(majorVersion)
	version.MinorVersion = int(minorVersion)
	io.Copy(res, filemanager.GetFile(version))
	filemanager.IncreaseDownloadCounter(version)
	http.Redirect(res, req, "/", http.StatusPermanentRedirect)
}
