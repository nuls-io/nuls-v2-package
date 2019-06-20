package cmd

import (
	"bytes"
	"log"

	"gopkg.in/urfave/cli.v2"
)

var commands = []*cli.Command{}

func init() {
	commands = append(commands, listModuleCommand)
}

func GetCommands() []*cli.Command {
	return commands
}

var listModuleCommand = &cli.Command {
	Name:      "s",
	Usage:     "查看配置的待打包模块列表",
	ArgsUsage: "",
	Action:    listModule,
	Description: "",
}

func listModule(ctx *cli.Context) error {

	if ctx.Bool("print-args") {

		var buffer bytes.Buffer
		buffer.WriteString("the args is : ")

		args := ctx.Args()
		for i := 0; i < args.Len(); i++ {
			buffer.WriteString(" \n")
			buffer.WriteString(args.Get(i))
		}

		log.Println(buffer.String())
	}

	log.Println("the message is : ", ctx.String("message"))
	return nil
}