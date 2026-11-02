package parse

import "github.com/zazzcode/zazzles/internal/graph"

// Contracts apply only after import/connection/cursor provenance resolves.
// Parameter values cannot replace identifiers or repair SQL assembled earlier.
func pythonDBAPIModel(driver string) *graph.APIModel {
	reference, placeholder := "", "%s or %(name)s"
	switch driver {
	case "sqlite3":
		reference = "https://docs.python.org/3/library/sqlite3.html#sqlite3.Cursor.execute"
		placeholder = "? or :name"
	case "psycopg":
		reference = "https://www.psycopg.org/psycopg3/docs/basic/params.html"
	case "psycopg2":
		reference = "https://www.psycopg.org/docs/usage.html#passing-parameters-to-sql-queries"
	case "pymysql":
		reference = "https://pymysql.readthedocs.io/en/latest/modules/cursors.html#pymysql.cursors.Cursor.execute"
	case "mysql.connector":
		reference = "https://dev.mysql.com/doc/connector-python/en/connector-python-api-mysqlcursor-execute.html"
	default:
		return nil
	}
	return &graph.APIModel{ID: "python-" + driver + "-dbapi-v1", Library: driver, Reference: reference, Semantics: "execute/executemany SQL text is argument 0; argument 1 supplies separate value parameters using " + placeholder + " placeholders for this driver. Placeholders must not be quoted. Bound values cannot protect SQL already interpolated or dynamic identifiers. SQLite executescript does not bind parameters. Connection/cursor provenance is supplied separately; this contract is not evidence that the caller uses it correctly."}
}
