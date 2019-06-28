package config

import (
	"log"

	"github.com/larspensjo/config"
	"github.com/pkg/errors"

	"github.com/nuls-io/nuls-v2-package/util"
)

func ReadConfigFile(file string) map[string]string {

	exists, _ := util.PathExists(file)
	if !exists {
		return nil
	}

	// 读取配置文件内容
	cfg,err := config.ReadDefault(file)

	if err != nil {
		log.Fatalf("Fail to find %v,%v", file, err)
	}

	configMap := make(map[string]string)
	if	cfg.HasSection("package") {
		options,err := cfg.SectionOptions("package")
		if err == nil {
			for _,v := range options{
				optionValue,err := cfg.String("package",v)
				if err == nil {
					configMap[v] =optionValue
				}
			}
		}
	} else {
		log.Fatalf("Fail to find section [package] in file %v", file)
	}

	return configMap
}

func LoadConfigFile(file string) (*config.Config, error) {

	exists, _ := util.PathExists(file)
	if !exists {
		return nil, errors.New("file not exists")
	}

	// 读取配置文件内容
	cfg,err := config.Read(file, "# ", "=", false, false)
	return cfg, err
}

