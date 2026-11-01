package parse

import "encoding/json"

func (s *moduleSet) packageVersion(name string) string {
	// A lockfile is authoritative for resolved versions. A literal manifest pin
	// is accepted only when no lock entry exists; ranges never establish a version.
	if file, ok := s.readText("package-lock.json"); ok {
		var lock struct {
			Packages map[string]struct {
				Version string `json:"version"`
			} `json:"packages"`
			Dependencies map[string]struct {
				Version string `json:"version"`
			} `json:"dependencies"`
		}
		if json.Unmarshal([]byte(file.Text), &lock) == nil {
			if p := lock.Packages["node_modules/"+name]; p.Version != "" {
				return p.Version
			}
			if p := lock.Dependencies[name]; p.Version != "" {
				return p.Version
			}
		}
	}
	if file, ok := s.readText("package.json"); ok {
		var manifest struct {
			Dependencies map[string]string `json:"dependencies"`
		}
		if json.Unmarshal([]byte(file.Text), &manifest) == nil {
			return manifest.Dependencies[name]
		}
	}
	return ""
}
