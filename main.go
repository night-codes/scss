package main

import (
	"bytes"
	"fmt"
	"github.com/alexflint/go-arg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	sass "github.com/wellington/go-libsass"
)

var (
	args struct {
		Input  string `arg:"positional" help:"Input SCSS file (\"scss/style.scss\")."`
		Output string `arg:"positional" help:"Output CSS file (\"style.css\")."`
		Watch  string `arg:"-w" help:"Watch directory for changes."`
		Import string `arg:"-i" help:"Specify a Scss import path."`
	}

	filelist = make(map[string]time.Time)
)

func main() {
	arg.MustParse(&args)
	sass.OutputStyle(sass.COMPRESSED_STYLE)

	if args.Output == "" {
		args.Output = "style.css"
	}
	if args.Input == "" {
		args.Input = "scss/style.scss"
	}
	if args.Watch == "" {
		makecss()
		os.Exit(0)
	}

	fmt.Printf("SCSS started to watch \"%s\"\n", args.Watch)

	for {
		changes := false
		filepath.Walk(args.Watch, func(path string, info os.FileInfo, err error) error {
			if err == nil && info != nil && !info.IsDir() {
				if t, ok := filelist[path]; !ok || t != info.ModTime() {
					filelist[path] = info.ModTime()
					if !ok {
						fmt.Println("Add \"" + path + "\"")
					} else {
						fmt.Println("Changed \"" + path + "\"")
					}
					changes = true
				}
			}
			return nil
		})
		if changes {
			makecss()
		}
		time.Sleep(time.Second / 10)
	}
}

func makecss() {
	cssbufer := bytes.Buffer{}
	scssfile, err := os.Open(args.Input)
	if err != nil {
		log.Println("Error: ", err)
	}
	comp, err := sass.New(&cssbufer, scssfile)

	if err != nil {
		log.Println("Error: ", err)
	}
	// configure @import paths
	if args.Import != "" {
		if err := comp.Option(sass.IncludePaths([]string{args.Import})); err != nil {
			log.Println("Error: ", err)
		}
	}
	if err := comp.Run(); err != nil {
		log.Println("Error: ", err)
	}
	if err := ioutil.WriteFile(args.Output, cssbufer.Bytes(), 0644); err != nil {
		log.Println("Error: ", err)
	}
}
