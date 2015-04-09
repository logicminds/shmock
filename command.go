package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"encoding/json"

)
var namespace string

type CommandResponse struct {
	Stdout   string `json:"Stdout"`
	Stderr   string `json:"Stderr"`
	Exitcode int    `json:"Exitcode"`
	Delay    int    `json:"Delay"`
	Stdin    string
}

// generates a sha1 hash based on cmd_stdin
// generally this should be the command and args
func generateCommandHash(cmd_stdin []string) string {
	full_command := strings.Join(cmd_stdin, " ")
	data := []byte(fmt.Sprintf("%U", full_command))
	hash := fmt.Sprintf("%x", sha1.Sum(data))
	return hash
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
