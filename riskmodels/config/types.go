package config

// ColorChart object
type ColorChart struct {
	Cholesterol bool
	Diabetes    bool
	Gender      string
	Smoker      bool
	Age         int
	Chart       [][]int
}

// RegionColorChart object
type RegionColorChart map[string][]ColorChart

// Settings object
type Settings struct {
	RegionColorChart RegionColorChart
}
