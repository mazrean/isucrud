package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mazrean/isucrud/dbdoc"
)

var (
	version  = "Unknown"
	revision = "Unknown"

	versionFlag                  bool
	dst                          string
	ignores                      sliceString
	ignorePrefixes               sliceString
	ignoreMain, ignoreInitialize bool
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "show version")

	flag.StringVar(&dst, "dst", "./dbdoc.md", "destination file")
	flag.Var(&ignores, "ignore", "ignore function")
	flag.Var(&ignorePrefixes, "ignorePrefix", "ignore function")
	flag.BoolVar(&ignoreMain, "ignoreMain", true, "ignore main function")
	flag.BoolVar(&ignoreInitialize, "ignoreInitialize", true, "ignore functions with 'initialize' in the name")
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

	err = dbdoc.Run(dbdoc.Config{
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
}
