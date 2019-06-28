package cmd

import (
	"github.com/nuls-io/nuls-v2-package/util"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v2"
)

var commands = []*cli.Command{}
var dir,configPath string

func init() {
	commands = append(commands, listModuleCommand)
	commands = append(commands, addModuleCommand)
	commands = append(commands, removeModuleCommand)

	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	configPath = dir + Separator + "package.ncf"
	exist, _ := util.PathExists(configPath)
	if !exist {
		// 复制默认配置文件
		util.CopyFile(dir + Separator + "package-base.ncf", configPath)
	}
}

func GetCommands() []*cli.Command {
	return commands
}

var listModuleCommand = &cli.Command {
	Name:      "show",
	Usage:     "查看配置的待打包模块列表",
	ArgsUsage: "",
	Action:    listModule,
	Description: "",
}

var addModuleCommand = &cli.Command{
	Name:  "add",
	Usage: "添加一个打包模块",
	ArgsUsage: "add [ModuleName]",
	Action:    addModule,
	Description: "",
}

var removeModuleCommand = &cli.Command{
	Name:  "remove",
	Usage: "移除一个打包模块",
	ArgsUsage: "remove [ModuleName]",
	Action:    removeModule,
	Description: "",
}

func listModule(ctx *cli.Context) error {

	configMap := readConfigFile(configPath)

	for k, v := range configMap {
		log.Println(k, " : ", v)
	}
	return nil
}

func addModule(ctx *cli.Context) error {

	if len(os.Args) < 3 {
		log.Println("缺少参数, 使用方法为./package add XXX, XXX为模块名称")
		return nil
	}

	cfg, _ := LoadConfigFile(configPath)

	moduleName := os.Args[2]
	v,_ := cfg.String("package", moduleName)
	if v == "" {
		log.Println("参数错误，没有名为[", moduleName, "]的模块")
		return nil
	}
	cfg.AddOption("package", moduleName, "1")

	cfg.WriteFile(configPath, 0644, "")
	log.Println("success")
	return nil
}

func removeModule(ctx *cli.Context) error {

	if len(os.Args) < 3 {
		log.Println("缺少参数, 使用方法为./package remove XXX, XXX为模块名称")
		return nil
	}

	cfg, _ := LoadConfigFile(configPath)

	moduleName := os.Args[2]
	v,_ := cfg.String("package", moduleName)
	if v == "" {
		log.Println("参数错误，没有名为[", moduleName, "]的模块")
		return nil
	}
	cfg.AddOption("package", moduleName, "0")

	cfg.WriteFile(configPath, 0644, "")
	log.Println("success")
	return nil
}