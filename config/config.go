package config

import (
	"cache/constant/env"
	clog "cache/pkg/log"
	"os"
	"path/filepath"
)

var ENV int

func Configinit(env int) {
	ENV = env
}

func ConfigPath() string {
	if ENV != env.PRODUCE_ENV {
		executablePath, err := os.Executable()
		if err != nil {
			clog.Panic("[config] ProjectPath Error : " + err.Error())
		}
		projectRoot := filepath.Dir(executablePath)
		// appYMLPath := filepath.Join(projectRoot, "config")
		return projectRoot
	}
	return "./config"
}
