package judge

import (
	"github.com/zazzcode/zazzles/internal/graph"
	"testing"
)

func TestLoggingSourceProvenance(t *testing.T) {
	for _, kind := range []string{"literal_sensitive_source", "storage_source", "unresolved_value"} {
		t.Run(kind, func(t *testing.T) {
			u := graph.StringUse{Kind: "log_sink_candidate", Context: []graph.Evidence{{Kind: "logger_binding"}}, Dependencies: []graph.Dependency{{Source: graph.Evidence{Kind: kind}}}}
			a := graph.Assessment{QuestionVersion: logQuestionVersion, Answers: map[string]string{"data_class": "personal_data", "content_exposure": "yes", "redaction": "no", "log_exposure": "vulnerable"}}
			applyLogPolicy(u, &a)
			want := "confirmed_vulnerability"
			if kind == "unresolved_value" {
				want = "needs_evidence"
			}
			if a.FindingStatus != want {
				t.Fatalf("status %s want %s", a.FindingStatus, want)
			}
		})
	}
}

func TestUnknownLoggingTransformationCannotBecomeConfirmed(t *testing.T) {
	u := graph.StringUse{Kind: "log_sink_candidate", Context: []graph.Evidence{{Kind: "logger_binding"}}, Dependencies: []graph.Dependency{{Source: graph.Evidence{Kind: "request_source"}}}, Gaps: []graph.Gap{{Kind: "log_transform_unresolved"}}}
	a := graph.Assessment{QuestionVersion: logQuestionVersion, Answers: map[string]string{"data_class": "credential", "content_exposure": "yes", "redaction": "no", "log_exposure": "vulnerable"}}
	applyLogPolicy(u, &a)
	if a.FindingStatus != "needs_evidence" {
		t.Fatal("model answer discharged unknown transformation")
	}
}

func TestPublicLiteralEmailIsNotAutomaticallySensitive(t *testing.T) {
	u := graph.StringUse{Kind: "log_sink_candidate", Context: []graph.Evidence{{Kind: "logger_binding"}}, Dependencies: []graph.Dependency{{Source: graph.Evidence{Kind: "literal_sensitive_source"}}}, Gaps: []graph.Gap{{Kind: "log_literal_sensitivity_unresolved"}}}
	a := graph.Assessment{QuestionVersion: logQuestionVersion, Answers: map[string]string{"data_class": "personal_data", "content_exposure": "yes", "redaction": "no", "log_exposure": "vulnerable"}}
	applyLogPolicy(u, &a)
	if a.FindingStatus != "needs_evidence" {
		t.Fatal("email syntax alone became a confirmed sensitive leak")
	}
}
