package config

type ColorChart struct {
	Cholesterol string
	Diabetes    string
	Gender      string
	Smoker      string
	Age         int
	Chart       [][]int
}

type RegionColorChart map[string][]ColorChart

type Settings struct {
	RegionColorChart RegionColorChart
}
