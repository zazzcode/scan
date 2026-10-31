package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zazzcode/scan/internal/graph"
	"github.com/zazzcode/scan/internal/judge"
	"github.com/zazzcode/scan/internal/parse"
	reportpkg "github.com/zazzcode/scan/internal/report"
)

func runStrings(root, entry, routeFilter, checksFile, provider, model string, limit int, asJSON, asSARIF bool, stdout, stderr io.Writer) int {
	var client *judge.Client
	if provider != "" {
		key := ""
		switch provider {
		case "openai":
			key = os.Getenv("OPENAI_API_KEY")
		case "jev":
			key = os.Getenv("TYPESAFE_API_KEY")
		}
		var err error
		client, err = judge.New(provider, model, key)
		if err != nil {
			fmt.Fprintf(stderr, "scan: %v\n", err)
			return exitError
		}
	}
	checks, err := loadChecks(checksFile)
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	report, err := parse.CheckedStrings(root, entry, checks)
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	if routeFilter != "" {
		report.SelectedRoute = routeFilter
		selected := []graph.Entry{}
		ids := map[string]bool{}
		for _, e := range report.Entries {
			if e.Route.PathPattern == routeFilter {
				selected = append(selected, e)
				ids[e.Route.ID] = true
			}
		}
		if len(selected) == 0 {
			fmt.Fprintln(stderr, "scan: selected route was not found in the supported entry inventory")
			return exitError
		}
		uses := []graph.StringUse{}
		for _, u := range report.Uses {
			if ids[u.RouteID] {
				uses = append(uses, u)
			}
		}
		report.Entries = selected
		report.Uses = uses
	}
	failed := false
	if client != nil {
		for i := range report.Uses {
			if i >= limit {
				report.Uses[i].Assessment.Status = "not_run"
				report.Uses[i].Assessment.Reason = "judge request limit reached"
				failed = true
				continue
			}
			report.Uses[i].Assessment = client.Assess(context.Background(), report.Uses[i])
			if report.Uses[i].Assessment.Status != "completed" {
				failed = true
			}
		}
	}
	if asJSON || asSARIF {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		var output any = report
		if asSARIF {
			output = reportpkg.SARIF(report)
		}
		if err := enc.Encode(output); err != nil {
			fmt.Fprintf(stderr, "scan: %v\n", err)
			return exitError
		}
	} else {
		fmt.Fprintf(stdout, "entries: %d candidate uses: %d\n", len(report.Entries), len(report.Uses))
		for _, u := range report.Uses {
			if u.Check != nil {
				fmt.Fprintf(stdout, "%s %s check=%s finding=%s protection=%s blockers=%d\n", u.Check.Method, u.Check.Path, u.Check.ID, u.Assessment.FindingStatus, u.Assessment.Protection, len(u.Assessment.BlockingGaps))
			}
			fmt.Fprintf(stdout, "%s %s:%d sources=%d gaps=%d judge=%s\n", u.Kind, u.Sink.Span.File, u.Sink.Span.StartLine, len(u.Dependencies), len(u.Gaps), u.Assessment.Status)
		}
		for _, g := range report.Gaps {
			fmt.Fprintf(stdout, "gap %s %s:%d %s\n", g.Kind, g.At.File, g.At.StartLine, g.Detail)
		}
	}
	if failed {
		return exitError
	}
	return exitOK
}
