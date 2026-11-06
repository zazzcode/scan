package analysispacks

import "github.com/zazzcode/zazzles/internal/analysis"

var sqlContractIDs = []string{"sequelize-raw-query-v1", "python-sqlite3-dbapi-v1", "python-psycopg-dbapi-v1", "python-psycopg2-dbapi-v1", "python-pymysql-dbapi-v1", "python-mysql.connector-dbapi-v1"}

// SQLContract returns a copy of the supported operation contract by verified model ID.
func SQLContract(id string) *analysis.APIModel {
	if id == "sequelize-raw-query-v1" {
		return &analysis.APIModel{ID: "sequelize-raw-query-v1", Library: "sequelize", Semantics: "Sequelize query takes executable SQL text as its first argument. bind and replacements supply separated data values according to their API. model, plain and result-format options do not parameterize an already-interpolated SQL string.", Reference: "https://sequelize.org/docs/v6/core-concepts/raw-queries/"}
	}
	for _, driver := range []string{"sqlite3", "psycopg", "psycopg2", "pymysql", "mysql.connector"} {
		v := PythonContract(driver)
		if v.ID == id {
			return v
		}
	}
	return nil
}

// PythonContract applies only after import/connection/cursor provenance resolves.
func PythonContract(driver string) *analysis.APIModel {
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
	return &analysis.APIModel{ID: "python-" + driver + "-dbapi-v1", Library: driver, Reference: reference, Semantics: "execute/executemany SQL text is argument 0; argument 1 supplies separate value parameters using " + placeholder + " placeholders for this driver. Placeholders must not be quoted. Bound values cannot protect SQL already interpolated or dynamic identifiers. SQLite executescript does not bind parameters. Connection/cursor provenance is supplied separately; this contract is not evidence that the caller uses it correctly."}
}
