package graph

// Evidence binds an observed operation to an immutable source file snapshot.
type Evidence struct {
	Kind   string `json:"kind"`
	Text   string `json:"text"`
	Span   Span   `json:"span"`
	SHA256 string `json:"sha256"`
}

// Dependency records a local value relationship; use-level gaps distinguish
// supported expression dependencies from uncertain transfers. It is not a runtime trace.
type Dependency struct {
	Source    Evidence   `json:"source"`
	Transfers []Evidence `json:"transfers"`
}

// StringUse is a candidate use needing context-specific assessment.
type StringUse struct {
	ControlEdges []ControlEdge `json:"control_edges,omitempty"`
	Check        *Check        `json:"check,omitempty"`
	Completeness *Completeness `json:"completeness,omitempty"`
	ID           string        `json:"id"`
	RouteID      string        `json:"route_id"`
	Kind         string        `json:"kind"`
	Sink         Evidence      `json:"sink"`
	Dependencies []Dependency  `json:"dependencies"`
	Context      []Evidence    `json:"context"`
	Gaps         []Gap         `json:"gaps"`
	Assessment   Assessment    `json:"assessment"`
}

// Assessment preserves the judge's result independently from structural evidence.
type Assessment struct {
	FindingStatus   string            `json:"finding_status,omitempty"`
	BlockingGaps    []Gap             `json:"blocking_gaps,omitempty"`
	Answers         map[string]string `json:"answers,omitempty"`
	Usage           map[string]int64  `json:"usage,omitempty"`
	Status          string            `json:"status"`
	Provider        string            `json:"provider,omitempty"`
	Model           string            `json:"model,omitempty"`
	QuestionVersion string            `json:"question_version,omitempty"`
	Protection      string            `json:"protection,omitempty"`
	Reason          string            `json:"reason,omitempty"`
}

// Entry links a syntactic route registration to a selected handler body.
type Entry struct {
	Route   Route     `json:"route"`
	Handler *Evidence `json:"resolved_handler,omitempty"`
	Gaps    []Gap     `json:"gaps"`
}

// StringReport is the bounded analysis artifact, separate from route-v1 JSON.
type StringReport struct {
	SelectedRoute string        `json:"selected_route,omitempty"`
	SchemaVersion int           `json:"schema_version"`
	Root          string        `json:"root"`
	Entries       []Entry       `json:"entries"`
	Uses          []StringUse   `json:"uses"`
	Gaps          []Gap         `json:"gaps"`
	Files         []Evidence    `json:"files"`
	Limitations   []string      `json:"limitations"`
	Practice      []PracticeRef `json:"practice,omitempty"`
}

// PracticeRef names the embedded language rule applied to this report.
// The full text is sent with a decision request, not copied into every use.
type PracticeRef struct {
	ID        string   `json:"id"`
	Version   string   `json:"version"`
	Languages []string `json:"languages"`
}

// Check declares a security requirement independently of the expected test answer.
type Check struct {
	ID         string   `json:"id"`
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	Kind       string   `json:"kind"`
	Fields     []string `json:"fields"`
	Validation string   `json:"validation"`
	Protection string   `json:"protection"`
}

// Completeness separates claim blockers from unrelated coverage limitations.
type Completeness struct {
	Required []string              `json:"required"`
	Evidence map[string][]Evidence `json:"evidence"`
	Blockers []Gap                 `json:"blockers"`
	Coverage []Gap                 `json:"coverage"`
	Scope    string                `json:"scope"`
}

// ControlEdge records a source-backed local successor, including both guard outcomes.
type ControlEdge struct {
	From Evidence `json:"from"`
	To   Evidence `json:"to"`
	Kind string   `json:"kind"`
}
