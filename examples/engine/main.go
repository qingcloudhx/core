package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"

	_ "github.com/qingcloudhx/core/data/expression/script"
	"github.com/qingcloudhx/core/engine"
)

var (
	cpuProfile    = flag.String("cpuprofile", "", "Writes CPU profile to the specified file")
	memProfile    = flag.String("memprofile", "", "Writes memory profile to the specified file")
	fileJson      = flag.String("conf", "", "app config")
	cfgJson       string
	cfgCompressed bool
)

func main() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create CPU profiling file: %v\n", err)
			os.Exit(1)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start CPU profiling: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}
	//os.Setenv("FLOGO_CONFIG_PATH", "/home/code/flowgo/core/examples/engine/flogo.json")
	if *fileJson != "" {
		data, err := ioutil.ReadFile(*fileJson)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read config file: %s\n", *fileJson)
			os.Exit(1)
		} else {
			cfgJson = string(data)
		}
	}
	cfg, err := engine.LoadAppConfig(cfgJson, cfgCompressed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	e, err := engine.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	code := engine.RunEngine(e)

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create memory profiling file: %v\n", err)
			os.Exit(1)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write memory profiling data: %v", err)
			os.Exit(1)
		}
		_ = f.Close()
	}

	os.Exit(code)
}
