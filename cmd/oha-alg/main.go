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

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/openhealthalgorithms/service/pkg/riskmodels"
	freRM "github.com/openhealthalgorithms/service/pkg/riskmodels/framingham"
	whoCvdRM "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
	"github.com/openhealthalgorithms/service/pkg/types"
)

var (
	cpuprofile = "./oha-service-cpu-prof.prof"
	memprofile = "./oha-service-mem-prof.prof"

	appName    = "oha-algorithms"
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
		// RiskModel name
		cli.StringFlag{
			Name:  "riskmodel",
			Usage: "Risk Model to use. REQUIRED.",
			Value: "whocvd",
		},
		// Param for algorithm/risk model
		cli.StringFlag{
			Name:  "param",
			Usage: "Param for the risk model. REQUIRED.",
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

	var riskModelName string
	var param string
	var listRiskModels bool
	var showConfig bool
	var cpuProf bool
	var memProf bool
	var debug bool

	flag.StringVar(&riskModelName, "riskmodel", "WhoCVDRiskModel", "risk model name")
	flag.StringVar(&param, "param", "gender:male,age:40,systolic1:120,systolic2:140,cholesterol:8,cholesterolUnit:mmol,smoker:true,diabetic:true,region:searb", "param for riskModel")
	flag.BoolVar(&listRiskModels, "list", false, "list available riskModels")
	flag.BoolVar(&showConfig, "showconfig", false, "show config for riskModels")
	flag.BoolVar(&cpuProf, "cpuprofile", false, "enable cpu profiling")
	flag.BoolVar(&memProf, "memprofile", false, "enable mem profiling")
	flag.BoolVar(&debug, "debug", false, "debug flag")
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

	riskModelsMap := map[string]interface{}{
		"fre":    freRM.New(),
		"whocvd": whoCvdRM.New(),
	}

	var riskModelsList []string
	for k := range riskModelsMap {
		riskModelsList = append(riskModelsList, k)
	}
	sort.Strings(riskModelsList)

	if listRiskModels {
		var buf bytes.Buffer
		for _, i := range riskModelsList {
			fmt.Fprintf(&buf, "%s\n", i)
		}
		log.Printf("available riskModels:\n%s", buf.String())
		os.Exit(0)
	}

	riskModelRaw, ok := riskModelsMap[riskModelName]
	if !ok {
		log.Errorf("risk model %s not found", riskModelName)
		os.Exit(1)
	}

	riskModel, ok := riskModelRaw.(riskmodels.RiskModeler)
	if !ok {
		log.Errorf("risk model %s doesn't implement RiskModeler interface", riskModelName)
		os.Exit(1)
	}

	v := types.NewValuesCtx()
	v.Params.Set("params", param)
	// if cliCtx.GlobalBool("debug") {
	if debug {
		v.Params.Set("debug", "true")
	}
	ctx := context.WithValue(context.Background(), types.KeyValuesCtx, &v)
	err = riskModel.Get(ctx)

	if err != nil {
		log.Fatal("error: ", err)
	}

	riskModelOut, _ := riskModel.Output()
	j, _ := json.MarshalIndent(riskModelOut, "", "  ")
	log.Info("risk model output\n", string(j))

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
