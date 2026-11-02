package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/analysis"
)

var pythonScript = regexp.MustCompile(`(?is)<script\b([^>]*)>(.*?)</script\s*>`)
var pythonJQueryCDN = regexp.MustCompile(`(?i)\bsrc\s*=\s*["']https://code\.jquery\.com/jquery-[0-9]+\.[0-9]+\.[0-9]+(?:\.min)?\.js["']`)

// browserTemplates hands literal inline JavaScript to the existing browser
// adapter. Masking preserves the original file offsets, lines and snapshot hash.
func (s *pythonSet) browserTemplates(out *analysis.AnalysisReport) error {
	parsers, err := newParsers()
	if err != nil {
		return err
	}
	set := &moduleSet{root: s.root, parsers: parsers, modules: map[string]*sourceModule{}, allowed: map[string]bool{}, textSources: s.text.textSources}
	defer set.close()
	files := []string{}
	for file := range s.text.textSources {
		files = append(files, file)
	}
	sort.Strings(files)
	for _, file := range files {
		snapshot := s.text.textSources[file]
		matches := pythonScript.FindAllStringSubmatchIndex(snapshot.Text, -1)
		if len(matches) == 0 {
			continue
		}
		original := []byte(snapshot.Text)
		masked := make([]byte, len(original))
		for i, b := range original {
			masked[i] = ' '
			if b == '\n' || b == '\r' {
				masked[i] = b
			}
		}
		globals := map[string]analysis.Evidence{}
		hasScript := false
		firstInline := false
		for _, match := range matches {
			tag := snapshot.Text[match[2]:match[3]]
			if pythonJQueryCDN.MatchString(tag) && !firstInline {
				binding := textRange(snapshot, match[0], match[4], "render_binding")
				globals["$"] = binding
				globals["jQuery"] = binding
				continue
			}
			if strings.Contains(strings.ToLower(tag), "src=") || strings.Contains(strings.ToLower(tag), "src =") {
				continue
			}
			if strings.Contains(strings.ToLower(tag), "application/json") {
				continue
			}
			firstInline = true
			copy(masked[match[4]:match[5]], original[match[4]:match[5]])
			hasScript = true
		}
		if !hasScript {
			continue
		}
		tree := parsers.js.Parse(masked, nil)
		if tree == nil {
			return fmt.Errorf("no inline JavaScript tree for %s", file)
		}
		set.modules[file] = &sourceModule{file: file, src: original, hash: snapshot.SHA256, tree: tree, funcs: map[string]*sitter.Node{}, exports: map[string]*sitter.Node{}, imports: map[string]importBinding{}, express: map[string]bool{}, receivers: map[string]bool{}, browserGlobals: globals}
		set.allowed[file] = true
	}
	return set.browser(out)
}
