package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	buildDir     string
	pagesDir     string
	templatesDir string
}

var c config

func main() {
	c = config{
		buildDir:     "site/build",
		pagesDir:     "site/pages",
		templatesDir: "site/templates",
	}

	build(c.pagesDir, c.buildDir)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		cleanUpDir(c.buildDir)
		build(c.pagesDir, c.buildDir)
		http.FileServer(http.Dir(c.buildDir)).ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cleanUpDir(path string) {
	checkError(os.RemoveAll(path))
	checkError(os.Mkdir(path, 0755))
}

func build(srcDir string, dstDir string) {
	dir, err := os.ReadDir(srcDir)
	checkError(err)
	cleanUpDir(dstDir)
	for _, file := range dir {
		srcFile := filepath.Join(srcDir, file.Name())
		dstFile := dstDir

		if file.IsDir() {
			build(srcFile, filepath.Join(dstDir, file.Name()))
			continue
		}

		fname := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		if fname != "index" {
			dstFile = filepath.Join(dstDir, fname)
		}

		tpl := template.Must(template.ParseFiles(srcFile))
		tpl = template.Must(template.Must(tpl.Clone()).ParseGlob(filepath.Join(c.templatesDir, "*.html")))
		checkError(os.MkdirAll(dstFile, 0775))
		index, err := os.Create(filepath.Join(dstFile, "index.html"))
		checkError(err)
		tpl.Execute(index, nil)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
