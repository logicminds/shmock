package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/stretchr/graceful"
  "github.com/unrolled/render"
  "text/template"
  "time"
  "encoding/json"
  "crypto/sha1"
  "bytes"
  //"os"

)
type CommandEnv struct {
  Home   string
}
type CommandResponse struct {
  Stdout string
  Stderr string
  Exitcode int
}

// Generate a unique hash for the given environment
func mainHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
}

func generateCommandHash(cmd_stdin string) string {
  data := []byte(fmt.Sprintf("%x",cmd_stdin))
  return fmt.Sprintf("%x", sha1.Sum(data))
}

func renderJson(file string, context CommandEnv) []byte {
  // Check for template file and make sure template exists, otherwise render blank template
  var doc bytes.Buffer 
  tmpl, err := template.ParseFiles( file)
  if err != nil { panic(err) }
  err = tmpl.Execute(&doc, context)
  if err != nil { panic(err) }
  return doc.Bytes()
}

func getCommandResponse(cmd_hash string, templatefile string, context CommandEnv) string {
  jsondata := renderJson(templatefile, context)
  //return jsondata
  var objmap map[string]CommandResponse
  // // Get json object
  err := json.Unmarshal(jsondata, &objmap)
  if err != nil {

  }
  fmt.Printf("Looking at %s", cmd_hash)
  cmd_json, de_err := json.Marshal(objmap[cmd_hash])
  if de_err != nil {
    // return nil if not found
  }
  return string(cmd_json)

}

func commandHandler(w http.ResponseWriter, r *http.Request) {
  commandname := mux.Vars(r)["command"]
  response_id := mux.Vars(r)["id"]
  fmt.Printf("%s", r.Body)
  template_dir := "templates"
  template_suffix := "tmpl"
  templatefile := fmt.Sprintf("%s/%s.%s", template_dir,commandname, template_suffix)
  cmd_env := CommandEnv{"/home/user1"}
  // throw 404 if not found
  fmt.Fprintf(w, getCommandResponse(response_id, templatefile, cmd_env))
}
// func listCommandsHandler(w http.ResponseWriter, r *http.Request) {
//   commandname := mux.Vars(r)["command"]
//   //response_id := mux.Vars(r)["id"]
//   fmt.Printf("%s", r.Body)
//   template_dir := "templates"
//   template_suffix := "tmpl"
//   templatefile := fmt.Sprintf("%s/%s.%s", template_dir,commandname, template_suffix)
//   sweaters := Inventory{"wool", 17}
//   // Check for template file and make sure template exists, otherwise render blank template
//   tmpl, err := template.ParseFiles( templatefile)
//   if err != nil { panic(err) }
//   err = tmpl.Execute(w, sweaters)
//   if err != nil { panic(err) }
// }
// func showCommandResponseHandler(w http.ResponseWriter, r *http.Request) {
//   commandname := mux.Vars(r)["command"]
//   //response_id := mux.Vars(r)["id"]
//   fmt.Printf("%s", r.Body)
//   template_dir := "templates"
//   template_suffix := "tmpl"
//   templatefile := fmt.Sprintf("%s/%s.%s", template_dir,commandname, template_suffix)
//   sweaters := Inventory{"wool", 17}
//   // Check for template file and make sure template exists, otherwise render blank template
//   tmpl, err := template.ParseFiles( templatefile)
//   if err != nil { panic(err) }
//   err = tmpl.Execute(w, sweaters)
//   if err != nil { panic(err) }
// }


func NewRenderer() (*render.Render) {

  r := render.New( render.Options{
    Directory: "templates", // Specify what path to load the templates from.
    Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
    Extensions: []string{".tmpl", ".html", ".json"}, // Specify extensions to load for templates.
    Delims: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
    Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
    IndentJSON: true, // Output human readable JSON.
    IndentXML: true, // Output human readable XML.
    //PrefixJSON: []byte(")]}',\n"), // Prefixes JSON responses with the given bytes.
    PrefixXML: []byte("<?xml version='1.0' encoding='UTF-8'?>"), // Prefixes XML responses with the given bytes.
    HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html".
    IsDevelopment: true, // Render will now recompile the templates on every HTML response.
  } )
  return r
}

func main() {
  router := mux.NewRouter()

  router.HandleFunc("/", mainHandler) 

  router.HandleFunc("/commands/{command}/{id:.+}", commandHandler)
  //router.HandleFunc("/commands/{command}", showCommandResponsesHandler)
  //router.HandleFunc("/commands", listCommandsHandler)
 
  n := negroni.Classic()

  n.UseHandler(router)
  graceful.Run(":3001",10*time.Second,n)

}


