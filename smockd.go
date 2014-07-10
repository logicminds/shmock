package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/unrolled/render"
  //"encoding/json"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    rd := NewRenderer()
    rd.JSON(w, http.StatusOK, map[string]string{"hello": "json"})
}
func NewRenderer() render {

  return render.New( render.Options{
    Directory: "templates", // Specify what path to load the templates from.
    Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
    Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
    Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
    Delims: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
    Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
    IndentJSON: true, // Output human readable JSON.
    IndentXML: true, // Output human readable XML.
    PrefixJSON: []byte(")]}',\n"), // Prefixes JSON responses with the given bytes.
    PrefixXML: []byte("<?xml version='1.0' encoding='UTF-8'?>"), // Prefixes XML responses with the given bytes.
    HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html".
    IsDevelopment: true, // Render will now recompile the templates on every HTML response.
  } )
}

func main() {
  router := mux.NewRouter()

  router.HandleFunc("/", mainHandler) 

  router.HandleFunc("/json", jsonHandler)
       
  n := negroni.Classic()

  n.UseHandler(router)
  n.Run(":3000")
}


