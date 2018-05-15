package main

import (
	"fmt"
	"os"
	"flag"
	"runtime"
	"github.com/xxlixin1993/LCS/graceful_exit"
	"github.com/xxlixin1993/LCS/configure"
	"github.com/xxlixin1993/LCS/logging"
)

const (
	KVersion = "0.0.1"
)

func main() {
	fmt.Println("build project")
}

// Initialize framework
func initFrame() {
	// Parsing configuration environment
	runMode := flag.String("m", "local", "Use -m <config mode>")
	configFile := flag.String("c", "./conf/app.ini", "use -c <config file>")
	version := flag.Bool("v", false, "Use -v <current version>")
	flag.Parse()

	// Show version
	if *version {
		fmt.Println("Version", KVersion, runtime.GOOS+"/"+runtime.GOARCH)
		os.Exit(0)
	}

	// Initialize exitList
	utils.InitExitList()

	// Initialize configure
	configErr := configure.InitConfig(*configFile, *runMode)
	if configErr != nil {
		fmt.Printf("Initialize Configure error : %s", configErr)
		os.Exit(configure.KInitConfigError)
	}

	// Initialize log
	logErr := logging.InitLog()
	if logErr != nil {
		fmt.Printf("Initialize log error : %s", logErr)
		os.Exit(configure.KInitLogError)
	}

	logging.Trace("Initialized frame")
}
