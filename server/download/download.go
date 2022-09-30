package download

import (
	"fmt"
	"net/http"
	"strings"
)

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	queryStrings := strings.Split(req.URL.Query().Get("file"), "/")
	fmt.Println(queryStrings[len(queryStrings)-1])
	http.Redirect(res, req, "/", http.StatusPermanentRedirect)
}
