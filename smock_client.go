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
	"bytes"
	"bufio"
	"regexp"

)

//=========================================================================================
type CommandResponse struct {
	Stdout string `json:"Stdout"`
	Stderr string `json:"Stderr"`
	Exitcode int `json:"Exitcode"`
	Delay int `json:"Delay"`
}

func main() {
	// two different ways to execute an http.GET
	//
	app := cli.NewApp()
	app.Name = "smock"
	app.Usage = "smock usage"
	app.Commands = []cli.Command{
	  {
		  Name:      "list",
		  ShortName: "l",
		  Usage:     "get a list of commands",
		  Action: func(c *cli.Context) {
		    println("Get a list of commands: ", c.Args().First())
		  },
  	  },
	  {
		  Name:      "environment",
		  ShortName: "e",
		  Usage:     "use passing in environment",
		  Action: func(c *cli.Context) {
			  println("Sets an environment: ", c.Args().First())
		  },
	  },
	  {
		  Name:      "namespace",
		  ShortName: "n",
		  Usage:     "use a command namespace",
		  Action: func(c *cli.Context) {
			  println("Sets a namespace when selecting commands: ", c.Args().First())
		  },
	  },
	}
	// start a interactive shell session
	if len(os.Args) < 2 {
		runShell()
		os.Exit(0)
	}
	command_name := os.Args[1]
	args_hash := generateCommandHash(os.Args[1:])
	app.Action = func(c *cli.Context) {
		endpoint := fmt.Sprintf("http://localhost:3001/commands/%s/%s", command_name, args_hash)
		cmd := doGet(endpoint)
		fmt.Println(cmd.Stdout)
		if cmd.Stderr != "" {
			fmt.Println(cmd.Stderr)
		}
		os.Exit(cmd.Exitcode)
		//doPost(endpoint, , " ")
	}
	app.Run(os.Args)
}
func runShell() {
	r, _ := regexp.Compile("[0-9A-Za-z_-]+")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("SmockShell> ")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			break
		}
		// Catch any bad characters instead of sending them to the server
		if !r.MatchString(line) {
			fmt.Print("SmockShell> ")
			continue
		}
		// look for set namespace
		// look for set environment
		// look for list
		// look for help
		args := strings.Split(line, " ")
		command_name := args[0]
		args_hash := generateCommandHash(args)
		endpoint := fmt.Sprintf("http://localhost:3001/commands/%s/%s", command_name, args_hash)
		cmd := doGet(endpoint)
		fmt.Println(cmd.Stdout)
		if cmd.Stderr != "" {
			fmt.Println(cmd.Stderr)
		}
		fmt.Println(cmd.Exitcode)
		fmt.Print("SmockShell> ")
	}
}

func generateCommandHash(cmd_stdin []string) string {
  data := []byte(fmt.Sprintf("%x",cmd_stdin))
  hash := fmt.Sprintf("%x", sha1.Sum(data))
  //fmt.Println(hash)
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

func doGet(url string) CommandResponse {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		// Check for 404 error, return command not found
		if err != nil {

			log.Fatal(err)
		}
		cmd := renderJson(contents)
		return cmd
	}
	return CommandResponse{}
}
func doPost(url string, json_body []byte, namespace string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_body))

	if namespace != "" {
		req.Header.Set("X-Command-Namespace", namespace)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

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



