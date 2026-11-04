package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zazzcode/zazzles/internal/analysis"
)

func TestAnalysisModeCompatibility(t *testing.T) {
	t.Setenv("SCAN_FINDINGS_DIR", "")
	root := "../../internal/parse/testdata/string-flow"
	var canonical []byte
	for _, mode := range []string{"--analyze", "--strings"} {
		var out, stderr bytes.Buffer
		if code := run([]string{mode, "--json", "--entry", "server.ts", root}, &out, &stderr); code != exitOK {
			t.Fatalf("%s: code=%d stderr=%s", mode, code, stderr.String())
		}
		var r analysis.AnalysisReport
		if err := json.Unmarshal(out.Bytes(), &r); err != nil {
			t.Fatal(err)
		}
		if r.SchemaVersion != analysis.ReportSchemaVersion || len(r.Uses) == 0 || len(r.AnalysisGuidelines) != 1 {
			t.Fatalf("missing current analysis metadata: %+v", r)
		}
		if canonical != nil && !bytes.Equal(canonical, out.Bytes()) {
			t.Fatal("compatibility alias changed analysis output")
		}
		canonical = append([]byte{}, out.Bytes()...)
	}
	for _, mode := range []string{"--analyze", "--strings"} {
		var out, stderr bytes.Buffer
		if code := run([]string{mode, "--routes", root}, &out, &stderr); code != exitUsage {
			t.Fatalf("conflicting modes accepted for %s", mode)
		}
	}
}

func TestSavedAnalysisSchemaExports(t *testing.T) {
	for _, input := range []string{
		`{"schema_version":1,"practice":[{"id":"python","version":"2"}],"entries":[],"uses":[],"files":[],"gaps":[]}`,
		`{"schema_version":2,"analysis_guidelines":[{"id":"python","version":"3"}],"entries":[],"uses":[],"files":[],"gaps":[]}`,
	} {
		file := filepath.Join(t.TempDir(), "saved.json")
		if err := os.WriteFile(file, []byte(input), 0600); err != nil {
			t.Fatal(err)
		}
		for _, mode := range []string{"--export-html", "--export-sarif"} {
			var out, stderr bytes.Buffer
			if code := run([]string{mode, file}, &out, &stderr); code != exitOK || out.Len() == 0 {
				t.Fatalf("%s: code=%d stderr=%s", mode, code, stderr.String())
			}
			if mode == "--export-sarif" && !json.Valid(out.Bytes()) {
				t.Fatal("invalid SARIF JSON")
			}
		}
		after, err := os.ReadFile(file)
		if err != nil || string(after) != input {
			t.Fatal("export mutated saved evidence")
		}
	}
}

func TestBrowserEventPipeline(t *testing.T) {
	t.Setenv("SCAN_FINDINGS_DIR", "")
	root := "../../internal/parse/testdata/javascript/browser-events"
	var out, stderr bytes.Buffer
	if code := run([]string{"--analyze", "--json", root}, &out, &stderr); code != exitOK {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	var r analysis.AnalysisReport
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if r.SchemaVersion != 3 || len(r.Entries) != 6 || len(r.Uses) != 4 || len(r.Findings) != 4 {
		t.Fatalf("bad actual event pipeline: %d/%d/%d", len(r.Entries), len(r.Uses), len(r.Findings))
	}
	if err := r.ValidateUses(); err != nil {
		t.Fatal(err)
	}
	saved := filepath.Join(t.TempDir(), "events.json")
	original := append([]byte{}, out.Bytes()...)
	if err := os.WriteFile(saved, original, 0600); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"--export-html", "--export-sarif"} {
		out.Reset()
		stderr.Reset()
		if code := run([]string{mode, saved}, &out, &stderr); code != exitOK {
			t.Fatalf("%s code=%d stderr=%s", mode, code, stderr.String())
		}
		if mode == "--export-sarif" && !json.Valid(out.Bytes()) {
			t.Fatal("invalid event SARIF")
		}
	}
	after, err := os.ReadFile(saved)
	if err != nil || !bytes.Equal(after, original) {
		t.Fatal("export modified evidence")
	}
}

func TestDynamicPythonVulnPyTraces(t *testing.T) {
	t.Setenv("SCAN_FINDINGS_DIR", "")
	var out, stderr bytes.Buffer
	if code := run([]string{"--analyze", "--json", "../../internal/parse/testdata/python/dynamic-vulnpy"}, &out, &stderr); code != exitOK {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	var r analysis.AnalysisReport
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if err := r.ValidateUses(); err != nil {
		t.Fatal(err)
	}
	paths := map[string]string{}
	for _, e := range r.Entries {
		if e.HTTP != nil {
			paths[e.ID] = e.HTTP.Path
		}
	}
	count := 0
	for _, f := range r.Findings {
		if strings.HasPrefix(paths[f.EntryID], "/vulnpy/") && f.Status == "poor_practice" {
			count++
		}
	}
	if count != 12 {
		t.Fatalf("generated contextual poor-practice findings=%d want12", count)
	}
}
