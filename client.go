package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"path"
	"os/exec"
	"syscall"
)

func getCommandList() []string {
	// path is normally reserved for File related things, but since its unix this will also work on url schemes
	endpoint := generateEndpoint("","")
	jsondata := Get(endpoint)
	commands := []string{}
	// render json to object

	de_err := json.Unmarshal(jsondata, &commands)
	check(de_err)
	return commands
}
func Get(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Add("X-Shmock-Namespace", namespace)
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
		return contents
	}
	return nil
}
func doGet(url string) CommandResponse {
	contents := Get(url)
	cmd := renderJson(contents)
	return cmd
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

// generates a valid http endpoint
func generateEndpoint(command_name string, args_hash string) string {
	return fmt.Sprintf("http://localhost:3001/%s", path.Join("commands",command_name, args_hash))
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
		fmt.Println(command.Stderr)
	}
	os.Exit(command.Exitcode)
}
