package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type config struct {
	buildDir     string
	pagesDir     string
	templatesDir string
}

type meta struct {
	Template string `yaml:"template"`
}

type props struct {
	Meta meta
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

		fext := filepath.Ext(file.Name())
		fname := strings.TrimSuffix(file.Name(), fext)
		if fname != "index" {
			dstFile = filepath.Join(dstDir, fname)
		}

		checkError(os.MkdirAll(dstFile, 0775))
		index, err := os.Create(filepath.Join(dstFile, "index.html"))
		checkError(err)

		fcontent, err := os.ReadFile(srcFile)
		checkError(err)

		var fm meta
		rest, err := frontmatter.Parse(bytes.NewReader(fcontent), &fm)

		var tpl *template.Template

		if fext == ".md" {
			var buf bytes.Buffer
			checkError(goldmark.Convert(rest, &buf))
			tpl = template.Must(template.New("page").Parse(buf.String()))
		} else {

			tpl = template.Must(template.New("page").Parse(string(rest)))
		}

		tpl = template.Must(tpl.ParseFiles(filepath.Join(c.templatesDir, fm.Template+".html")))
		checkError(tpl.ExecuteTemplate(index, fm.Template, nil))

	}
}

func checkError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
