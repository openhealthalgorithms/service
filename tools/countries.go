package tools

import (
	"encoding/json"
	"io/ioutil"
)

// Country object
type Country struct {
	Name   string
	Region string
}

// CountryCode map
type CountryCode struct {
	Countries map[string]Country `json:"countries"`
}

// Countries returns list of countries
func Countries(countriesPath string) CountryCode {
	countryCodes := CountryCode{}

	countriesFile, err := ioutil.ReadFile(countriesPath)
	if err == nil {
		if er := json.Unmarshal(countriesFile, &countryCodes); er != nil {
			return countryCodes
		}
	}

	return countryCodes
}
