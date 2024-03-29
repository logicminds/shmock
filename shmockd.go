package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/stretchr/graceful"
  "github.com/unrolled/render"
  "text/template"
  "path/filepath"
  "time"
  "encoding/json"
  "bytes"
  "os"
  "log"
  "strings"
  "path"

)
var template_dir string = "templates"
var template_suffix string = ".tmpl"

type CommandEnv struct {
  Home   string
}

// Generate a unique hash for the given environment
func mainHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
}

func renderTemplate(file string, context CommandEnv) []byte {
  // Check for template file and make sure template exists, otherwise render blank template
  var doc bytes.Buffer
  tmpl, err := template.ParseFiles(file)
  if err != nil {
    panic(err)
  }
  err = tmpl.Execute(&doc, context)
  if err != nil {
    panic(err)
  }
  return doc.Bytes()
}
// We render the json because its a template and could contain variables that would need to be
// rendered first.
// cmd_hash is the representation of that specific command in numeric form
// templatefile is the template file name
// context is a hash of variables that could be used to render everything together
func getCommandResponse(cmd_hash string, templatefile string, context CommandEnv) string {
  jsondata := renderTemplate(templatefile, context)
  var objmap map[string]CommandResponse
  // Get json object
  err := json.Unmarshal(jsondata, &objmap)
  if err != nil {
    // Invalid JSON
     fmt.Println(err)
  }
  //fmt.Printf("Looking at %s", cmd_hash)
  // we need to marshal again so we can just get the specific command hash
  value, ok := objmap[cmd_hash]
  if !ok {
    return renderNotFoundError()
  } else {
    cmd_json, de_err := json.Marshal(value)
    if de_err != nil {
      return renderNotFoundError()
    }
    return string(cmd_json)
  }
}
func renderNotFoundError() string {
  template_file := fmt.Sprintf("%s/%s.%s", template_dir,"internal_responses", "json")
  response_id := "invalid_command"
  cmd_env := CommandEnv{"/home/user1"}
  return getCommandResponse(response_id, template_file, cmd_env)
}
func commandHandler(w http.ResponseWriter, r *http.Request) {
  namespace := r.Header.Get("X-Shmock-Namespace")
  command_name := mux.Vars(r)["command"]
  response_id := mux.Vars(r)["id"]
  template_file := fmt.Sprintf("%s%s", path.Join(template_dir,namespace,command_name), template_suffix)
  cmd_env := CommandEnv{"/home/user1"}
  log.Printf("Looking up template: %s", template_file)
  if _, err := os.Stat(template_file); os.IsNotExist(err) {
    // once we start to pass data, we need to add the previous command called so its rendered here
    log.Printf("%s file not found", template_file)
    fmt.Fprintf(w, renderNotFoundError())
  } else {
    fmt.Fprintf(w, getCommandResponse(response_id, template_file, cmd_env))
  }
}
func listCommandsHandler(w http.ResponseWriter, r *http.Request) {
  namespace := r.Header.Get("X-Shmock-Namespace")
  template_path := path.Join(template_dir,namespace)
  path_size := len(template_path)
  prefix_size := len(template_suffix)
  fileList := []string{}
  _ = filepath.Walk(template_path, func(path string, f os.FileInfo, err error) error {
      if strings.Contains(path, template_suffix) {
        // only gets the /usr/local/sbin/echo in templates/biz/logicminds/rubyipmi/usr/local/sbin/echo.tmpl
        p := path[path_size:(len(path) - prefix_size)]
        fileList = append(fileList, p)
      }
      return nil
    })
  cmd_json, de_err := json.Marshal(fileList)
  if de_err != nil {
     panic(de_err)
  }
  fmt.Fprintf(w, string(cmd_json))

}
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
  //http://localhost:3001/commands/usr/bin/which/b7bdcd4486792996bf130fe8675f95daf905d13f
  router.HandleFunc("/commands/{command:.*}/{id:[0-9a-z]+}", commandHandler)
  // list all the commands available
  router.HandleFunc("/commands", listCommandsHandler)

  //router.HandleFunc("/commands/{command}", showCommandResponsesHandler)

  n := negroni.Classic()

  n.UseHandler(router)
  graceful.Run(":3001",10*time.Second,n)

}


