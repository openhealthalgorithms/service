package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"

	"github.com/openhealthalgorithms/service/pkg/algorithms"
	heartsAlg "github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"github.com/openhealthalgorithms/service/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	cpuprofile = "./oha-service-cpu-prof.prof"
	memprofile = "./oha-service-mem-prof.prof"

	appName    = "oha-service"
	appVersion = "v0.1"
	appCommit  = "0000000"
)

func main() {

	// To be able to exit properly in case of error.
	var err error
	// Deferred calls are run in LIFO order.
	// os.Exit does not run any other deferred calls.
	// This way allows us to exit the app releasing resources properly.
	defer func() {
		if err != nil && err != context.Canceled {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// Create an instance of the cli.App to process commands and options (flags).
	// For each mode we have a command.
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "Provides an API service for Open Health Algorithms"
	app.Version = appVersion
	// Create Metadata to carry commit version
	app.Metadata = make(map[string]interface{})
	app.Metadata["commit"] = appCommit
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Open Health Algorithms",
			Email: "contact@openhealthalgorithms.org",
		},
	}

	// The list of global flags.
	app.Flags = []cli.Flag{
		// Debug mode sets level for logger to debug.
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug mode makes output more verbose. Default - off",
		},
		// Algorithm name
		cli.StringFlag{
			Name:  "algorithm",
			Usage: "Algorithm to use. REQUIRED.",
			Value: "who",
		},
		// Param for algorithm
		cli.StringFlag{
			Name:  "param",
			Usage: "Param for the algorithm. REQUIRED.",
			Value: "",
		},
		// Local mode makes agent send data to localhost.
		cli.BoolFlag{
			Name:  "grace",
			Usage: "Grace mode enables graceful shutdown. Default - off",
		},
	}

	// The list of commands.
	app.Commands = []cli.Command{
		// Run command.
		cli.Command{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the algorithm",
			Action: func(c *cli.Context) error {
				return setupAndRun(c)
			},
		},
		// Daemon command.
		cli.Command{
			Name:    "daemon",
			Aliases: []string{"d"},
			Usage:   "run algorithm in daemon mode",
			Action: func(c *cli.Context) error {
				return setupAndRun(c)
			},
		},
	}

	// Run the app
	// The error is handled by the deferred function
	err = app.Run(os.Args)
}

func setupAndRun(cliCtx *cli.Context) error {
	var err error
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	var algorithmName string
	var param string
	var listAlgorithms bool
	var showConfig bool
	var cpuProf bool
	var memProf bool

	flag.StringVar(&algorithmName, "algorithm", "HeartsAlgorithm", "algorithm name")
	flag.StringVar(&param, "param", "{gender:male,age:40,systolic1:120,systolic2:140,cholesterol:8,cholesterolUnit:mmol,smoker:true,diabetic:true,region:searb}", "param for algorithm")
	flag.BoolVar(&listAlgorithms, "list", false, "list available algorithms")
	flag.BoolVar(&showConfig, "showconfig", false, "show config for algorithms")
	flag.BoolVar(&cpuProf, "cpuprofile", false, "enable cpu profiling")
	flag.BoolVar(&memProf, "memprofile", false, "enable mem profiling")
	flag.Parse()

	if cpuProf {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	algorithmsMap := map[string]interface{}{
		"hearts": heartsAlg.New(),
	}

	var algorithmsList []string
	for k := range algorithmsMap {
		algorithmsList = append(algorithmsList, k)
	}
	sort.Strings(algorithmsList)

	if listAlgorithms {
		var buf bytes.Buffer
		for _, i := range algorithmsList {
			fmt.Fprintf(&buf, "%s\n", i)
		}
		log.Printf("available algorithms:\n%s", buf.String())
		os.Exit(0)
	}

	algorithmRaw, ok := algorithmsMap[algorithmName]
	if !ok {
		log.Errorf("algorithm %s not found", algorithmName)
		os.Exit(1)
	}

	algorithm, ok := algorithmRaw.(algorithms.Algorithmer)
	if !ok {
		log.Errorf("algorithm %s doesn't implement Algorithmer interface", algorithmName)
		os.Exit(1)
	}

	v := types.NewValuesCtx()
	v.Params.Set("params", param)
	ctx := context.WithValue(context.Background(), types.KeyValuesCtx, &v)
	err = algorithm.Get(ctx)

	if err != nil {
		log.Fatal("error: ", err)
	}

	algorithmOut, _ := algorithm.Output()
	j, _ := json.MarshalIndent(algorithmOut, "", "  ")
	log.Info("algorithm output\n", string(j))

	if memProf {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

	return err
}
