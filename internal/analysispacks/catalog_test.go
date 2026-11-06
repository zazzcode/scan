package analysispacks

import "testing"

func catalogItems(t *testing.T) []Component {
	t.Helper()
	items, err := descriptors()
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func TestCatalogValidation(t *testing.T) {
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []struct {
		name   string
		change func([]Component) []Component
	}{
		{"duplicate", func(c []Component) []Component { return append(c, c[0]) }},
		{"missing dependency", func(c []Component) []Component {
			c[0].Dependencies = []Requirement{{ID: "missing", Version: "1"}}
			return c
		}},
		{"dependency cycle", func(c []Component) []Component {
			c[0].Dependencies = []Requirement{{ID: c[0].ID, Version: c[0].Version}}
			return c
		}},
		{"unknown protocol", func(c []Component) []Component {
			c[0].Compatibility = []string{"future-components-v1"}
			return c
		}},
	} {
		t.Run(mutate.name, func(t *testing.T) {
			if _, err := newCatalog(mutate.change(catalogItems(t))); err == nil {
				t.Fatal("invalid catalog accepted")
			}
		})
	}
}

func TestCatalogMissingSelectorRoots(t *testing.T) {
	for _, missing := range []string{"pack.sql", "language.javascript", "language.python", "entry.http", "entry.graphql", "entry.grpc", "api.sequelize-raw-query-v1", "api.python-sqlite3-dbapi-v1", "api.python-psycopg-dbapi-v1", "api.python-psycopg2-dbapi-v1", "api.python-pymysql-dbapi-v1", "api.python-mysql.connector-dbapi-v1", "pack.nosql", "claim.nosql", "evidence.nosql", "language.nosql.javascript", "api.marsdb-collection-v1", "api.mongoose-model-query-v1", "api.mongoose-model-query-v2", "pack.validation", "claim.validation", "evidence.validation", "language.validation.javascript", "language.validation.python"} {
		t.Run(missing, func(t *testing.T) {
			items := catalogItems(t)
			filtered := items[:0]
			for _, item := range items {
				if item.ID != missing {
					filtered = append(filtered, item)
				}
			}
			if _, err := newCatalog(filtered); err == nil {
				t.Fatal("missing selector root accepted as valid embedded configuration")
			}
		})
	}
}
