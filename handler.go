package isucontools

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func InitStaticFiles(callback func(urlpath string, handler http.Handler), prefix string) {
	wf := func(path string, info os.FileInfo, err error) error {
		log.Println(path, info, err)
		if path == prefix {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		urlpath := path[len(prefix):]
		if urlpath[0] != '/' {
			urlpath = "/" + urlpath
		}
		log.Println("Registering", urlpath, path)
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			return nil
		}
		content := make([]byte, info.Size())
		f.Read(content)
		f.Close()

		handler := func(w http.ResponseWriter, r *http.Request) {
			if path[len(path)-4:] == ".css" {
				w.Header().Set("Content-Type", "text/css")
			} else if path[len(path)-3:] == ".js" {
				w.Header().Set("Content-Type", "application/javascript")
			}
			w.Write(content)
		}
		callback(urlpath, http.HandlerFunc(handler))
		return nil
	}
	filepath.Walk(prefix, wf)
}
