package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type wasmServer struct {
	fs   http.Handler
	wasm string
}

func (s *wasmServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/", "/index.html":
		s.index(w, req)
	case "/wasm_exec.js":
		s.execJS(w, req)
	default:
		s.fs.ServeHTTP(w, req)
	}
}

func (s *wasmServer) index(w http.ResponseWriter, req *http.Request) {
	html := strings.Replace(indexHtml, "main.wasm", s.wasm, -1)
	io.WriteString(w, html)
}

func (s *wasmServer) execJS(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	w.Write(execJSBytes)
}

var port = flag.Int("p", 8080, "web server binding port")
var initFlag = flag.Bool("init", false, "create html and javascript files in current directory")
var fileServerFlag = flag.Bool("fs", false, "run as FileServer only")

func main() {
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	if *initFlag {
		initWebFiles()
		return
	}

	fileServer := http.FileServer(http.Dir("."))
	if *fileServerFlag {
		err := http.ListenAndServe(addr, fileServer)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if flag.NArg() == 0 {
		log.Fatal("missing wasm file")
	}

	wasm := flag.Args()[0]
	if _, err := os.Stat(wasm); os.IsNotExist(err) {
		log.Fatal("wasm file not exist: ", wasm)
	}

	server := &wasmServer{
		fs:   fileServer,
		wasm: wasm,
	}
	http.Handle("/", server)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initWebFiles() error {
	names := []string{"index.html", "wasm_exec.js"}
	data := []string{indexHtml, string(execJSBytes)}
	for i, name := range names {
		fmt.Println("generate:", name)
		err := func() error {
			var fp *os.File
			var err error
			if fp, err = os.Create(name); err != nil {
				return err
			}
			defer fp.Close()
			io.WriteString(fp, data[i])
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
