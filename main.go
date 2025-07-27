package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type config struct {
	buildDir     string
	pagesDir     string
	templatesDir string
}

var c config

func main() {
	c = config{
		buildDir:     "example/build/",
		pagesDir:     "example/pages/",
		templatesDir: "example/templates/",
	}

	checkError(os.RemoveAll(c.buildDir))
	checkError(os.Mkdir(c.buildDir, 0755))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		build()
		http.FileServer(http.Dir(c.buildDir)).ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func build() {
	dir, err := os.ReadDir(c.pagesDir)
	checkError(err)
	for _, file := range dir {
		fname := file.Name()
		pname := fname[0 : len(fname)-5]
		if pname == "index" {
			pname = ""
		}
		tpl := template.Must(template.ParseFiles(c.pagesDir + fname))
		tpl = template.Must(template.Must(tpl.Clone()).ParseGlob(c.templatesDir + "*.html"))
		os.Mkdir(c.buildDir+pname, 0775)
		index, err := os.Create(c.buildDir + pname + "/index.html")
		checkError(err)
		tpl.Execute(index, nil)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
