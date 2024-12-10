package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mazrean/isucrud/dbdoc"
	"github.com/mazrean/isucrud/internal/ui"
)

var (
	version  = "Unknown"
	revision = "Unknown"

	versionFlag                  bool
	dst                          string
	ignores                      sliceString
	ignorePrefixes               sliceString
	ignoreMain, ignoreInitialize bool
	web                          bool
	addr                         string
	base                         string
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "show version")

	flag.StringVar(&dst, "dst", "./dbdoc.md", "destination file")
	flag.Var(&ignores, "ignore", "ignore function")
	flag.Var(&ignorePrefixes, "ignorePrefix", "ignore function")
	flag.BoolVar(&ignoreMain, "ignoreMain", true, "ignore main function")
	flag.BoolVar(&ignoreInitialize, "ignoreInitialize", true, "ignore functions with 'initialize' in the name")
	flag.BoolVar(&web, "web", false, "run as web server")
	flag.StringVar(&addr, "addr", "localhost:7070", "address to listen on")
	flag.StringVar(&base, "base", "/", "base for serving the web server")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Printf("iwrapper %s (revision: %s)\n", version, revision)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}

	nodes, err := dbdoc.Run(dbdoc.Config{
		WorkDir:             wd,
		BuildArgs:           flag.Args(),
		IgnoreFuncs:         ignores,
		IgnoreFuncPrefixes:  ignorePrefixes,
		IgnoreMain:          ignoreMain,
		IgnoreInitialize:    ignoreInitialize,
		DestinationFilePath: dst,
	})
	if err != nil {
		panic(fmt.Errorf("failed to run dbdoc: %w", err))
	}

	if web {
		mux := http.NewServeMux()
		mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("Content-Type", "text/html")

			targetNodeID := r.URL.Query().Get("node")

			err := ui.RenderHTML(w, nodes, targetNodeID, base)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}))

		log.Printf("open http://%s/ in your browser\n", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			panic(fmt.Errorf("server exit: %w", err))
		}
	} else {
		err = ui.RenderMarkdown(dst, nodes)
		if err != nil {
			panic(fmt.Errorf("failed to render markdown: %w", err))
		}
	}
}
