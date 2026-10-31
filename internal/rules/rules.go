// Package rules is the practice pack the scanner applies while it runs.
// The files are embedded in the binary. .zazz standards are the SDLC for this
// repo and are not consulted here.
package rules

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/zazzcode/scan/internal/lang"
)

//go:embed javascript.md python.md go.md
var files embed.FS

// Practice is one language standard, loaded from the embedded rule pack.
type Practice struct {
	ID        string   `json:"id"`
	Version   string   `json:"version"`
	Languages []string `json:"languages"`
	Body      string   `json:"body"`
}

var pack []Practice

func init() {
	names := []string{"javascript.md", "python.md", "go.md"}
	for _, name := range names {
		raw, err := files.ReadFile(name)
		if err != nil {
			panic(err)
		}
		item, err := parse(string(raw))
		if err != nil {
			panic(name + ": " + err.Error())
		}
		pack = append(pack, item)
	}
}

func parse(raw string) (Practice, error) {
	item := Practice{Body: raw}
	for _, line := range strings.Split(raw, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			break
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "id":
			item.ID = value
		case "version":
			item.Version = value
		case "languages":
			for _, part := range strings.Split(value, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					item.Languages = append(item.Languages, part)
				}
			}
		default:
			return item, nil
		}
	}
	if item.ID == "" || item.Version == "" || len(item.Languages) == 0 {
		return Practice{}, fmt.Errorf("practice file needs id, version, and languages")
	}
	return item, nil
}

// ForLanguage returns the practice pack for a detected source language.
func ForLanguage(language lang.Language) (Practice, bool) {
	for _, item := range pack {
		for _, name := range item.Languages {
			if name == string(language) {
				return item, true
			}
		}
	}
	return Practice{}, false
}

// ForPaths returns each applicable practice once, ordered by id.
func ForPaths(paths []string) []Practice {
	seen := map[string]Practice{}
	for _, path := range paths {
		language, ok := lang.Detect(path)
		if !ok {
			continue
		}
		item, ok := ForLanguage(language)
		if !ok {
			continue
		}
		seen[item.ID] = item
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Practice, 0, len(ids))
	for _, id := range ids {
		item := seen[id]
		item.Body = ""
		out = append(out, item)
	}
	return out
}

// Body returns the full standard for a language, for the decision request.
func Body(language lang.Language) (Practice, bool) {
	return ForLanguage(language)
}
