package config

// NewSettings returns a valid settings object
func NewSettings() Settings {
	settings := Settings{
		RegionColorChart: regionalColorChart,
	}

	return settings
}
