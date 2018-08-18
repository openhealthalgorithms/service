package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	regions = []string{
		"AFRD", "AFRE", "AMRA", "AMRB", "AMRD", "EMRB", "EMRD", "EURA", "EURB", "EURC", "SEARB", "SEARD", "WPRA", "WPRB",
	}

	regionalColorChart RegionColorChart
)

func init() {
	populateColorChart()
}

func populateColorChart() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	regionalColorChart = make(RegionColorChart, 0)

	for _, r := range regions {
		fLocation := filepath.Join(pwd, "pkg", "riskmodels", "common", "config", "color_charts", r)
		if _, err := os.Stat(fLocation); err != nil {
			if !os.IsNotExist(err) {
				continue
			}
		}

		files, err := ioutil.ReadDir(fLocation)
		if err != nil {
			continue
		}

		colorCharts := make([]ColorChart, 0)
		for _, file := range files {
			fullFilePath := filepath.Join(fLocation, file.Name())
			names := strings.Split(file.Name(), ".")
			nameParts := strings.Split(names[0], "_")
			contents := getFileContents(fullFilePath)
			age, _ := strconv.Atoi(nameParts[4])

			colorChart := ColorChart{
				Cholesterol: nameParts[0],
				Diabetes:    nameParts[1],
				Gender:      nameParts[2],
				Smoker:      nameParts[3],
				Age:         age,
				Chart:       contents,
			}

			colorCharts = append(colorCharts, colorChart)
		}

		regionalColorChart[r] = colorCharts
	}

	return nil
}

func getFileContents(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	lines := make([][]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			elements := strings.Split(line, ",")
			elems := make([]int, 0)
			for _, element := range elements {
				el, _ := strconv.Atoi(element)
				elems = append(elems, el)
			}
			lines = append(lines, elems)
		}
	}

	return lines
}
