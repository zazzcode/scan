package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range []string{"src/server.ts", "src/app.py", "README.md"} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte("// fixture\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	return root
}

func TestRunListText(t *testing.T) {
	root := fixture(t)
	var stdout, stderr bytes.Buffer

	if code := run([]string{"--list", root}, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}

	out := stdout.String()
	for _, want := range []string{"src/app.py", "src/server.ts", "supported:   2", "unsupported: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "README.md") {
		t.Errorf("unsupported file listed as a source file:\n%s", out)
	}
}

func TestRunListJSON(t *testing.T) {
	root := fixture(t)
	var stdout, stderr bytes.Buffer

	if code := run([]string{"--list", "--json", root}, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}

	var got struct {
		Supported   int            `json:"supported"`
		Unsupported int            `json:"unsupported"`
		ByLanguage  map[string]int `json:"by_language"`
		Files       []struct {
			Path     string `json:"path"`
			Language string `json:"language"`
		} `json:"files"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, stdout.String())
	}

	if got.Supported != 2 || got.Unsupported != 1 {
		t.Errorf("supported=%d unsupported=%d, want 2 and 1", got.Supported, got.Unsupported)
	}
	if got.ByLanguage["typescript"] != 1 || got.ByLanguage["python"] != 1 {
		t.Errorf("by_language = %v, want one typescript and one python", got.ByLanguage)
	}
	if len(got.Files) != 2 || got.Files[0].Path != "src/app.py" {
		t.Errorf("files = %+v, want src/app.py first", got.Files)
	}
}

func TestRunUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", nil},
		{"list without a path", []string{"--list"}},
		{"unknown flag", []string{"--nope"}},
		{"two paths", []string{"--list", "a", "b"}},
		{"routes and list", []string{"--routes", "--list", "a"}},
		{"routes without a path", []string{"--routes"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != exitUsage {
				t.Errorf("exit code = %d, want %d", code, exitUsage)
			}
		})
	}
}

func TestRoutesCLI(t *testing.T) {
	root := filepath.Join("..", "..", "internal", "parse", "testdata", "express-login")
	var stdout, stderr bytes.Buffer

	if code := run([]string{"--routes", root}, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"POST /rest/user/login login.js:", "GET /ts route.ts:"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "/hidden") {
		t.Errorf("node_modules route listed:\n%s", out)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"--routes", "--json", root}, &stdout, &stderr); code != exitOK {
		t.Fatalf("json exit code = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
	var got struct {
		SchemaVersion int `json:"schema_version"`
		Routes        []struct {
			Method      string `json:"method"`
			PathPattern string `json:"path_pattern"`
			Framework   string `json:"framework"`
		} `json:"routes"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, stdout.String())
	}
	if got.SchemaVersion != 1 {
		t.Errorf("schema_version = %d, want 1", got.SchemaVersion)
	}
	found := false
	for _, route := range got.Routes {
		if route.Method == "POST" && route.PathPattern == "/rest/user/login" && route.Framework == "express" {
			found = true
		}
	}
	if !found {
		t.Fatalf("json routes = %+v, want POST /rest/user/login", got.Routes)
	}

	stdout.Reset()
	stderr.Reset()
	code := run([]string{"--routes", filepath.Join(t.TempDir(), "does-not-exist")}, &stdout, &stderr)
	if code != exitError {
		t.Errorf("missing dir exit = %d, want %d", code, exitError)
	}
}

func TestRunMissingPath(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--list", filepath.Join(t.TempDir(), "does-not-exist")}, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit code = %d, want %d", code, exitError)
	}
	if !strings.Contains(stderr.String(), "scan:") {
		t.Errorf("stderr = %q, want a scan:-prefixed error", stderr.String())
	}
}
