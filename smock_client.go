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
)

//=========================================================================================
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
	app.Action = func(c *cli.Context) {
		doGet("http://localhost:3001/commands/echo/1234")
	}
	app.Run(os.Args)
}
func generateCommandHash(cmd_stdin string) string {
  data := []byte(fmt.Sprintf("%x",cmd_stdin))
  return fmt.Sprintf("%x", sha1.Sum(data))
}

func doGet(url string) {
	response, err := http.Get(url)
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
		fmt.Println(string(contents))
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

