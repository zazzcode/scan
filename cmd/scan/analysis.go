package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/zazzcode/zazzles/internal/analysis"
	"github.com/zazzcode/zazzles/internal/judge"
	"github.com/zazzcode/zazzles/internal/parse"
	reportpkg "github.com/zazzcode/zazzles/internal/report"
)

func runAnalysis(root, entry, routeFilter, checksFile, provider, model string, limit int, asJSON, asSARIF bool, outputDir string, stdout, stderr io.Writer) int {
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
	report, err := parse.AnalyzeWithChecks(root, entry, checks)
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	if routeFilter != "" {
		report.SelectedRoute = routeFilter
		selected := []analysis.EntryPoint{}
		ids := map[string]bool{}
		for _, e := range report.Entries {
			if e.HTTP != nil && e.HTTP.Path == routeFilter {
				selected = append(selected, e)
				ids[e.HTTPRoute().ID] = true
			}
		}
		if len(selected) == 0 {
			fmt.Fprintln(stderr, "scan: selected route was not found in the supported entry inventory")
			return exitError
		}
		uses := []analysis.CandidateUse{}
		for _, u := range report.Uses {
			if ids[u.EntryID] {
				uses = append(uses, u)
			}
		}
		report.SelectEntries(ids)
		if err := report.BindUses(); err != nil {
			fmt.Fprintln(stderr, err)
			return exitError
		}
	}
	failed := false
	if client != nil {
		order := make([]int, len(report.Uses))
		for i := range order {
			order[i] = i
		}
		sort.SliceStable(order, func(a, b int) bool { return queryUse(report.Uses[order[a]]) && !queryUse(report.Uses[order[b]]) })
		for requestIndex, i := range order {
			if requestIndex >= limit {
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
	reportpkg.Enrich(report)
	bundleDir := ""
	if outputDir != "" {
		bundleDir, err = reportpkg.WriteBundle(outputDir, report)
		if err != nil {
			fmt.Fprintf(stderr, "scan: output bundle incomplete at %s: %v\n", bundleDir, err)
			return exitError
		}
		if asJSON || asSARIF {
			fmt.Fprintf(stderr, "scan: findings saved to %s\n", bundleDir)
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
	} else if bundleDir != "" {
		fmt.Fprintf(stdout, "Findings saved to %s\nJSON: findings.json\nHTML: findings.html\nSARIF: findings.sarif\n", bundleDir)
	} else {
		fmt.Fprintf(stdout, "entries: %d result records: %d\n", len(report.Entries), len(report.Uses))
		for i, u := range report.Uses {
			f := report.Findings[i]
			fmt.Fprintf(stdout, "%s: %s [%s]\n%s\n", f.RecordType, f.Title, f.Status, f.Summary)
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

func queryUse(use analysis.CandidateUse) bool {
	return use.Kind == "sql_query_candidate" || use.Kind == "nosql_query_candidate"
}
