package analysis

import "encoding/json"

// UnmarshalJSON reads historical schema-1 guideline metadata without retaining
// obsolete Go fields. New reports serialize only analysis_guidelines.
func (r *AnalysisReport) UnmarshalJSON(data []byte) error {
	type report AnalysisReport
	var wire struct {
		*report
		LegacyGuidelines []GuidelinesRef `json:"practice"`
	}
	var decoded AnalysisReport
	wire.report = (*report)(&decoded)
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if decoded.SchemaVersion == 1 && decoded.AnalysisGuidelines == nil {
		decoded.AnalysisGuidelines = wire.LegacyGuidelines
	}
	*r = decoded
	return nil
}
