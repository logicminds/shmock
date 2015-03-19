package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"
	"crypto/sha1"
  	"github.com/codegangsta/cli"
	"encoding/json"

)

//=========================================================================================
type CommandResponse struct {
	Stdout string `json:"Stdout"`
	Stderr string `json:"Stderr"`
	Exitcode int `json:"Exitcode"`
}

func main() {
	// two different ways to execute an http.GET
	//
	app := cli.NewApp()
	app.Name = "smock"
	app.Usage = "smock usage"
	app.Commands = []cli.Command{
	  {
	    Name:      "get",
	    ShortName: "g",
	    Usage:     "get a command response",
	    Action: func(c *cli.Context) {
	      println("Get a command response: ", c.Args().First())
	    },
  	},
 	}
	command_name := os.Args[1]
	args_hash := generateCommandHash(os.Args[1:])
	app.Action = func(c *cli.Context) {
		endpoint := fmt.Sprintf("http://localhost:3001/commands/%s/%s", command_name, args_hash)
		doGet(endpoint)
	}
	app.Run(os.Args)
}
func generateCommandHash(cmd_stdin []string) string {
  data := []byte(fmt.Sprintf("%x",cmd_stdin))
  hash := fmt.Sprintf("%x", sha1.Sum(data))
  fmt.Println(hash)
  return hash
}

func renderJson(jsondata []byte) CommandResponse {
	res := &CommandResponse{}
	// render json to object

	de_err := json.Unmarshal(jsondata, &res)

	if de_err != nil {
		//return nil if not found
		fmt.Println(string(jsondata))
		//panic(de_err)
		res.Exitcode = 1
		res.Stdout = "Invalid Json in Command Response"
	}
	return *res
}

func doGet(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		// Check for 404 error, return command not found
		if err != nil {

			log.Fatal(err)
		}
		cmd := renderJson(contents)
		fmt.Println(cmd.Stdout)
		os.Exit(cmd.Exitcode)
	}
}
func doPut(url string) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader("<golang>really</golang>"))
	request.SetBasicAuth("admin", "admin")
	request.ContentLength = 23
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("The calculated length is:", len(string(contents)), "for the url:", url)
		fmt.Println("   ", response.StatusCode)
		hdr := response.Header
		for key, value := range hdr {
			fmt.Println("   ", key, ":", value)
		}
		fmt.Println(contents)
	}
}



