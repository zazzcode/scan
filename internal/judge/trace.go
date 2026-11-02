package judge

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zazzcode/zazzles/internal/analysis"
)

// DecisionTrace contains the exact JSON bodies used in one assessment, never
// authentication headers. Callers decide where source-containing traces are saved.
type DecisionTrace struct {
	QuestionVersion string          `json:"question_version"`
	Request         json.RawMessage `json:"request,omitempty"`
	Response        json.RawMessage `json:"response,omitempty"`
	HTTPStatus      int             `json:"http_status"`
	RequestBytes    int             `json:"request_bytes"`
	ResponseBytes   int             `json:"response_bytes"`
}

// Replay applies the current decoder/policy to a saved decision only when its
// exact request still matches the supplied use. It makes no network request.
func (c *Client) Replay(use analysis.CandidateUse, trace DecisionTrace) analysis.Assessment {
	version := protocolVersion(use)
	a := analysis.Assessment{Status: "failed", Provider: c.provider, Model: c.model, QuestionVersion: version}
	payload, err := c.payload(use)
	if err != nil || !sameRequest(payload, trace.Request) || trace.HTTPStatus != 200 || trace.QuestionVersion != version {
		a.Reason = "Saved decision request does not exactly match the current evidence/protocol, or the response failed."
		return a
	}
	answers, model, usage, err := c.decodeQuestions(trace.Response, questionsFor(use))
	if err != nil {
		a.Reason = err.Error()
		return a
	}
	details, err := answerDetails(trace.Response, answers, questionsFor(use))
	if err != nil {
		a.Reason = err.Error()
		return a
	}
	a.Status, a.Answers, a.Model, a.Usage, a.AnswerDetails = "completed", answers, model, usage, details
	applyPolicy(use, &a)
	return a
}

// Preview prepares the exact JSON request without making a network call.
func (c *Client) Preview(use analysis.CandidateUse) (DecisionTrace, error) {
	payload, err := c.payload(use)
	version := protocolVersion(use)
	return DecisionTrace{QuestionVersion: version, Request: payload, RequestBytes: len(payload)}, err
}

// AssessTrace evaluates one use and returns its request/response evidence for
// bounded calibration. It has the same limits and no-retry behavior as Assess.
func (c *Client) AssessTrace(ctx context.Context, use analysis.CandidateUse) (analysis.Assessment, DecisionTrace) {
	var trace DecisionTrace
	assessment := c.assess(ctx, use, &trace)
	return assessment, trace
}

func sameRequest(a, b []byte) bool {
	var x, y bytes.Buffer
	if json.Compact(&x, a) != nil || json.Compact(&y, b) != nil {
		return false
	}
	return bytes.Equal(x.Bytes(), y.Bytes())
}
