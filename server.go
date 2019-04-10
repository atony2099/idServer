/*
 * @Author: atony2099
 * @Date: 2019-04-10 02:59:49
 * @Last Modified by: atony2099
 * @Last Modified time: 2019-04-10 12:37:26
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/atony2099/idServer/core"
)

var (
	configFile string
)

func initCmd() {
	flag.StringVar(&configFile, "config", "./config.json", "where config.json is")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initEnv()
	initCmd()

	var err error

	if err = core.LoadConf(configFile); err != nil {
		goto ERROR
	}
	if err = core.InitMysql(); err != nil {
		goto ERROR
	}

	if err = core.InitAlloc(); err != nil {
		goto ERROR
	}

	if err = core.StartServer(); err != nil {
		goto ERROR
	}

	os.Exit(0)
ERROR:
	fmt.Println(err)
	os.Exit(-1)
}
