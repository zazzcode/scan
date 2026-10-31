package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zazzcode/scan/internal/graph"
	"github.com/zazzcode/scan/internal/report"
)

func runExportSARIF(path string, stdout, stderr io.Writer) int {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, (64<<20)+1))
	if err != nil || len(data) > 64<<20 {
		fmt.Fprintln(stderr, "scan: saved report read failed or exceeds 64 MiB")
		return exitError
	}
	var saved graph.StringReport
	if err := json.Unmarshal(data, &saved); err != nil || saved.SchemaVersion != 1 || saved.Uses == nil {
		fmt.Fprintln(stderr, "scan: expected a schema-v1 canonical string report")
		return exitError
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report.SARIF(&saved)); err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	return exitOK
}
