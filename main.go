package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	sass "github.com/wellington/go-libsass"
)

var (
	inputFile  = flag.String("i", "scss/style.scss", "Input SCSS file.")
	outputFile = flag.String("o", "css/style.css", "Output CSS file.")
	watchDir   = flag.String("watch", "", "Watch directory for changes.")
	filelist   = make(map[string]time.Time)
)

func main() {
	flag.Parse()
	sass.OutputStyle(sass.COMPRESSED_STYLE)
	
	if *watchDir == "" {
		makecss()
		os.Exit(0)
	}


	fmt.Printf("SCSS started to track \"%s\"\n", *watchDir)
	
	for {
		changes := false
		filepath.Walk(*watchDir, func(path string, info os.FileInfo, err error) error {
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
	scssfile, err := os.Open(*inputFile)
	if err != nil {
		log.Println("Error: ", err)
	}
	comp, err := sass.New(&cssbufer, scssfile)

	if err != nil {
		log.Println("Error: ", err)
	}
	// configure @import paths
	if err := comp.Option(sass.IncludePaths([]string{*watchDir})); err != nil {
		log.Println("Error: ", err)
	}
	if err := comp.Run(); err != nil {
		log.Println("Error: ", err)
	}
	if err := ioutil.WriteFile(*outputFile, cssbufer.Bytes(), 0644); err != nil {
		log.Println("Error: ", err)
	}
}
