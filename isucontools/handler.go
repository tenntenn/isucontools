package isucontools

import (
	"filepath"
	"net/http"
	"os"
)

func InitStaticFiles(f func(urlpath string, handler http.Handler), prefix string) {
	wf := func(path string, info os.FileInfo, err error) error {
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
		f(urlpath, handler)
		return nil
	}
	filepath.Walk(prefix, wf)
}
