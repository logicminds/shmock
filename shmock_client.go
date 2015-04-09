package main

import (
	"fmt"
	"log"
	"strings"
	"os"
  	"github.com/codegangsta/cli"
	"bufio"
	"regexp"
	"path"

)

//=========================================================================================

func main() {
	// two different ways to execute an http.GET
	//
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Author = "Corey Osman"
	app.Email = "corey@logicminds.biz"
	app.Name = "shmock"
	app.Usage = "shmock usage"
	app.Flags = []cli.Flag {

		cli.StringFlag{
			Name:      "namespace",
			Usage:     "Pass in a namespace to be used when running or capturing commands.\n\tThis helps separate sets of commands that might have different use cases\n\tie. biz/logicminds/rubyipmi",
			EnvVar:    "SHMOCK_COMMAND_NAMESPACE",
		},

	}
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "get a list of commands",
			Action: func(c *cli.Context) {
				namespace = c.String("namespace")
				commands := getCommandList()
				for _, file := range commands {
					fmt.Println(file)
				}
				os.Exit(0)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:      "namespace",
					Usage:     "Pass in a namespace to be used when running or capturing commands.\n\tThis helps separate sets of commands that might have different use cases\n\tie. biz/logicminds/rubyipmi",
					EnvVar:    "SHMOCK_COMMAND_NAMESPACE",
				},
			},
		},
		{
			Name:      "install",
			Usage:     "Syncs the shmocks defined in the Shmockfile",
			Action:    func(c *cli.Context) {
				sync_shmock_file(c.String("template-path"), c.String("shmockfile"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:      "shmockfile",
					Usage:     "The path to the file that contains the shmock sets",
					EnvVar:    "SHMOCK_FILE",
					Value:     "Shmockfile",
				},
				cli.StringFlag{
					Name:      "template-path",
					Usage:     "Where to store templates during command captures",
					Value:     "shmocks",
					EnvVar:    "SHMOCK_TEMPLATE_PATH",
				},
			},

		},
		{
			Name:      "commandhash",
			Usage:     "Returns the hash of the args",
			Action: func(c *cli.Context) {
				println(generateCommandHash(c.Args()))
			},
		},
		{
			Name:      "capture",
			Usage:     "run real commands and capture their output",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:      "template-path",
					Usage:     "Where to store templates during command captures",
					Value:     "shmock-templates",
					EnvVar:    "SHMOCK_TEMPLATE_PATH",
				},
				cli.StringFlag{
					Name:      "namespace",
					Usage:     "Pass in a namespace to be used when running or capturing commands.\n\tThis helps separate sets of commands that might have different use cases\n\tie. biz/logicminds/rubyipmi",
					EnvVar:    "SHMOCK_COMMAND_NAMESPACE",
				},
				cli.BoolFlag{
					Name:      "query-first",
					Usage:     "Check if the command exists before capturing",
					EnvVar:    "SHMOCK_QUERY_FIRST",
				},

			},
			Action: func(c *cli.Context) {
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
			},
		},
	}
	app.Action = func(c *cli.Context) {
		// start a interactive shell session
		namespace = c.String("namespace")
		if len(os.Args) < 2 {
			runShell()
			os.Exit(0)
		}
		i := c.NumFlags() + 1
		command_name := strings.Split(os.Args[i], " ")[0]
		command_args := os.Args[i:]
		args_hash := generateCommandHash(command_args)
		getCommand(command_name, args_hash)

	}
	file, err := os.OpenFile("/tmp/shmock.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", "file.txt", ":", err)
	}
	log.SetOutput(file)
	app.Run(os.Args)
}

func runShell() {
	r, _ := regexp.Compile("[0-9A-Za-z_-]+")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("ShmockShell> ")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			break
		}
		// Catch any bad characters instead of sending them to the server
		if !r.MatchString(line) {
			fmt.Print("ShmockShell> ")
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
		fmt.Print("ShmockShell> ")
	}
}





