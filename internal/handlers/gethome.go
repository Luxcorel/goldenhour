package handlers

import (
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=60")
	http.ServeFile(w, r, "internal/templates/index.html")
}
