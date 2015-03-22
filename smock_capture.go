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
	"os/exec"
	"syscall"
	"path"

)

//=========================================================================================
type CommandResponse struct {
	Stdout string `json:"Stdout"`
	Stderr string `json:"Stderr"`
	Exitcode int `json:"Exitcode"`
	Delay int `json:"Delay"`
	Stdin string
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
var namespace string

func main() {
	// two different ways to execute an http.GET
	//
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Author = "Corey Osman"
	app.Email = "corey@logicminds.biz"
	app.Name = "smock"
	app.Usage = "smock usage"
	app.Flags = []cli.Flag {

		cli.StringFlag{
			Name:      "template-path",
			Usage:     "Where to store templates during command captures",
			Value:     "smock-templates",
			EnvVar:    "SMOCK_TEMPLATE_PATH",
		},
		cli.StringFlag{
			Name:      "namespace",
			Usage:     "Pass in a namespace to be used when running or capturing commands.\n\tThis helps separate sets of commands that might have different use cases\n\tie. biz/logicminds/rubyipmi",
			EnvVar:    "SMOCK_COMMAND_NAMESPACE",
		},
		cli.BoolFlag{
			Name:      "query-first",
			Usage:     "Check if the command exists before capturing",
			EnvVar:    "SMOCK_QUERY_FIRST",
		},

	}
	app.Commands = []cli.Command{

		{
			Name:      "commandhash",
			Usage:     "Returns the hash of the args",
			Action: func(c *cli.Context) {
				println(generateCommandHash(c.Args()))
				os.Exit(0)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		// there can be two ways of entering commands
		// ie. echo hello  or
		// ie. 'echo hello'
		var full_cmd []string
		if len(c.Args()) > 1 {
			full_cmd = c.Args()[:]
		} else if len(c.Args()) == 1 {
			full_cmd = strings.Split(c.Args()[0], " ")
		} else {
			os.Exit(0)
		}
		namespace = c.String("namespace")
		command_name := full_cmd[0]
		command_args := full_cmd[0:]
		// if forced, check to see if the command already exists, then exit
		if c.Bool("query-first") {
			args_hash := generateCommandHash(command_args)
			endpoint := generateEndpoint(command_name, args_hash)
			cmd := doGet(endpoint)
			if cmd.Stdin != "invalid_command"  {
				log.Printf("Command already exists, not capturing")
				printToConsole(cmd)
			}
		}
		resp := captureCommand(command_args)
		cmd_map := generateResponseMap(resp)
		template_path := c.String("template-path")
		mergeToCaptureFile(fmt.Sprintf("%s.tmpl", path.Join(template_path, namespace,command_name)), cmd_map)
		printToConsole(resp)

	}
	file, err := os.OpenFile("/tmp/smock.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", "file.txt", ":", err)
	}
	log.SetOutput(file)
	app.Run(os.Args)
}

// generates a sha1 hash based on cmd_stdin
// generally this should be the command and args
func generateCommandHash(cmd_stdin []string) string {
	full_command := strings.Join(cmd_stdin, " ")
	data := []byte(fmt.Sprintf("%U",full_command))
	hash := fmt.Sprintf("%x", sha1.Sum(data))
	return hash
}
// Creates the given path, if path has an extension it gets the base name of the file and generates that path
func createTemplatePath(filepath string) {
	filepath = path.Dir(filepath)
	err := os.MkdirAll(filepath, 0744)
	check(err)
}
// dumps the cmd and exit status to the console then exits
func printToConsole(command CommandResponse) {
	if command.Exitcode == 0 {
		fmt.Print(command.Stdout)
	} else {
		fmt.Print(command.Stderr)
	}
	os.Exit(command.Exitcode)
}
// write the map output to a file
// if a file already exists with a map inside, read the contents first and then merge the contents and overwrite the file
func mergeToCaptureFile(filepath string, m map[string]CommandResponse) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// file does not exist, no need to merge
		createTemplatePath(filepath)
	} else {
		// file already exists we need to read in the data and then merge the two together
		prev := generateResponseMapFromFile(filepath)
		for k,v := range prev {
			m[k] = v
		}
	}
	b, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		log.Fatalf("error encoding command %v", err)
	}
	err = ioutil.WriteFile(filepath, b, 0644)
	check(err)
	log.Printf("Wrote to file: %s", filepath)
}
func getCommand(command_name string, args_hash string) {
	// path is normally reserved for File related things, but since its unix this will also work on url schemes
	endpoint := generateEndpoint(command_name, args_hash)
	cmd := doGet(endpoint)
	printToConsole(cmd)
}
// run the command on the os and capture the output, return a CommandResponse Object
func captureCommand(args []string) CommandResponse {
	command := args[0]
	args = args[1:]
	log.Printf("Command: %s %s", command, args)
	var exitcode = 0
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// There is no plattform independent way to retrieve
			// the exit code, but the following will work on Unix
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// for some odd reason this is output exit status 1
				exitcode = status.ExitStatus()
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	//usertime := cmd.ProcessState.UserTime()
	//fmt.Printf("Milli [%v]", usertime.Seconds())

	log.Printf("Exit Status: %d", exitcode )
	// store the original command call in stdin
	stdin := append([]string{command}, args...)
	command_response := CommandResponse{Stdout: stdout.String(),
		Stderr: stderr.String(),
		Exitcode: exitcode,
		Delay:0,
		Stdin: strings.Join(stdin, " "),
	}
	return command_response
}
func generateResponseMap(command CommandResponse) map[string]CommandResponse {
	var m = make(map[string]CommandResponse)
	hash := generateCommandHash(strings.Split(command.Stdin, " "))
	m[hash] = command
	return m
}
func generateResponseMapFromFile(filepath string) map[string]CommandResponse {
	var m = map[string]CommandResponse{}
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	if err := json.Unmarshal(dat,&m) ; err != nil {
		log.Fatalf("Cannot read file %s", filepath)
	}
	return m
}
// generates a valid http endpoint
func generateEndpoint(command_name string, args_hash string) string {
	return fmt.Sprintf("http://localhost:3001/%s", path.Join("commands",command_name, args_hash))
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
		endpoint := generateEndpoint(command_name, args_hash)
		cmd := doGet(endpoint)
		if cmd.Exitcode == 0 {
			fmt.Println(cmd.Stdout)
		} else {
			fmt.Println(cmd.Stderr)
		}
		fmt.Println(cmd.Exitcode)
		fmt.Print("SmockShell> ")
	}
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Add("X-Smock-Namespace", namespace)
	response, err := client.Do(req)
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



