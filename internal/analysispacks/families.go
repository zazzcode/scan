package analysispacks

import (
	"fmt"

	"github.com/zazzcode/zazzles/internal/analysis"
	"github.com/zazzcode/zazzles/internal/lang"
)

// Question/envelope protocols. A changed request composition for a family
// takes a new protocol identifier; historical answers under an older protocol
// are retained but cannot replay onto changed requests.
const (
	SQLProtocol        = "sql-components-v2"
	NoSQLProtocol      = "nosql-components-v1"
	ValidationProtocol = "validation-components-v1"
)

// Family registers one claim family's selection contract. Registrations are
// scanner configuration: an unknown candidate kind keeps its existing runtime
// path and never falls through to a generic trusted security judgment.
type Family struct {
	Name string
	// Label names the family in selection-failure messages.
	Label string
	// Kinds lists the candidate kinds eligible for this family.
	Kinds []string
	// QuestionNames defines the family protocol by exact question names, not
	// by a universal count.
	QuestionNames []string
	// Languages lists the source languages with verified adapters.
	Languages []lang.Language
	// Entries lists required entry capabilities; standalone claims need none.
	Entries []string
	// RequiresAPI marks families whose operations need a verified contract.
	RequiresAPI bool
	// PythonHTTPOnly restricts Python candidates to HTTP entries.
	PythonHTTPOnly bool
	Pack           string
	Claim          string
	Evidence       string
	PolicyOwner    string
}

func (f Family) supports(language lang.Language) bool {
	for _, l := range f.Languages {
		if l == language {
			return true
		}
	}
	return false
}

func (f Family) rootIDs() []string {
	ids := []string{f.Pack, f.Claim, f.Evidence, "flow.local"}
	for _, language := range f.Languages {
		ids = append(ids, languageID(&f, language))
	}
	ids = append(ids, f.Entries...)
	if f.RequiresAPI {
		for _, id := range contractIDs(&f) {
			ids = append(ids, "api."+id)
		}
	}
	return ids
}

func languageID(f *Family, language lang.Language) string {
	if f.Name == "sql" {
		// The SQL projection keeps the accepted language component identities.
		if language == lang.Python {
			return "language.python"
		}
		return "language.javascript"
	}
	name := "javascript"
	if language == lang.Python {
		name = "python"
	}
	return "language." + f.Name + "." + name
}

// Families returns the registered component families in stable order.
func Families() []Family {
	return []Family{
		{
			Name: "sql", Label: "SQL", Kinds: []string{"sql_query_candidate"}, QuestionNames: sqlQuestionNames(),
			Languages:      []lang.Language{lang.JavaScript, lang.TypeScript, lang.Python},
			Entries:        []string{"entry.http", "entry.graphql", "entry.grpc"},
			RequiresAPI:    true,
			PythonHTTPOnly: true,
			Pack:           "pack.sql", Claim: "claim.sql", Evidence: "evidence.sql",
			PolicyOwner: "internal/judge/policy.go",
		},
		{
			Name: "nosql", Label: "NoSQL", Kinds: []string{"nosql_query_candidate"}, QuestionNames: noSQLQuestionNames(),
			Languages:   []lang.Language{lang.JavaScript, lang.TypeScript},
			Entries:     []string{"entry.http", "entry.graphql", "entry.grpc"},
			RequiresAPI: true,
			Pack:        "pack.nosql", Claim: "claim.nosql", Evidence: "evidence.nosql",
			PolicyOwner: "internal/judge/policy.go",
		},
		{
			Name: "validation", Label: "Standalone validation", Kinds: []string{"input_validation_candidate"}, QuestionNames: validationQuestionNames(),
			Languages: []lang.Language{lang.JavaScript, lang.TypeScript, lang.Python},
			// A standalone field claim needs an explicit/source-backed
			// requirement and enforcing path, never a database API contract.
			Pack: "pack.validation", Claim: "claim.validation", Evidence: "evidence.validation",
			PolicyOwner: "internal/judge/validation.go",
		},
	}
}

// FamFamily returns the registered family for a candidate kind, or nil.
func FamFamily(kind string) *Family {
	for _, f := range Families() {
		for _, k := range f.Kinds {
			if k == kind {
				value := f
				return &value
			}
		}
	}
	return nil
}

// ProtocolFor returns the request protocol of a registered family kind.
func ProtocolFor(kind string) string {
	if f := FamFamily(kind); f != nil {
		return familyProtocol(f.Name)
	}
	return ""
}

func familyProtocol(name string) string {
	switch name {
	case "sql":
		return SQLProtocol
	case "nosql":
		return NoSQLProtocol
	case "validation":
		return ValidationProtocol
	}
	return ""
}

func familyContract(f *Family, id string) *analysis.APIModel {
	switch f.Name {
	case "sql":
		return SQLContract(id)
	case "nosql":
		return NoSQLContract(id)
	}
	return nil
}

func contractIDs(f *Family) []string {
	if f.Name == "sql" {
		return sqlContractIDs
	}
	return noSQLContractIDs
}

// claimQuestions decodes and validates one family's claim component against its
// registered question protocol: exact names in order, each with choices.
func (c *catalog) claimQuestions(f *Family) ([]Question, map[string][]map[string]string, error) {
	component, ok := c.components[f.Claim]
	if !ok {
		return nil, nil, fmt.Errorf("family %s claim component missing", f.Name)
	}
	var payload claimMaterial
	if err := decode(component.Material, &payload); err != nil {
		return nil, nil, fmt.Errorf("family %s claim material: %w", f.Name, err)
	}
	names := make([]string, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		names = append(names, q.Name)
		if len(payload.Choices[q.Name]) == 0 {
			return nil, nil, fmt.Errorf("family %s question %s has no choices", f.Name, q.Name)
		}
	}
	if len(names) != len(f.QuestionNames) {
		return nil, nil, fmt.Errorf("family %s claim questions %v do not match registration %v", f.Name, names, f.QuestionNames)
	}
	for i := range names {
		if names[i] != f.QuestionNames[i] {
			return nil, nil, fmt.Errorf("family %s claim questions %v do not match registration %v", f.Name, names, f.QuestionNames)
		}
	}
	return payload.Questions, payload.Choices, nil
}
