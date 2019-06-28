package cmd

import (
	"github.com/nuls-io/nuls-v2-package/util"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	Separator = string(os.PathSeparator)

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
	//是否跳过执行mvn打包命令从新打包各模块
	doSkipPackage = false
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
	//是否压缩输出目录
	zip = false
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

	doSkipPackage = ctx.Bool("i")

	// Get the current project directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	// Determine if the current directory is in the source directory
	log.Println("dir : ", dir)
	isInSource := strings.LastIndex(dir, Separator + "build") != -1

	if !isInSource {
		log.Println("Make sure the packager is in the correct directory [nuls-v2/build]")
		os.Exit(0)
	}
	projectDir = dir[:strings.LastIndex(dir, Separator)]
	log.Println("the project home is ：", projectDir)

	// the output dir
	buildPath = projectDir + Separator + "build"
	packageConfig = buildPath + Separator + "package.ncf"

	outputPath = ctx.String("o")
	if outputPath == "" {
		outputPath = buildPath + Separator + "output"
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

	if ctx.Bool("z") {
		zip = true
	}

	// Packaged name
	nulsWalletName = ctx.String("n")
	if nulsWalletName == "" {
		nulsWalletName = "NULS-Walltet-linux64"
	}
	nulsWalletPath = outputPath + Separator + nulsWalletName
	modulesPath = nulsWalletPath + Separator + "Modules"

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
	tempDir := outputPath + Separator + "tmp" + Separator
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
	util.CopyDir(tempDir + "Release" + Separator + "Modules", nulsWalletPath + Separator + "Modules")
	util.CopyDir(tempDir + "Release" + Separator + "Libraries", nulsWalletPath + Separator + "Libraries")

	if osVersion == "Linux" {
		//str := "
		//	#!/bin/bash
		//	cd \`dirname $0\`
		//	"
	}

	// Delete the extracted file
	os.RemoveAll(tempDir + Separator + "Release")
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

	if doSkipPackage {
		return
	}

	log.Println("start maven package , exec mvn package")

	sysType := runtime.GOOS

	cmd := "bash"
	args := []string {"-c", "cd .. && mvn clean package -Dmaven.test.skip=true"}
	if sysType == "windows" {
		// windows系统
		cmd = "cmd.exe"
		args = []string {"/C", "cd .. && mvn clean package -Dmaven.test.skip=true"}
	}
	result := util.Command(cmd, args)
	log.Println(result)
	if !result {
		os.Exit(-1)
	}
}

func doCopy() {

	// 读取当前的打包配置
	exist, _ := util.PathExists(packageConfig)
	if !exist {
		// 复制默认配置文件
		util.CopyFile(buildPath + Separator + "package-base.ncf", packageConfig)
	}

	configMap := readConfigFile(packageConfig)

	// 根据配置，把需要打包的模块进行复制
	for k, v := range configMap {
		if v == "0" {
			continue
		}
		if v == "mykernel" && !doMock {
			continue
		}
		path, version := findModulePath(projectDir, k)
		if path != "" && version != "" {
			doCopyModule(k, path, version)
		}
	}
}

func doCopyModule(moduleName string, modulePath string, moduleVersion string) {
	log.Println("do copy Module : ", moduleName, modulePath, moduleVersion)

	// 定义目标目录
	destDir := modulesPath + Separator + "Nuls" + Separator + moduleName + Separator + moduleVersion + Separator
	// 公共jar目标路径
	libsDir := modulesPath + Separator + "Nuls" + Separator + "libs"

	// 创建模块的lib目录
	libDir := destDir + "lib"
	log.Println("create lib dir of ", libDir)
	os.MkdirAll(libDir, 0755)

	// 遍历mvn package之后的目录
	targetDir := modulePath + Separator + "target"
	rd, _ := ioutil.ReadDir(targetDir)
	for _, fi := range rd {
		if fi.IsDir() {
			// 复制依赖的jar包，并生成依赖列表文件
			if fi.Name() == "libs" {
				util.CopyDir(targetDir + Separator + fi.Name(), libsDir)

				dependentFile := destDir + "dependent.conf"
				dependentFileContent := ""
				libRd, _ := ioutil.ReadDir(targetDir + Separator + fi.Name())
				for _, libFi := range libRd {
					if dependentFileContent != "" {
						dependentFileContent += "\n"
					}
					dependentFileContent += libFi.Name()
				}
				ioutil.WriteFile(dependentFile, []byte(dependentFileContent), 0666)
			}
			continue
		}

		// 复制jar文件到目标目录
		fileSuffix := path.Ext(fi.Name())
		if fileSuffix == ".jar" {
			util.CopyFile(targetDir + Separator + fi.Name(), destDir + Separator + moduleName + "-" + moduleVersion + fileSuffix)
			continue
		}
	}

	// 拷贝启动停止脚本
	copyScripts(moduleName, modulePath, moduleVersion, destDir)

	// 合并并生成配置文件
	mergeConfig(modulePath, destDir)
}

func mergeConfig(modulePath string, destDir string) {
	defCfg, _ := LoadConfigFile(buildPath + Separator + "module-prod.ncf")

	cfg, _ := LoadConfigFile(modulePath + Separator + "module.ncf")

	for _, selectName := range defCfg.Sections() {
		opts, err := defCfg.SectionOptions(selectName)
		if err != nil {
			continue
		}
		for _, opt := range opts {
			value1, _ := defCfg.String(selectName, opt)
			value2, _ := cfg.String(selectName, opt)
			if value1 != "" && value2 != "" && value1 != value2 {
				cfg.AddOption(selectName, opt, value1)
			}
		}
	}
	cfg.WriteFile(destDir + "module.ncf", 0644, "")
}

func copyScripts(moduleName string, modulePath string, moduleVersion string, destDir string) {
	//拷贝start, stop脚本
	//读取配置文件
	VERSION := moduleVersion
	APP_NAME := moduleName
	MAIN_CLASS := ""
	JOPT_XMS := ""
	JOPT_XMX := ""
	JOPT_MAXMETASPACESIZE := ""
	JOPT_METASPACESIZE := ""

	loadLanguage := ""
	managed := ""

	cfg, _ := LoadConfigFile(modulePath + Separator + "module.ncf")

	MAIN_CLASS,_ = cfg.String("JAVA", "MAIN_CLASS")
	JOPT_XMS,_ = cfg.String("JAVA", "JOPT_XMS")
	JOPT_XMX,_ = cfg.String("JAVA", "JOPT_XMX")
	JOPT_MAXMETASPACESIZE,_ = cfg.String("JAVA", "JOPT_MAXMETASPACESIZE")
	JOPT_METASPACESIZE,_ = cfg.String("JAVA", "JOPT_METASPACESIZE")

	managed,_ = cfg.String("Core", "Managed")
	loadLanguage,_ = cfg.String("Core", "loadLanguage")

	// 拷贝Language资源
	if loadLanguage == "1" {
		util.CopyDir(buildPath + Separator + "gen_languages", destDir + "languages")
	}

	// 追加管理配置
	if managed == "1" {
		managerFile := nulsWalletPath + Separator + ".module"
		exists,_ := util.PathExists(managerFile)
		if !exists {
			util.CreateFile(managerFile)
		}
		managerContentBytes, _ := util.ReadAllIntoMemory(managerFile)
		managerContent := string(managerContentBytes)
		if strings.Index(managerContent, moduleName) == -1 {
			if managerContent != "" {
				managerContent += "\n"
			}
			managerContent += moduleName
			ioutil.WriteFile(managerFile, []byte(managerContent), 0666)
		}
	}

	//如果模块目录下存在script文件夹，则拷贝文件夹下的内容，否则拷贝start,stop脚本
	exist, _ := util.PathExists(modulePath + Separator + "script")
	if exist {
		util.CopyDir(modulePath + Separator + "script", destDir)

		//sysType := runtime.GOOS
		//cmd := "bash"
		//args := []string {"-c", "chmod u+x " + destDir + "start"}
		//if sysType == "windows" {
		//	// windows系统
		//	cmd = "cmd.exe"
		//	args = []string {"/C", "cd .. && mvn clean package -Dmaven.test.skip=true"}
		//}
		//util.Command(cmd, args)

		return
	}

	if JOPT_XMS == "" {
		JOPT_XMS = "256"
	}
	if JOPT_XMX == "" {
		JOPT_XMX = "256"
	}
	if JOPT_MAXMETASPACESIZE == "" {
		JOPT_MAXMETASPACESIZE = "128"
	}
	if JOPT_METASPACESIZE == "" {
		JOPT_METASPACESIZE = "256"
	}

	// 替换start脚本
	// 读取start-temp内容
	startTempContent, _ := util.ReadAllIntoMemory(buildPath + Separator + "start-temp")
	newStartTempContent := string(startTempContent)
	newStartTempContent = strings.Replace(newStartTempContent, "%VERSION%", VERSION, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%APP_NAME%", APP_NAME, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%MAIN_CLASS%", MAIN_CLASS, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%JOPT_XMS%", JOPT_XMS, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%JOPT_XMX%", JOPT_XMX, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%JOPT_METASPACESIZE%", JOPT_METASPACESIZE, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%JOPT_MAXMETASPACESIZE%", JOPT_MAXMETASPACESIZE, -1)
	newStartTempContent = strings.Replace(newStartTempContent, "%JAVA_OPTS%", "", -1)

	startBatTempContent, _ := util.ReadAllIntoMemory(buildPath + Separator + "start-temp.bat")
	newStartBatTempContent := string(startBatTempContent)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%VERSION%", VERSION, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%APP_NAME%", APP_NAME, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%MAIN_CLASS%", MAIN_CLASS, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%JOPT_XMS%", JOPT_XMS, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%JOPT_XMX%", JOPT_XMX, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%JOPT_METASPACESIZE%", JOPT_METASPACESIZE, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%JOPT_MAXMETASPACESIZE%", JOPT_MAXMETASPACESIZE, -1)
	newStartBatTempContent = strings.Replace(newStartBatTempContent, "%JAVA_OPTS%", "", -1)

	ioutil.WriteFile(destDir + "start", []byte(newStartTempContent), 0666)
	ioutil.WriteFile(destDir + "start.bat", []byte(newStartBatTempContent), 0666)

	stopTempContent, _ := util.ReadAllIntoMemory(buildPath + Separator + "stop-temp")
	newStopTempContent := string(stopTempContent)
	newStopTempContent = strings.Replace(newStopTempContent, "%APP_NAME%", APP_NAME, -1)
	ioutil.WriteFile(destDir + "stop", []byte(newStopTempContent), 0666)

	stopBatTempContent, _ := util.ReadAllIntoMemory(buildPath + Separator + "stop-temp.bat")
	newStopBatTempContent := string(stopBatTempContent)
	newStopBatTempContent = strings.Replace(newStopBatTempContent, "%APP_NAME%", APP_NAME, -1)
	ioutil.WriteFile(destDir + "stop.bat", []byte(newStopBatTempContent), 0666)

	// 修改start和stop文件的权限
	sysType := runtime.GOOS
	cmd := "bash"
	args := []string {"-c", "chmod u+x " + destDir + "start"}
	if sysType == "windows" {
		// windows系统
		cmd = "cmd.exe"
		args = []string {"/C", ""}
	}
	util.Command(cmd, args)

	args = []string {"-c", "chmod u+x " + destDir + "stop"}
	util.Command(cmd, args)

	log.Println("copy scripts complete!")
}

func findModulePath(baseDir string, moduleName string) (string, string) {

	path := ""
	version := ""

	rd, _ := ioutil.ReadDir(baseDir)
	for _, fi := range rd {
		if fi.IsDir() {
			path , version = findModulePath(baseDir + Separator + fi.Name(), moduleName)
			if path != "" && version != "" {
				return path, version
			}
		} else {
			// 跳过最外层目录下的配置文件
			if baseDir == projectDir {
				continue
			}

			if fi.Name() == "module.ncf" {

				cfg, err := LoadConfigFile(baseDir + Separator + fi.Name())
				if err == nil && cfg.HasSection("JAVA") {
					options,err := cfg.SectionOptions("JAVA")
					if err == nil {
						for _,v := range options {
							if v == "APP_NAME" {
								optionValue,err := cfg.String("JAVA",v)
								if err == nil && optionValue == moduleName {
									path = baseDir
									version, _ = cfg.String("JAVA", "VERSION")
									return path, version
								}
							}
						}
					}
				}
				return path, version
			}
		}
	}
	return path, version
}

func doTar() {
	if !zip {
		return
	}

	log.Println("Start compressing the output directory")
	zipName := outputPath + Separator + nulsWalletName + ".tar.gz"

	file,err := os.Open(nulsWalletPath)
	if err != nil {
		log.Panic(err)
	}
	util.Compress([]*os.File{file}, zipName)

	log.Println("Compression is complete, generate file is : ", zipName)
}
