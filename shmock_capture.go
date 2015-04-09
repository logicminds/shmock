package main

import (
	"fmt"
	"log"
	"strings"
	"os"
  	"github.com/codegangsta/cli"
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





