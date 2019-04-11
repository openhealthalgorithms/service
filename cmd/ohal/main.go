package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"

	"github.com/openhealthalgorithms/service/pkg/config"

	"github.com/openhealthalgorithms/service/pkg"
	"github.com/openhealthalgorithms/service/pkg/algorithms"
	heartsAlg "github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"github.com/openhealthalgorithms/service/pkg/riskmodels"
	freRM "github.com/openhealthalgorithms/service/pkg/riskmodels/framingham"
	whoCvdRM "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	cpuprofile = "./ohas-cpu-prof.prof"
	memprofile = "./ohas-mem-prof.prof"

	appName    = "ohal"
	appVersion = pkg.GetVersion()
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
			Value: "hearts",
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
			Usage: "Param file. REQUIRED.",
			Value: "sample-request.json",
		},
		cli.StringFlag{
			Name:  "project",
			Usage: "Project Name",
			Value: "default-json",
		},
		cli.StringFlag{
			Name:  "guide",
			Usage: "Guideline file. REQUIRED.",
			Value: "guideline_hearts.json",
		},
		cli.StringFlag{
			Name:  "guidecontent",
			Usage: "Guideline Content file. REQUIRED.",
			Value: "guideline_hearts_content.json",
		},
		cli.StringFlag{
			Name:  "goal",
			Usage: "Goal file. REQUIRED.",
			Value: "goal_hearts.json",
		},
		cli.StringFlag{
			Name:  "goalcontent",
			Usage: "Goal Content file. REQUIRED.",
			Value: "goal_hearts_content.json",
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
	var listAlgorithms bool
	var riskModelName string
	var listRiskModels bool
	var param string
	var guideline string
	var guidelineContent string
	var goal string
	var goalContent string
	var projectName string
	var showConfig bool
	var cpuProf bool
	var memProf bool
	var debug bool

	flag.StringVar(&algorithmName, "algorithm", "hearts", "algorithm name")
	flag.BoolVar(&listAlgorithms, "listalgorithms", false, "list available algorithms")
	flag.StringVar(&riskModelName, "riskmodel", "whocvd", "risk model name")
	flag.BoolVar(&listRiskModels, "listriskmodels", false, "list available riskModels")
	// flag.StringVar(&param, "param", "gender:male,age:40,systolic1:120,systolic2:140,cholesterol:8,cholesterolUnit:mmol,smoker:true,diabetic:true,region:searb", "param for riskModel")
	flag.StringVar(&param, "param", "sample-request.json", "param file")
	flag.StringVar(&guideline, "guide", "guideline_hearts.json", "guideline file")
	flag.StringVar(&guidelineContent, "guidecontent", "guideline_hearts_content.json", "guideline content file")
	flag.StringVar(&goal, "goal", "goals_hearts.json", "goal file")
	flag.StringVar(&goalContent, "goalcontent", "goals_hearts_content.json", "goal content file")
	flag.StringVar(&projectName, "project", "default-json", "project name")
	flag.BoolVar(&showConfig, "showconfig", false, "show config for riskModels")
	flag.BoolVar(&cpuProf, "cpuprofile", false, "enable cpu profiling")
	flag.BoolVar(&memProf, "memprofile", false, "enable mem profiling")
	flag.BoolVar(&debug, "debug", false, "debug flag")
	flag.Parse()

	currentSettings := config.CurrentSettings()
	// log.Println(currentSettings)
	if len(projectName) == 0 {
		projectName = ""
	}

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

	_, rok := riskModelRaw.(riskmodels.RiskModeler)
	if !rok {
		log.Errorf("risk model %s doesn't implement RiskModeler interface", riskModelName)
		os.Exit(1)
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

	content, err := ioutil.ReadFile(param)
	if err != nil {
		log.Errorf("param file error: %#v", err)
		os.Exit(1)
	}

	paramObj, err := tools.ParseParams(content)
	if err != nil {
		log.Errorf("param file error: %#v", err)
		os.Exit(1)
	}

	v := types.NewValuesCtx()
	v.Params.Set("params", paramObj)

	if currentSettings.CloudEnable {
		v.Params.Set("cloud", "yes")
		v.Params.Set("project", projectName)
		v.Params.Set("bucket", currentSettings.CloudBucket)
		v.Params.Set("configfile", currentSettings.CloudConfigFile)
	} else {
		v.Params.Set("cloud", "no")
		v.Params.Set("guide", currentSettings.GuidelineFile)
		v.Params.Set("guidecontent", currentSettings.GuidelineContentFile)
		v.Params.Set("goal", currentSettings.GoalFile)
		v.Params.Set("goalcontent", currentSettings.GoalContentFile)
	}

	if debug {
		v.Params.Set("debug", "true")
	}
	ctx := context.WithValue(context.Background(), types.KeyValuesCtx, &v)

	// err = riskModel.Get(ctx)
	//
	// if err != nil {
	// 	log.Fatal("error: ", err)
	// }
	//
	// riskModelOut, _ := riskModel.Output()
	// j, _ := json.MarshalIndent(riskModelOut, "", "  ")
	// if debug {
	// 	log.Info("risk model output\n", string(j))
	// }

	err = algorithm.Get(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}

	algorithmOut, _ := algorithm.Output()
	al, _ := json.MarshalIndent(algorithmOut, "", "  ")
	// dst := new(bytes.Buffer)
	// json.HTMLEscape(dst, al)
	log.Info("algorithm output\n")
	log.Println(string(al))

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
