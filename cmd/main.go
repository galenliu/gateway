package main

import (
	"flag"
	"fmt"
	"os"
	"smartassistant"
)

var (
	cfgDir        string
	logRotateDays int
	verbose       bool

	showVersion bool
	reset       bool
)

func init() {
	flag.StringVar(&cfgDir, "config", smartassistant.GetDefaultConfigDir(), "Profile directory")
	flag.IntVar(&logRotateDays, "log_rotate_days", 0, "Enables daily log rotation and keeps up to the specified days")
	flag.BoolVar(&reset, "reset", false, "reset")
	flag.BoolVar(&showVersion, "version", false, "version")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose log to file")
}

func main() {

	// 首先解析命令行参数
	flag.Parse()

	//if version command then print version
	if showVersion {
		fmt.Print(smartassistant.Version)
		return
	}

	runtimeConfig := smartassistant.NewRuntimeConfig(cfgDir, logRotateDays, verbose, reset)

	CheckError(smartassistant.Start(runtimeConfig))

}

func CheckError(err error) {
	if err != nil {
		if err != nil {
			os.Exit(1)
		}
	}
}
