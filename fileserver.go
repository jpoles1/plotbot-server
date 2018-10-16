package main

import (
	"net/http"
	"path"
	"strings"
)

func fileServer(publicDir string, production bool) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !production {
			if origin := req.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}

		_path := req.URL.Path

		// static files
		if strings.Contains(_path, ".") || _path == "/" {
			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join(publicDir, "/index.html"))
	})
}
