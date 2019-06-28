package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"gopkg.in/urfave/cli.v2"

	"github.com/nuls-io/nuls-v2-package/cmd"
	"github.com/nuls-io/nuls-v2-package/config"
	"github.com/nuls-io/nuls-v2-package/util"
)

func main() {

	//test()
	//os.Exit(0)

	if err := setupAPP().Run(os.Args); err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
}

func test() {
	log.Println("start maven package , exec mvn package")

	sysType := runtime.GOOS

	cmd := "bash"
	args := []string {"-c", "cd .. && mvn clean package -Dmaven.test.skip=true"}
	if sysType == "windows" {
		// windows系统
		cmd = "cmd"
	}
	//cmd = "cd .. && mvn clean package -Dmaven.test.skip=true"
	res := util.Command(cmd, args)
	log.Println(res)

	//err, out, errout := util.Shellout("cd .. && mvn clean package -Dmaven.test.skip=true")
	//if err != nil {
	//	log.Printf("error: %v\n", err)
	//}
	//fmt.Println("--- stdout ---")
	//fmt.Println(out)
	//fmt.Println("--- stderr ---")
	//fmt.Println(errout)
}

func setupAPP() *cli.App {
	app := cli.App{}
	app.Name = "package"
	app.Usage = "run ./package"
	app.Version = "1.0"
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "nuls core team",
			Email: "dev@nuls.io",
		},
	}
	app.Action = startup
	app.Before = initConfig
	app.Commands = cmd.GetCommands()
	app.Flags = config.GetFlags()

	return &app
}

func initConfig(context *cli.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return nil
}

func startup(context *cli.Context) error {
	return cmd.DoPackage(context)
}