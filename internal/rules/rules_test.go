package rules

import (
	"strings"
	"testing"

	"github.com/zazzcode/scan/internal/lang"
)

func TestPracticePack(t *testing.T) {
	js, ok := Body(lang.TypeScript)
	if !ok || js.ID != "javascript" || js.Version != "1" || !strings.Contains(js.Body, "Free text") {
		t.Fatalf("javascript practice=%+v", js)
	}
	if _, ok := Body(lang.JavaScript); !ok {
		t.Fatal("javascript files do not share the typescript practice")
	}
	py, ok := Body(lang.Python)
	if !ok || py.ID != "python" || !strings.Contains(py.Body, "DB-API") {
		t.Fatalf("python practice=%+v", py)
	}
	g, ok := Body(lang.Go)
	if !ok || g.ID != "go" || !strings.Contains(g.Body, "database/sql") {
		t.Fatalf("go practice=%+v", g)
	}
	if _, ok := Body(lang.Java); ok {
		t.Fatal("java has no practice pack")
	}
	got := ForPaths([]string{"src/server.ts", "handler.js", "app.py"})
	if len(got) != 2 || got[0].ID != "javascript" || got[1].ID != "python" || got[0].Body != "" {
		t.Fatalf("paths=%+v", got)
	}
}
