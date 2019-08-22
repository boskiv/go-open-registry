package parser

// CrateDependency Struct to unmarshal Deps from cargo request
type crateDependency struct {
	Name            string   `json:"name"`
	Features        []string `json:"features"`
	Optional        bool     `json:"optional"`
	DefaultFeatures bool     `json:"default_features"`
	Target          *string  `json:"target"`
	Kind            string   `json:"kind"`
	Registry        *string  `json:"registry"`
	Package         *string  `json:"package"`
}

type inCrateDependency struct {
	crateDependency
	VersionReq string `json:"version_req"`
}

type outCrateDependency struct {
	crateDependency
	Req string `json:"req"`
}

// CrateJSON Struct to unmarshal JSON from cargo request
type CrateJSON struct {
	Name     string      `json:"name"`
	Vers     string      `json:"vers"`
	Features interface{} `json:"features"`
	Yanked   bool        `json:"yanked"`
	Links    *string     `json:"links"`
}

// InCrateJSON input structure to unmarshal
type InCrateJSON struct {
	*CrateJSON
	Deps []inCrateDependency `json:"deps"`
}

// OutCrateJSON output structure to marshal
type OutCrateJSON struct {
	*CrateJSON
	Deps  []outCrateDependency `json:"deps"`
	Cksum string               `json:"cksum"`
}

// Convert method to convert input to output
// todo: refactor it
func (out OutCrateJSON) Convert(in *InCrateJSON) OutCrateJSON {
	out = OutCrateJSON{
		CrateJSON: in.CrateJSON,
		Deps:      []outCrateDependency{},
		Cksum:     "",
	}

	for i := range in.Deps {
		out.Deps = append(out.Deps, outCrateDependency{})
		out.Deps[i].crateDependency = in.Deps[i].crateDependency
		out.Deps[i].Req = in.Deps[i].VersionReq
	}

	return out
}
