package parse

import (
	"strings"

	"github.com/zazzcode/zazzles/internal/graph"
)

// linkCrossTierXSS links a browser output that reads an API response field to the
// one server route whose JSON response provides that field. The browser path is
// joined to the server origin; unmatched, ambiguous or constant fields stay Review.
func (s *moduleSet) linkCrossTierXSS(report *graph.StringReport) {
	for i := range report.Uses {
		u := &report.Uses[i]
		x := u.XSS
		if u.Kind != "xss_candidate" || x == nil || x.Type != "cross_tier" || x.OriginResolved || x.ResponsePath == "" {
			continue
		}
		if x.ResponseField == "" || x.ResponseField == "*" {
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "response_field_unresolved", Detail: "The browser value's response field is not a single established JSON path.", At: u.Sink.Span})
			continue
		}
		var matched []jsonProducer
		routes := map[string]bool{}
		for _, p := range s.producers {
			if strings.EqualFold(p.route.Method, x.ResponseMethod) && routePatternMatch(p.route.PathPattern, x.ResponsePath) {
				matched = append(matched, p)
				routes[p.route.ID] = true
			}
		}
		if len(routes) != 1 {
			detail := "No server JSON response matches " + x.ResponseMethod + " " + x.ResponsePath + "."
			if len(routes) > 1 {
				detail = "More than one server route matches " + x.ResponseMethod + " " + x.ResponsePath + "."
			}
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "response_producer_unresolved", Detail: detail, At: u.Sink.Span})
			continue
		}
		var origins []graph.Dependency
		var sinks []graph.Evidence
		for _, p := range matched {
			field, ok := fieldValue(p.value, x.ResponseField)
			if !ok {
				continue
			}
			for _, d := range field.deps {
				origins = append(origins, d)
				sinks = append(sinks, p.sink)
			}
		}
		if len(origins) == 0 {
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "response_field_unresolved", Detail: "The producing route does not establish request-derived data in response field " + x.ResponseField + ".", At: u.Sink.Span})
			continue
		}
		var deps []graph.Dependency
		for _, d := range u.Dependencies {
			if d.Source.Kind != "api_response" {
				deps = append(deps, d)
				continue
			}
			for j, o := range origins {
				o.Transfers = append(append(append([]graph.Evidence{}, o.Transfers...), sinks[j], d.Source), d.Transfers...)
				deps = append(deps, o)
			}
		}
		if len(deps) > 64 {
			deps = deps[:64]
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "dependency_limit", Detail: "more than 64 linked cross-tier paths; remaining paths unexamined", At: u.Sink.Span})
		}
		u.Dependencies = deps
		x.OriginResolved = true
		for id := range routes {
			x.RelatedEntries = append(x.RelatedEntries, id)
		}
		f := flow{set: s}
		if x.Protection != "encoded" && x.Protection != "sanitized" || u.Sink.Kind == "angular_output" {
			if u.Sink.Kind == "angular_output" {
				x.Protection, x.ProtectionReason, x.Paths = f.angularPathProtection(value{deps: deps}, x.OutputContext)
			} else {
				var protections []string
				x.Protection, x.ProtectionReason, protections = f.xssPathProtection(value{deps: deps}, x.OutputContext, false)
				x.Paths = xssPaths(protections, nil)
			}
		}
	}
}

// routePatternMatch compares a server route pattern with a client request path;
// a parameter segment on either side matches one segment of the other.
func routePatternMatch(server, client string) bool {
	a := strings.Split(strings.Trim(server, "/"), "/")
	b := strings.Split(strings.Trim(client, "/"), "/")
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.HasPrefix(a[i], ":") || strings.HasPrefix(b[i], ":") {
			continue
		}
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
