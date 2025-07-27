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

func main() {
	c := config{
		buildDir:     "example/build/",
		pagesDir:     "example/pages/",
		templatesDir: "example/templates/",
	}
	checkError(os.RemoveAll(c.buildDir))
	checkError(os.Mkdir(c.buildDir, 0755))

	dir, err := os.ReadDir(c.pagesDir)
	checkError(err)
	for _, file := range dir {
		fname := file.Name()
		pname := fname[0 : len(fname)-5]
		tpl := template.Must(template.ParseFiles(c.pagesDir + file.Name()))
		tpl = template.Must(template.Must(tpl.Clone()).ParseGlob(c.templatesDir + "*.tmpl"))
		os.Mkdir(c.buildDir+pname, 0775)
		index, err := os.Create(c.buildDir + pname + "/index.html")
		checkError(err)
		tpl.Execute(index, nil)
	}

	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(c.buildDir))))
}

func checkError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
