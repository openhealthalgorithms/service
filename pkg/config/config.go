// Package config provides configuration constants and variables
package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// CurrentSettings function selects and responds correct settings
func CurrentSettings() Settings {
	v := viper.New()
	defaults := defaultSettings()

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetConfigName("ohas")
	v.AddConfigPath("/etc/ohas/")
	v.AddConfigPath("/usr/local/ohas/")
	v.AddConfigPath("/usr/local/etc/ohas/")
	v.AddConfigPath("/var/lib/ohas/")
	v.AddConfigPath(".")

	v.ReadInConfig()

	return configSettings(v)
}

func defaultSettings() map[string]interface{} {
	settings := map[string]interface{}{
		"server.port":                  "9595",
		"files.guideline_file":         "guideline_hearts.json",
		"files.guideline_content_file": "guideline_hearts_content.json",
		"files.goal_file":              "goals_hearts.json",
		"files.goal_content_file":      "goals_hearts_content.json",
		"files.log_file":               "logs.db",
		"directories.guideline_path":   "",
		"directories.goal_path":        "",
		"directories.log_file_path":    "",
	}

	return settings
}

func configSettings(v *viper.Viper) Settings {
	settings := Settings{}

	settings.Port = v.GetString("server.port")
	settings.GuidelineFile = filepath.Join(v.GetString("directories.guideline_path"), v.GetString("files.guideline_file"))
	settings.GuidelineContentFile = filepath.Join(v.GetString("directories.guideline_path"), v.GetString("files.guideline_content_file"))
	settings.GoalFile = filepath.Join(v.GetString("directories.goal_path"), v.GetString("files.goal_file"))
	settings.GoalContentFile = filepath.Join(v.GetString("directories.goal_path"), v.GetString("files.goal_content_file"))
	settings.LogFile = filepath.Join(v.GetString("directories.log_file_path"), v.GetString("files.log_file"))

	return settings
}
