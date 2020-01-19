package config

type ColorChart struct {
	Cholesterol bool
	Diabetes    bool
	Gender      string
	Smoker      bool
	Age         int
	Chart       [][]int
}

type RegionColorChart map[string][]ColorChart

type Settings struct {
	RegionColorChart RegionColorChart
}
