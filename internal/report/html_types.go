package report

import (
	"sort"
)

type htmlTypePill struct {
	ID, Label, Status, Tone                       string
	Count, Observed, Review, Potential, Protected int
}

func rowTypePills(items []htmlItem) []htmlTypePill {
	indexes := map[string]int{}
	pills := []htmlTypePill{}
	for _, item := range items {
		for _, id := range item.Finding.TypeIDs {
			i, ok := indexes[id]
			if !ok {
				i = len(pills)
				indexes[id] = i
				pills = append(pills, htmlTypePill{ID: id, Label: findingTypesByID[id].Label})
			}
			p := &pills[i]
			p.Count++
			switch item.Finding.Status {
			case "protected_for_claim":
				p.Protected++
			case "observed_pattern":
				p.Observed++
			case "supported_potential":
				p.Potential++
			default:
				p.Review++
			}
		}
	}
	for i := range pills {
		p := &pills[i]
		switch {
		case p.Potential == p.Count:
			p.Status = "potential issue"
			p.Tone = "observed"
		case p.Observed == p.Count:
			p.Status = "observed"
			p.Tone = "observed"
		case p.Protected == p.Count:
			p.Status = "protected"
			p.Tone = "protected"
		case p.Review == p.Count:
			p.Status = "review"
			p.Tone = "review"
		default:
			p.Status = "mixed status"
			p.Tone = "review"
		}
	}
	sort.Slice(pills, func(i, j int) bool { return pills[i].Label < pills[j].Label })
	return pills
}
