package config

// ColorChart object
type ColorChart struct {
	Cholesterol bool    `json:"cholesterol"`
	Diabetes    bool    `json:"diabetes"`
	Gender      string  `json:"gender"`
	Smoker      bool    `json:"smoker"`
	Age         int     `json:"age"`
	Chart       [][]int `json:"chart"`
}

// RegionColorChart object
type RegionColorChart map[string][]ColorChart

// Meta object
type Meta struct {
	Version    string `json:"version"`
	WHOVersion string `json:"whoversion"`
}

// WHOColorChart object
type WHOColorChart struct {
	Meta        Meta             `json:"meta"`
	ColorCharts RegionColorChart `json:"colorchart"`
}
