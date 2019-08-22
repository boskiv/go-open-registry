package parser

// CrateDependency Struct to unmarshal Deps from cargo request
type CrateDependency struct {
	Name            string   `json:"name"`
	Req             string   `json:"req"`
	VersionReq      string   `json:"version_req"`
	Features        []string `json:"features"`
	Optional        bool     `json:"optional"`
	DefaultFeatures bool     `json:"default_features"`
	Target          *string  `json:"target"`
	Kind            string   `json:"kind"`
	Registry        *string  `json:"registry"`
	Package         *string  `json:"package"`
}

// CrateJSON Struct to unmarshal JSON from cargo request
type CrateJSON struct {
	Name     string            `json:"name"`
	Vers     string            `json:"vers"`
	Deps     []CrateDependency `json:"deps"`
	Cksum    string            `json:"cksum"`
	Features interface{}       `json:"features"`
	Yanked   bool              `json:"yanked"`
	Links    *string           `json:"links"`
}
