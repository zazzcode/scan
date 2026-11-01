package graph

// ClientEntry identifies browser code without manufacturing an HTTP route.
type ClientEntry struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Framework string `json:"framework,omitempty"`
	Location  Span   `json:"location"`
}

// XSSFacts describes a rendering claim independently of the source language or
// framework. Parser adapters establish these facts; views do not infer them.
type XSSFacts struct {
	Type             string   `json:"type"`
	EntryKind        string   `json:"entry_kind"`
	OutputContext    string   `json:"output_context"`
	ContentType      string   `json:"content_type"`
	Protection       string   `json:"protection"`
	ProtectionReason string   `json:"protection_reason"`
	RelatedEntries   []string `json:"related_entries,omitempty"`
	StorageKey       string   `json:"storage_key,omitempty"`
	OriginResolved   bool     `json:"origin_resolved"`
	ResponsePath     string   `json:"response_path,omitempty"`
	ResponseMethod   string   `json:"response_method,omitempty"`
}
