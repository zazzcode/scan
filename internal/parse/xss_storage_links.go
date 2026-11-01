package parse

import (
	"sort"

	"github.com/zazzcode/zazzles/internal/graph"
)

func (s *moduleSet) linkStoredXSS(report *graph.StringReport) {
	for i := range report.Uses {
		u := &report.Uses[i]
		if u.XSS == nil || u.XSS.Type != "stored" {
			continue
		}
		var dependencies []graph.Dependency
		var entries []string
		related := map[string]bool{}
		resolved := true
		for _, read := range u.Dependencies {
			if read.StorageKey == "" {
				dependencies = append(dependencies, read)
				entries = append(entries, "")
				continue
			}
			writes := append([]storageWrite{}, s.storageWrites[read.StorageKey]...)
			sort.SliceStable(writes, func(i, j int) bool {
				if writes[i].route != writes[j].route {
					return writes[i].route < writes[j].route
				}
				return writes[i].sink.Span.StartByte < writes[j].sink.Span.StartByte
			})
			if len(writes) == 0 {
				resolved = false
				dependencies = append(dependencies, read)
				entries = append(entries, "")
				continue
			}
			u.XSS.StorageKey = read.StorageKey
			for _, write := range writes {
				related[write.route] = true
				for _, origin := range write.value.deps {
					if len(dependencies) >= 64 {
						resolved = false
						continue
					}
					origin.Transfers = append(append([]graph.Evidence{}, origin.Transfers...), write.sink, read.Source)
					origin.Transfers = append(origin.Transfers, read.Transfers...)
					dependencies = append(dependencies, origin)
					entries = append(entries, write.route)
				}
			}
		}
		u.Dependencies = dependencies
		u.XSS.OriginResolved = resolved
		for id := range related {
			u.XSS.RelatedEntries = append(u.XSS.RelatedEntries, id)
		}
		sort.Strings(u.XSS.RelatedEntries)
		f := flow{set: s}
		if u.XSS.Protection != "encoded" {
			var protections []string
			u.XSS.Protection, u.XSS.ProtectionReason, protections = f.xssPathProtection(value{deps: dependencies}, u.XSS.OutputContext, false)
			u.XSS.Paths = xssPaths(protections, entries)
		}
	}
}
