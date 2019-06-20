package cmd

import (
	"fmt"
	"github.com/nuls-io/nuls-v2-package/util"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	addNulstar = false
	nulstarFileName = ""
	nulstarUrl = "http://pub-readingpal.oss-cn-hangzhou.aliyuncs.com/"
	//获取参数
	//输出目录
	nulsWalletName = ""
	modulesPath = ""
	//是否马上更新代码
	doPull = false
	//是否生成mykernel模块
	doMock = false
	//更新代码的 git 分支
	gitBranch = ""
	//项目根目录
	projectDir = ""
	//打包工作目录
	buildPath = ""
	// 输出目录
	outputPath = ""
	// 编译出来钱包的路径
	nulsWalletPath = ""
	//打包配置文件
	packageConfig = ""
	//编译版本的运行系统
	osVersion = "Linux"
)

func DoPackage(ctx *cli.Context) error {

	//test()

	// 检查系统环境，是否支持打包
	check()

	// 初始化参数
	initialize(ctx)

	// 下载nulstar
	doDownload()

	// 更新代码
	doUpdateCode()

	// 打包jre
	doMvn()

	doCopy()

	doTar()

	return nil
}

func test() {
	check()

	args := []string{"/C", "cd ..", "&&", "dir"}
	err, out, errOut := util.ExecCommand("cmd", args)
	if err != nil || errOut != "" {
		log.Println(2)
		log.Println(err)
		log.Println(errOut)
		os.Exit(-1)
	}
	log.Println(out)
	os.Exit(-1)
}

// Check the packaging environment
func check() {
	// check git command
	args := []string{"--version"}
	err, _, errOut := util.ExecCommand("git", args)
	if err != nil || errOut != "" {
		log.Println("The system can't find the git command, please confirm that git is installed.")
		os.Exit(-1)
	}
	log.Println("check git : ok")

	//check mvn command
	args = []string{"-v"}
	err, _, errOut = util.ExecCommand("mvn", args)
	if err != nil || errOut != "" {
		log.Println("The system cannot find the mvn command. Please confirm that maven is installed and the environment variables are configured correctly.")
		os.Exit(-1)
	}
	log.Println("check mvn : ok")
}

// initialize params
func initialize(ctx *cli.Context) {

	// Compiled system
	if ctx.String("os") != "" {
		osVersion = ctx.String("os")
		if osVersion != "Linux" && osVersion != "MacOs" && osVersion != "Windows" {
			log.Println("Unsupported parameters，os=", osVersion)
			os.Exit(-1)
		}
	}

	// Whether to package nulstar
	addNulstar = ctx.Bool("N")
	// Nulstar download address
	nulstarFileName = ctx.String("nsn")
	nulstarUrl += nulstarFileName

	// Get the current project directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	// Determine if the current directory is in the source directory
	isInSource := dir[strings.LastIndex(dir, "/") : ] == "/build"

	if !isInSource {
		log.Println("Make sure the packager is in the correct directory [nuls-v2/build]")
		os.Exit(0)
	}
	projectDir = dir[:strings.LastIndex(dir, "/")]
	log.Println("the project home is ：", projectDir)

	// the output dir
	buildPath = projectDir + "/build"
	packageConfig = buildPath + "/package.ncf"

	outputPath = ctx.String("o")
	if outputPath == "" {
		outputPath = buildPath + "/output"
	}

	// Github configuration
	if ctx.String("b") != "" {
		gitBranch = ctx.String("b")
		doPull = true
	}

	if ctx.Bool("p") {
		gitBranch = "master"
		doPull = true
	}

	// Packaged name
	nulsWalletName = ctx.String("n")
	if nulsWalletName == "" {
		nulsWalletName = "NULS-Walltet-linux64"
	}
	nulsWalletPath = outputPath + "/" + nulsWalletName
	modulesPath = nulsWalletPath + "/Modules"

	exists,err := util.PathExists(nulsWalletPath)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	if !exists {
		os.MkdirAll(nulsWalletPath, 0755)
	}
}

func doDownload() {
	if !addNulstar {
		return
	}

	// Determine if nulstar has a different version
	tempDir := outputPath + "/tmp/"
	nulstartFilePath := tempDir + nulstarFileName
	exists, err := util.PathExists(nulstartFilePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	if !exists {
		log.Println("Start downloading nulstar···")
		// download file
		err = util.DownlaodFile(nulstarUrl, nulstartFilePath)
		if err != nil {
			log.Print(err)
			return
		}
		log.Println("Download nulstar to complete")
	}

	log.Println("Unzip nulstar")
	err = util.DeCompress(nulstartFilePath, tempDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Println("Unzip nulstar successfully")

	// copy files
	util.CopyDir(tempDir + "Release/Modules", nulsWalletPath + "/Modules")
	util.CopyDir(tempDir + "Release/Libraries", nulsWalletPath + "/Libraries")

	if osVersion == "Linux" {
		//str := "
		//	#!/bin/bash
		//	cd \`dirname $0\`
		//	"
	}

	// Delete the extracted file
	os.RemoveAll(tempDir + "/Release")
}

func doUpdateCode() {
	if !doPull {
		return
	}
	// Switch branch
	log.Println("Switch git branch to [", gitBranch, "]")
	args := []string{"checkout", gitBranch}
	err, out, errOut := util.ExecCommand("git", args)
	if err != nil || errOut != "" {
		log.Println(errOut)
		os.Exit(-1)
	}
	log.Println(out)

	// Update branch
	log.Println("Update branch to [", gitBranch, "]")
	args = []string{"pull", "origin", gitBranch}
	err, out, errOut = util.ExecCommand("git", args)
	if err != nil || errOut != "" {
		log.Println(errOut)
		os.Exit(-1)
	}
	log.Println(out)
}

func doMvn() {

	log.Println("start ")

	sysType := runtime.GOOS
	fmt.Println(sysType)

	cmd := "sh"
	args := []string {}
	if sysType == "windows" {
		// windows系统
	}
	err, out, errOut := util.ExecCommand(cmd, args)
	if err != nil || errOut != "" {
		log.Println(errOut)
		os.Exit(-1)
	}
	log.Println(out)
}

func doCopy() {

}

func doTar() {

}