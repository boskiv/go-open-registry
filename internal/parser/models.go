package parser

type CrateDependency struct {
	Name            string   `json:"name"`
	Req             string   `json:"req"`
	Features        []string `json:"features"`
	Optional        bool     `json:"optional"`
	DefaultFeatures bool     `json:"default_features"`
	Target          string   `json:"target"`
	Kind            string   `json:"kind"`
	Registry        string   `json:"registry"`
	Package         string   `json:"package"`
}

type CrateJson struct {
	Name     string            `json:"name"`
	Vers     string            `json:"vers"`
	Deps     []CrateDependency `json:"deps"`
	Cksum    string            `json:"cksum"`
	Features interface{}       `json:"features"`
	Yanked   bool              `json:"yanked"`
	Links    string            `json:"links"`
}
