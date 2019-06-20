package config

import "gopkg.in/urfave/cli.v2"

var flags = []cli.Flag{}

func init() {
	flags = append(flags, branchFlag)
	flags = append(flags, masterBranchFlag)
	flags = append(flags, outputDirFlag)
	flags = append(flags, javaHomeFlag)
	flags = append(flags, jreFlag)
	flags = append(flags, makeTarFlag)
	flags = append(flags, skipMvnPackageFlag)
	flags = append(flags, addNulstarFlag)
	flags = append(flags, nulstarFileNameFlag)
	flags = append(flags, addModuleFlag)
	flags = append(flags, removeModuleFlag)
	flags = append(flags, osFlag)
}

func GetFlags() []cli.Flag {
	return flags
}

var branchFlag = &cli.StringFlag{
	Name:  "b",
	Usage: "-b <branch> 打包前同步最新代码 参数为同步的远程分支名称",
	Value: "",
}

var masterBranchFlag = &cli.BoolFlag{
	Name:  "p",
	Usage: "-p 打包前同步最新代码 从master分支拉取",
	Value: false,
}

var outputDirFlag = &cli.StringFlag{
	Name:  "o",
	Usage: "-o <目录>  指定输出目录",
	Value: "",
}

var javaHomeFlag = &cli.StringFlag{
	Name:  "j",
	Usage: "-j JAVA_HOME",
	Value: "",
}

var jreFlag = &cli.StringFlag{
	Name:  "J",
	Usage: "-J 输出的jvm虚拟机目录，脚本将会把这个目录复制到程序依赖中",
	Value: "",
}

var makeTarFlag = &cli.BoolFlag{
	Name:  "z",
	Usage: "-z 生成压缩包",
	Value: true,
}

var skipMvnPackageFlag = &cli.BoolFlag{
	Name:  "i",
	Usage: "-i 跳过mvn打包",
	Value: false,
}

var addNulstarFlag = &cli.BoolFlag{
	Name:  "N",
	Usage: "-N 打包时加入Nulstar模块",
	Value: true,
}

var nulstarFileNameFlag = &cli.StringFlag{
	Name:  "nsn",
	Usage: "打包时如果加入Nulstar模块，则需要指定最新版本的Nulstar名称",
	Value: "nulstar-20190529.tar.gz",
}

var addModuleFlag = &cli.StringFlag{
	Name:  "a",
	Usage: "-a 添加一个打包模块",
	Value: "",
}

var removeModuleFlag = &cli.StringFlag{
	Name:  "r",
	Usage: "-r 移除一个打包模块",
	Value: "",
}

var osFlag = &cli.StringFlag{
	Name:  "os",
	Usage: "-os <Linux|MacOs|Windows> 编译对应系统的版本",
	Value: "Linux",
}